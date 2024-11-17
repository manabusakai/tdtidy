package tdtidy

import (
	"context"
	"errors"
	"flag"
	"time"

	"github.com/manabusakai/tdtidy/internal/ecs"
)

var (
	dryRun          = flag.Bool("dry-run", false, "Turn on dry-run. Output the target task definitions.")
	retentionPeriod = flag.Int("retention-period", 0, "The retention period for task definitions is specified in days. The unit is the number of days, and the default value is zero.")
	familyPrefix    = flag.String("family-prefix", "", "Specify the family name of the task definitions. If specified, filter by family name.")
)

func New(ctx context.Context) (*App, error) {
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

type command string

const (
	Deregister command = "deregister"
	Delete     command = "delete"
)

type option struct {
	subcommand      command
	dryRun          *bool
	retentionPeriod *int
	familyPrefix    *string
}

func (opt *option) threshold() time.Time {
	return time.Now().AddDate(0, 0, -(*opt.retentionPeriod)).UTC()
}

type App struct {
	ctx context.Context
	opt *option
}

func (app *App) Run() {
	debug.Printf("options: {subcommand: %s, dryRun: %t, retentionPeriod: %d, familyPrefix: %q}",
		app.opt.subcommand,
		*app.opt.dryRun,
		*app.opt.retentionPeriod,
		*app.opt.familyPrefix,
	)
	debug.Printf("threshold: %s", app.opt.threshold().Format(time.DateTime))

	client := ecs.NewClient(app.ctx)
	processor := NewProcessor(client, app.opt)
	processor.Process()
}
