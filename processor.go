package main

import (
	"log"
	"os"
)

type Processor struct {
	Rules []Rule
}

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
			log.Printf("Supported by rule")
		}
	}
}
