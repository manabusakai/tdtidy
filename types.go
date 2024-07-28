package tdtidy

import (
	"fmt"
	"time"
)

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
