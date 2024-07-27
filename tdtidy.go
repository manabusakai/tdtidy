package tdtidy

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

var (
	ecsClient *ecs.Client
)

type App struct {
	ctx context.Context
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	ecsClient = ecs.NewFromConfig(cfg)
	return &App{
		ctx: ctx,
	}, nil
}

func (app *App) Run(dryRun bool, retentionPeriod int, familyPrefix string) {
	opts := options{
		dryRun:    dryRun,
		threshold: time.Now().AddDate(0, 0, -retentionPeriod).UTC(),
	}
	log.Printf("[info] threshold datetime: %s", opts.threshold.Format(time.RFC3339))

	if familyPrefix == "" {
		opts.familyPrefix = nil
	} else {
		opts.familyPrefix = &familyPrefix
	}

	if _, err := app.deregister(opts); err != nil {
		log.Fatal(err)
	}

	if _, err := app.delete(opts); err != nil {
		log.Fatal(err)
	}
}

func (app *App) deregister(opts options) (bool, error) {
	tds, err := app.getTaskDefinitions(types.TaskDefinitionStatusActive, opts.familyPrefix)
	if err != nil {
		return false, err
	}

	families, err := app.selectTaskDefinitions(opts.threshold, tds)
	if err != nil {
		return false, err
	}

	tdNames := make([]string, 0)
	for _, family := range families {
		// Keep the latest revision.
		family = family[:len(family)-1]

		tdNames = append(tdNames, family...)
	}

	if len(tdNames) == 0 {
		return true, nil
	}

	for _, tdName := range tdNames {
		if opts.dryRun {
			log.Printf("[dry-run] deregister task definition: %s", tdName)
			continue
		}

		if _, err := ecsClient.DeregisterTaskDefinition(app.ctx, &ecs.DeregisterTaskDefinitionInput{
			TaskDefinition: &tdName,
		}); err != nil {
			return false, err
		}
		log.Printf("[notice] deregister task definition: %s", tdName)

		// Avoid request throttling.
		sleep()
	}

	return true, nil
}

func (app *App) delete(opts options) (bool, error) {
	tds, err := app.getTaskDefinitions(types.TaskDefinitionStatusInactive, opts.familyPrefix)
	if err != nil {
		return false, err
	}

	families, err := app.selectTaskDefinitions(opts.threshold, tds)
	if err != nil {
		return false, err
	}

	tdNames := make([]string, 0)
	for _, family := range families {
		tdNames = append(tdNames, family...)
	}

	if len(tdNames) == 0 {
		return true, nil
	}

	chunkSize := 10
	for _, tdNames := range chunk(tdNames, chunkSize) {
		if opts.dryRun {
			log.Printf("[dry-run] delete task definitions: %v", tdNames)
			continue
		}

		if _, err := ecsClient.DeleteTaskDefinitions(app.ctx, &ecs.DeleteTaskDefinitionsInput{
			TaskDefinitions: tdNames,
		}); err != nil {
			return false, err
		}
		log.Printf("[notice] delete task definitions: %v", tdNames)

		// Avoid request throttling.
		sleep()
	}

	return true, nil
}

func (app *App) getTaskDefinitions(status types.TaskDefinitionStatus, familyPrefix *string) ([]taskdef, error) {
	p := ecs.NewListTaskDefinitionsPaginator(ecsClient, &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: familyPrefix,
		Status:       status,
	})

	tds := make([]taskdef, 0)
	for p.HasMorePages() {
		res, err := p.NextPage(app.ctx)
		if err != nil {
			return nil, err
		}

		for _, tdArn := range res.TaskDefinitionArns {
			res, err := ecsClient.DescribeTaskDefinition(app.ctx, &ecs.DescribeTaskDefinitionInput{
				TaskDefinition: &tdArn,
			})
			if err != nil {
				return nil, err
			}

			tds = append(tds, taskdef{
				family:         *res.TaskDefinition.Family,
				revision:       res.TaskDefinition.Revision,
				registeredAt:   res.TaskDefinition.RegisteredAt,
				deregisteredAt: res.TaskDefinition.DeregisteredAt,
			})
		}
	}

	return tds, nil
}

func (app *App) selectTaskDefinitions(threshold time.Time, tds []taskdef) (families, error) {
	families := make(families)
	for _, td := range tds {
		// Old task definitions do not have RegisteredAt.
		if td.registeredAt == nil {
			continue
		}

		// Exclude task definitions by registeredAt or deregisteredAt.
		if td.deregisteredAt == nil && td.registeredAt.After(threshold) {
			continue
		}
		if td.deregisteredAt != nil && td.deregisteredAt.After(threshold) {
			continue
		}

		families[td.family] = append(families[td.family], td.name())
	}

	return families, nil
}
