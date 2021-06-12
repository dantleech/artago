package main

import (
	"testing"

	"github.com/dantleech/artag/test"
)

func TestApplyRules(t *testing.T) {
	workspace := test.Workspace{
		WorkspacePath: "./workspace",
	}
	workspace.Put("test.json", []byte("foobar"))

	processor := Processor{
		Rules: []Rule{
			{
				Predicate: "true",
				Actions: []Action{
					{
						Type: "echo",
						Params: map[string]interface{}{
							"string": "Foobar",
						},
					},
				},
			},
		},
		Actions: map[string]ActionHandler{
			"echo": EchoAction,
		},
	}
	processor.process(Artifact{file: workspace.Open("test.json")})
}
