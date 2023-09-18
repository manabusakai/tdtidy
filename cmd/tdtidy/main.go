package main

import (
	"context"
	"flag"
	"log"

	"github.com/manabusakai/tdtidy"
)

var (
	dryRun          bool
	retentionPeriod int
)

func init() {
	flag.BoolVar(&dryRun, "dry-run", false, "Turn on dry-run.")
	flag.IntVar(&retentionPeriod, "retention-period", 0, "Retention period for task definitions.")
	flag.Parse()
}

func main() {
	ctx := context.TODO()
	app, err := tdtidy.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	app.Run(ctx, dryRun, retentionPeriod)
}
