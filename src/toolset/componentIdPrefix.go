package uiToolset

import (
	"fmt"
	"sync/atomic"
)

type ComponentIdPrefixGenerator struct {
	componentName string
	sequence      atomic.Uint64
}

func NewComponentIdPrefixGenerator(componentName string) *ComponentIdPrefixGenerator {
	return &ComponentIdPrefixGenerator{componentName: componentName}
}

func (generator *ComponentIdPrefixGenerator) GenerateNext() string {
	return fmt.Sprintf(
		"%s-%d", generator.componentName, generator.sequence.Add(1),
	)
}
