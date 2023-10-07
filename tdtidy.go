package tdtidy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type options struct {
	dryRun       bool
	threshold    time.Time
	familyPrefix *string
}

type families map[string][]string

type taskdef struct {
	family         string
	revision       int32
	registeredAt   *time.Time
	deregisteredAt *time.Time
}

func (td taskdef) name() string {
	return fmt.Sprintf("%s:%d", td.family, td.revision)
}

type App struct {
	ecs *ecs.Client
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &App{
		ecs: ecs.NewFromConfig(cfg),
	}, nil
}

func (app *App) Run(ctx context.Context, dryRun bool, retentionPeriod int, familyPrefix string) {
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

	if _, err := app.deregister(ctx, opts); err != nil {
		log.Fatal(err)
	}

	if _, err := app.delete(ctx, opts); err != nil {
		log.Fatal(err)
	}
}

func (app *App) deregister(ctx context.Context, opts options) (bool, error) {
	tdArns, err := app.getTaskDefinitionArns(ctx, types.TaskDefinitionStatusActive, opts.familyPrefix)
	if err != nil {
		return false, err
	}

	families, err := app.selectTaskDefinitions(ctx, opts.threshold, tdArns)
	if err != nil {
		return false, err
	}

	tds := make([]string, 0)
	for _, family := range families {
		// Keep the latest revision.
		family = family[:len(family)-1]

		tds = append(tds, family...)
	}

	if len(tds) == 0 {
		return true, nil
	}

	for _, td := range tds {
		if opts.dryRun {
			log.Printf("[dry-run] deregister task definition: %s", td)
			continue
		}

		if _, err := app.ecs.DeregisterTaskDefinition(ctx, &ecs.DeregisterTaskDefinitionInput{
			TaskDefinition: &td,
		}); err != nil {
			return false, err
		}
		log.Printf("[notice] deregister task definition: %s", td)
	}

	return true, nil
}

func (app *App) delete(ctx context.Context, opts options) (bool, error) {
	tdArns, err := app.getTaskDefinitionArns(ctx, types.TaskDefinitionStatusInactive, opts.familyPrefix)
	if err != nil {
		return false, err
	}

	families, err := app.selectTaskDefinitions(ctx, opts.threshold, tdArns)
	if err != nil {
		return false, err
	}

	tds := make([]string, 0)
	for _, family := range families {
		tds = append(tds, family...)
	}

	if len(tds) == 0 {
		return true, nil
	}

	chunkSize := 10
	for _, chunk := range chunk(tds, chunkSize) {
		if opts.dryRun {
			log.Printf("[dry-run] delete task definitions: %v", chunk)
			continue
		}

		if _, err := app.ecs.DeleteTaskDefinitions(ctx, &ecs.DeleteTaskDefinitionsInput{
			TaskDefinitions: chunk,
		}); err != nil {
			return false, err
		}
		log.Printf("[notice] delete task definitions: %v", chunk)
	}

	return true, nil
}

func (app *App) getTaskDefinitionArns(ctx context.Context, status types.TaskDefinitionStatus, familyPrefix *string) ([]string, error) {
	p := ecs.NewListTaskDefinitionsPaginator(app.ecs, &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: familyPrefix,
		Status:       status,
	})

	tdArns := make([]string, 0)
	for p.HasMorePages() {
		res, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		tdArns = append(tdArns, res.TaskDefinitionArns...)
	}

	return tdArns, nil
}

func (app *App) selectTaskDefinitions(ctx context.Context, threshold time.Time, tdArns []string) (families, error) {
	families := make(families)
	for _, tdArn := range tdArns {
		res, err := app.ecs.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
			TaskDefinition: &tdArn,
		})
		if err != nil {
			return nil, err
		}

		td := taskdef{
			family:         *res.TaskDefinition.Family,
			revision:       res.TaskDefinition.Revision,
			registeredAt:   res.TaskDefinition.RegisteredAt,
			deregisteredAt: res.TaskDefinition.DeregisteredAt,
		}

		// Exclude task definitions by registeredAt or deregisteredAt
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
