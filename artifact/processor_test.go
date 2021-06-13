package artifact

import (
	"testing"

	"github.com/dantleech/artago/config"
	"github.com/stretchr/testify/assert"
)

func TestApplyRules(t *testing.T) {
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
	processor.Process(Artifact{Name: "test.json"})
	assert.Equal(t, "test.json", ta.artifact.Name)
}

func TestCombineResults(t *testing.T) {
	ta1 := TestAction{
		result: ActionResult{Section: "Foobar", Result: map[string]interface{}{"foobar": "barfoo"}},
	}
	ta2 := TestAction{
		result: ActionResult{Section: "Foobar", Result: map[string]interface{}{"barfoo": "foobar"}},
	}
	ta3 := TestAction{
		result: ActionResult{Section: "Barfoo", Result: map[string]interface{}{"bazboo": "boobaz"}},
	}
	processor := Processor{
		Rules: []config.Rule{
			{
				Predicate: "true",
				Actions: []config.Action{
					{
						Type: "do_ta1",
						Params: map[string]interface{}{
							"string": "Foobar",
						},
					},
					{
						Type: "do_ta2",
						Params: map[string]interface{}{
							"string": "Foobar",
						},
					},
					{
						Type: "do_ta3",
						Params: map[string]interface{}{
							"string": "Foobar",
						},
					},
				},
			},
		},
		Actions: map[string]ActionHandler{
			"do_ta1": ta1.DoSomethingAction,
			"do_ta2": ta2.DoSomethingAction,
			"do_ta3": ta3.DoSomethingAction,
		},
	}
	result := processor.Process(Artifact{Name: "test.json"})
	assert.Equal(t, map[string]map[string]interface{}{
		"Foobar": {
			"barfoo": "foobar",
			"foobar": "barfoo",
		},
		"Barfoo": {
			"bazboo": "boobaz",
		},
	}, result)
}

type TestAction struct {
	result   ActionResult
	artifact Artifact
}

func (ta *TestAction) DoSomethingAction(artifact Artifact, action config.Action) ActionResult {
	ta.artifact = artifact
	return ta.result
}

func TestResolveArtifactParameter(t *testing.T) {
	assert.Equal(t, "foobar", ResolveArtifactParameter(Artifact{}, "foobar"))
	assert.Equal(t, "-- foobar --", ResolveArtifactParameter(Artifact{
		Name: "foobar",
	}, "-- %artifact.Name% --"))
}
