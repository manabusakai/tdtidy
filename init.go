package tdtidy

import (
	"context"
	"errors"
	"flag"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

var (
	ecsClient *ecs.Client
)

var (
	dryRun          = flag.Bool("dry-run", false, "Turn on dry-run. Output the target task definitions.")
	retentionPeriod = flag.Int("retention-period", 0, "The retention period for task definitions is specified in days. The unit is the number of days, and the default value is zero.")
	familyPrefix    = flag.String("family-prefix", "", "Specify the family name of the task definitions. If specified, filter by family name.")
)

const (
	Deregister command = "deregister"
	Delete     command = "delete"
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

func initOption() (*option, error) {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		return nil, errors.New("subcommand not found")
	}
	cmd, args := args[0], args[1:]
	flag.CommandLine.Parse(args)

	debug.Printf("subcommand: %s", cmd)
	debug.Printf("dryRun: %t, retentionPeriod: %d, familyPrefix: %q", *dryRun, *retentionPeriod, *familyPrefix)

	if *familyPrefix == "" {
		familyPrefix = nil
	}

	return &option{
		subcommand:      command(cmd),
		dryRun:          dryRun,
		retentionPeriod: retentionPeriod,
		familyPrefix:    familyPrefix,
	}, nil
}
