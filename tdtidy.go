package tdtidy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type App struct {
	ctx context.Context
	opt *option
}

func (app *App) Run() {
	var err error
	switch app.opt.subcommand {
	case Deregister:
		_, err = app.deregister()
	case Delete:
		_, err = app.delete()
	default:
		err = fmt.Errorf("unknown subcommand %q", app.opt.subcommand)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func (app *App) deregister() (bool, error) {
	tds, err := app.getTaskDefinitions(types.TaskDefinitionStatusActive, app.opt.familyPrefix)
	if err != nil {
		return false, err
	}

	families, err := app.selectTaskDefinitions(app.opt.threshold(), tds)
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
		if *app.opt.dryRun {
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

func (app *App) delete() (bool, error) {
	tds, err := app.getTaskDefinitions(types.TaskDefinitionStatusInactive, app.opt.familyPrefix)
	if err != nil {
		return false, err
	}

	families, err := app.selectTaskDefinitions(app.opt.threshold(), tds)
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
		if *app.opt.dryRun {
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
	debug.Printf("threshold: %s", threshold.Format(time.DateTime))

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
