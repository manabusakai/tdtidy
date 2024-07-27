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
	familyPrefix    string
)

func init() {
	flag.BoolVar(&dryRun, "dry-run", false, "Turn on dry-run. List the target task definitions.")
	flag.IntVar(&retentionPeriod, "retention-period", 0, "Retention period for task definitions. Unit is number of days. The default value is zero.")
	flag.StringVar(&familyPrefix, "family-prefix", "", "Family name of task definitions. If specified, filter by family name.")
	flag.Parse()
}

func main() {
	app, err := tdtidy.New(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	app.Run(dryRun, retentionPeriod, familyPrefix)
}
