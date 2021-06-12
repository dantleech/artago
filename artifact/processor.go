package artifact

import (
	"log"

	"github.com/antonmedv/expr"
	config "github.com/dantleech/artag/config"
	"gopkg.in/yaml.v2"
)

type Processor struct {
	Rules   []config.Rule
	Actions map[string]ActionHandler
}

type ActionHandler func(Artifact, config.Action)

func isRuleSatisfied(rule config.Rule, artifact Artifact) bool {
	env := map[string]interface{}{
		"artifact": artifact,
	}

	program, err := expr.Compile(rule.Predicate, expr.Env(env), expr.AsBool())

	if err != nil {
		log.Fatalf("Could not evaluate expression `%s`: %s", rule.Predicate, err)
	}

	result, err := expr.Run(program, env)

	return result.(bool)

}

func (p Processor) Process(artifact Artifact) {
	for _, rule := range p.Rules {
		if isRuleSatisfied(rule, artifact) {
			log.Printf("...applying rule `%s`", rule.Predicate)
			p.applyActions(artifact, rule.Actions)
		}
	}
}

func (p Processor) applyActions(artifact Artifact, actions []config.Action) {
	for _, action := range actions {
		p.applyAction(artifact, action)
	}
}

func (p Processor) applyAction(artifact Artifact, action config.Action) {
	if _, ok := p.Actions[action.Type]; !ok {
		log.Fatalf("Unknown action type `%v`,", action.Type)
	}

	log.Printf("...applying action `%s` with params `%s`", action.Type, action.Params)
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
