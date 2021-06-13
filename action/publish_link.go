package action

import (
	art "github.com/dantleech/artago/artifact"
	"github.com/dantleech/artago/config"
)

type publishLinkParams struct {
	Name     string `yaml:"name"`
	Template string `yaml:"template"`
}

func PublishLinkAction(artifact art.Artifact, action config.Action) art.ActionResult {
	params := publishLinkParams{}
	art.UnmarshallParams(action.Params, &params)
	url := art.ResolveArtifactParameter(artifact, params.Template)

	m := map[string]interface{}{
		params.Name: url,
	}

	return art.ActionResult{
		Section: "links",
		Result:  m,
	}
}
