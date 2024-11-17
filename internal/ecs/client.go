package ecs

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/samber/lo"
)

const (
	// The refill rate for API actions per second.
	// https://docs.aws.amazon.com/AmazonECS/latest/APIReference/request-throttling.html
	refillRate = 1
)

type TaskDefinition struct {
	Arn            *string
	Family         *string
	Revision       int32
	RegisteredAt   *time.Time
	DeregisteredAt *time.Time
}

type TaskDefinitionStatus types.TaskDefinitionStatus

const (
	TaskDefinitionStatusActive   TaskDefinitionStatus = "ACTIVE"
	TaskDefinitionStatusInactive TaskDefinitionStatus = "INACTIVE"
)

type Client interface {
	ListTaskDefinitionStatus() []TaskDefinitionStatus
	ListTaskDefinitions(familyPrefix *string, status TaskDefinitionStatus) ([]TaskDefinition, error)
	DeregisterTaskDefinitions(tds []TaskDefinition) ([]TaskDefinition, error)
	DeleteTaskDefinitions(tds []TaskDefinition) ([]TaskDefinition, error)
}

type client struct {
	ctx       context.Context
	ecsClient *ecs.Client
}

func NewClient(ctx context.Context) Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	ecsClient := ecs.NewFromConfig(cfg)
	return &client{
		ctx:       ctx,
		ecsClient: ecsClient,
	}
}

func (c *client) ListTaskDefinitionStatus() []TaskDefinitionStatus {
	return []TaskDefinitionStatus{
		TaskDefinitionStatusActive,
		TaskDefinitionStatusInactive,
	}
}

func (c *client) ListTaskDefinitions(familyPrefix *string, status TaskDefinitionStatus) ([]TaskDefinition, error) {
	out := make([]TaskDefinition, 0)
	p := ecs.NewListTaskDefinitionsPaginator(c.ecsClient, &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: familyPrefix,
		Status:       types.TaskDefinitionStatus(status),
	})
	for p.HasMorePages() {
		res, err := p.NextPage(c.ctx)
		if err != nil {
			return nil, err
		}
		for _, arn := range res.TaskDefinitionArns {
			res, err := c.ecsClient.DescribeTaskDefinition(c.ctx, &ecs.DescribeTaskDefinitionInput{
				TaskDefinition: &arn,
			})
			if err != nil {
				return nil, err
			}
			out = append(out, mapToTaskDefinition(res.TaskDefinition))
		}
	}
	return out, nil
}

func (c *client) DeregisterTaskDefinitions(tds []TaskDefinition) ([]TaskDefinition, error) {
	out := make([]TaskDefinition, 0)
	for _, td := range tds {
		res, err := c.ecsClient.DeregisterTaskDefinition(c.ctx, &ecs.DeregisterTaskDefinitionInput{
			TaskDefinition: td.Arn,
		})
		if err != nil {
			return nil, err
		}
		out = append(out, mapToTaskDefinition(res.TaskDefinition))
		sleep()
	}
	return out, nil
}

func (c *client) DeleteTaskDefinitions(tds []TaskDefinition) ([]TaskDefinition, error) {
	out := make([]TaskDefinition, 0)
	maxDeletions := 10
	chunks := lo.Chunk(tds, maxDeletions)
	for _, tds := range chunks {
		arns := lo.Map(tds, func(td TaskDefinition, _ int) string {
			return *td.Arn
		})
		res, err := c.ecsClient.DeleteTaskDefinitions(c.ctx, &ecs.DeleteTaskDefinitionsInput{
			TaskDefinitions: arns,
		})
		if err != nil {
			return nil, err
		}
		for _, td := range res.TaskDefinitions {
			out = append(out, mapToTaskDefinition(&td))
		}
		sleep()
	}
	return out, nil
}

func mapToTaskDefinition(td *types.TaskDefinition) TaskDefinition {
	return TaskDefinition{
		Arn:            td.TaskDefinitionArn,
		Family:         td.Family,
		Revision:       td.Revision,
		RegisteredAt:   td.RegisteredAt,
		DeregisteredAt: td.DeregisteredAt,
	}
}

func sleep() {
	time.Sleep(refillRate * time.Second)
}
