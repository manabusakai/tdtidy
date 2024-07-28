package tdtidy

import (
	"context"
	"flag"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

var (
	ecsClient *ecs.Client
)

func New(ctx context.Context) (*App, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	ecsClient = ecs.NewFromConfig(cfg)

	opt, err := initOption()
	if err != nil {
		return nil, err
	}

	return &App{
		ctx: ctx,
		opt: opt,
	}, nil
}

type option struct {
	dryRun          *bool
	retentionPeriod *int
	familyPrefix    *string
}

func (opt *option) threshold() time.Time {
	return time.Now().AddDate(0, 0, -(*opt.retentionPeriod)).UTC()
}

func initOption() (*option, error) {
	var (
		dryRun          = flag.Bool("dry-run", false, "Turn on dry-run. List the target task definitions.")
		retentionPeriod = flag.Int("retention-period", 0, "Retention period for task definitions. Unit is number of days. The default value is zero.")
		familyPrefix    = flag.String("family-prefix", "", "Family name of task definitions. If specified, filter by family name.")
	)
	flag.Parse()

	if *familyPrefix == "" {
		familyPrefix = nil
	}

	return &option{
		dryRun:          dryRun,
		retentionPeriod: retentionPeriod,
		familyPrefix:    familyPrefix,
	}, nil
}
