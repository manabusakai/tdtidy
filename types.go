package tdtidy

import (
	"fmt"
	"time"
)

type command string

type option struct {
	subcommand      command
	dryRun          *bool
	retentionPeriod *int
	familyPrefix    *string
}

func (opt *option) threshold() time.Time {
	return time.Now().AddDate(0, 0, -(*opt.retentionPeriod)).UTC()
}

type taskdef struct {
	arn            string
	family         string
	revision       int32
	registeredAt   *time.Time
	deregisteredAt *time.Time
}

func (td taskdef) name() string {
	return fmt.Sprintf("%s:%d", td.family, td.revision)
}
