package tdtidy

import (
	"fmt"
	"log"
	"time"

	"github.com/manabusakai/tdtidy/internal/ecs"
)

type Processor struct {
	client ecs.Client
	opt    *option
}

func NewProcessor(client ecs.Client, opt *option) *Processor {
	return &Processor{
		client: client,
		opt:    opt,
	}
}

func (p *Processor) Process() {
	tds := make(map[ecs.TaskDefinitionStatus][]ecs.TaskDefinition)
	for _, status := range p.client.ListTaskDefinitionStatus() {
		res, err := p.client.ListTaskDefinitions(p.opt.familyPrefix, status)
		if err != nil {
			log.Fatal(err)
		}
		tds[status] = p.filterTaskDefinitions(res)
	}

	err := p.executeSubcommand(tds)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Processor) filterTaskDefinitions(tds []ecs.TaskDefinition) []ecs.TaskDefinition {
	filteredTds := make([]ecs.TaskDefinition, 0)
	for _, td := range tds {
		if p.isValidTaskDefinition(&td, p.opt.threshold()) {
			filteredTds = append(filteredTds, td)
		}
	}
	return filteredTds
}

func (p *Processor) isValidTaskDefinition(td *ecs.TaskDefinition, threshold time.Time) bool {
	// Check if RegisteredAt is missing
	if td.RegisteredAt == nil {
		return false
	}
	// If not deregistered, RegisteredAt should be before threshold
	if td.DeregisteredAt == nil && td.RegisteredAt.After(threshold) {
		return false
	}
	// If deregistered, DeregisteredAt should be before threshold
	if td.DeregisteredAt != nil && td.DeregisteredAt.After(threshold) {
		return false
	}
	return true
}

func (p *Processor) executeSubcommand(tds map[ecs.TaskDefinitionStatus][]ecs.TaskDefinition) error {
	var (
		targetTds []ecs.TaskDefinition
		action    func([]ecs.TaskDefinition) ([]ecs.TaskDefinition, error)
	)
	switch p.opt.subcommand {
	case Deregister:
		targetTds = tds[ecs.TaskDefinitionStatusActive]
		action = p.client.DeregisterTaskDefinitions
	case Delete:
		targetTds = tds[ecs.TaskDefinitionStatusInactive]
		action = p.client.DeleteTaskDefinitions
	default:
		return fmt.Errorf("unknown subcommand %q", p.opt.subcommand)
	}

	if len(targetTds) == 0 {
		return nil
	}
	if *p.opt.dryRun {
		p.printSummary(targetTds)
		return nil
	}

	res, err := action(targetTds)
	if err != nil {
		return err
	}
	p.printSummary(res)
	return nil
}

func (p *Processor) printSummary(tds []ecs.TaskDefinition) {
	logPrefixes := map[command]string{
		Deregister: "Deregistered",
		Delete:     "Deleted",
	}
	logPrefix, ok := logPrefixes[p.opt.subcommand]
	if !ok {
		logPrefix = "Unknown"
	}
	if *p.opt.dryRun {
		logPrefix = "[dry-run] " + logPrefix
	}
	for _, td := range tds {
		fmt.Printf("%s: %s:%d\n", logPrefix, *td.Family, td.Revision)
	}
}
