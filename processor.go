package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Processor struct {
	Rules   []Rule
	Actions map[string]ActionHandler
}

type ActionHandler func(Artifact, Action)

type Artifact struct {
	file *os.File
}

func NewArtifact(file *os.File) Artifact {
	return Artifact{
		file: file,
	}
}

func (r Rule) isSatisfied() bool {
	return true
}

func (p Processor) process(artifact Artifact) {
	for _, rule := range p.Rules {
		if rule.isSatisfied() {
			p.applyActions(artifact, rule.Actions)
		}
	}
}

func (p Processor) applyActions(artifact Artifact, actions []Action) {
	for _, action := range actions {
		p.applyAction(artifact, action)
	}
}

func (p Processor) applyAction(artifact Artifact, action Action) {
	if _, ok := p.Actions[action.Type]; !ok {
		panic(fmt.Sprintf("Unknown action type `%v`,", action.Type))
	}

	handler := p.Actions[action.Type]
	handler(artifact, action)
}

func UnmarshallParams(params, object interface{}) {
	m, err := yaml.Marshal(params)

	if err != nil {
		panic(err)
	}

	e2 := yaml.Unmarshal(m, object)
	if e2 != nil {
		panic(e2)
	}
}
