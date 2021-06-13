package artifact

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/antonmedv/expr"
	config "github.com/dantleech/artago/config"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

type Processor struct {
	Rules   []config.Rule
	Actions map[string]ActionHandler
}

type ActionResult struct {
	Section string
	Result  map[string]interface{}
}

type ActionHandler func(Artifact, config.Action) ActionResult

func ResolveArtifactParameter(artifact Artifact, parameter string) string {
	re := regexp.MustCompile("%.*?%")
	return re.ReplaceAllStringFunc(parameter, func(m string) string {
		env := map[string]interface{}{
			"artifact": artifact,
		}
		e := strings.Trim(m, "%")
		program, err := expr.Compile(e, expr.Env(env))

		if err != nil {
			log.Fatalf("Could not evaluate expression `%s`: %s", e, err)
		}

		result, err := expr.Run(program, env)
		return fmt.Sprintf("%s", result)
	})
}

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

func (p Processor) Process(artifact Artifact) map[string]map[string]interface{} {
	results := []ActionResult{}
	for _, rule := range p.Rules {
		if isRuleSatisfied(rule, artifact) {
			log.Printf("...applying rule `%s`", rule.Predicate)
			results = p.applyActions(artifact, rule.Actions, results)
		}
	}

	return combineResults(results)
}

func combineResults(results []ActionResult) map[string]map[string]interface{} {
	r := map[string]map[string]interface{}{}

	for _, ar := range results {
		if _, ok := r[ar.Section]; !ok {
			r[ar.Section] = map[string]interface{}{}
		}
		newResult := r[ar.Section]
		mergo.Merge(&newResult, ar.Result)
		r[ar.Section] = newResult
	}
	return r
}

func (p Processor) applyActions(artifact Artifact, actions []config.Action, results []ActionResult) []ActionResult {
	for _, action := range actions {
		ar := p.applyAction(artifact, action)
		if ar.Section != "" {
			results = append(results, ar)
		}
	}
	return results
}

func (p Processor) applyAction(artifact Artifact, action config.Action) ActionResult {
	if _, ok := p.Actions[action.Type]; !ok {
		log.Fatalf("Unknown action type `%v`,", action.Type)
	}

	log.Printf("...applying action `%s` with params `%s`", action.Type, action.Params)
	handler := p.Actions[action.Type]
	return handler(artifact, action)
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
