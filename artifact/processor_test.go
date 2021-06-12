package artifact

import (
	"testing"

	"github.com/dantleech/artag/config"
	"github.com/dantleech/artag/test"
	"github.com/stretchr/testify/assert"
)

func TestApplyRules(t *testing.T) {
	workspace := test.Workspace{
		WorkspacePath: "./../workspace",
	}
	workspace.Put("test.json", []byte("foobar"))

	ta := TestAction{}
	processor := Processor{
		Rules: []config.Rule{
			{
				Predicate: "true",
				Actions: []config.Action{
					{
						Type: "do_something",
						Params: map[string]interface{}{
							"string": "Foobar",
						},
					},
				},
			},
		},
		Actions: map[string]ActionHandler{
			"do_something": ta.DoSomethingAction,
		},
	}
	processor.Process(Artifact{
		Name: "test.json",
		Path: workspace.Path("test.json"),
		Size: 0,
	})
	assert.Equal(t, "test.json", ta.artifact.Name)
}

type TestAction struct {
	artifact Artifact
}

func (ta *TestAction) DoSomethingAction(artifact Artifact, action config.Action) {
	ta.artifact = artifact
}

func TestResolveArtifactParameter(t *testing.T) {
	assert.Equal(t, "foobar", ResolveArtifactParameter(Artifact{}, "foobar"))
	assert.Equal(t, "-- foobar --", ResolveArtifactParameter(Artifact{
		Name: "foobar",
	}, "-- %artifact.Name% --"))
}
