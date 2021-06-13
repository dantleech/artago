package action

import (
	"testing"

	art "github.com/dantleech/artago/artifact"
	"github.com/dantleech/artago/config"
	"github.com/stretchr/testify/assert"
)

func TestPublishLink(t *testing.T) {
	artifact := art.Artifact{
		Name: "foobar",
	}

	result := PublishLinkAction(artifact, config.Action{
		Type: "publishLink",
		Params: map[string]interface{}{
			"name":     "link1",
			"template": "http://%artifact.Name%",
		},
	})
	assert.Equal(t, "links", result.Section)
	assert.Equal(t, map[string]interface{}{
		"link1": "http://foobar",
	}, result.Result)
}
