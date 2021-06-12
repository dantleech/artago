package action

import (
	"testing"

	"github.com/dantleech/artag/config"
	"github.com/dantleech/artag/processor"
	"github.com/dantleech/artag/test"
	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	workspace := test.Workspace{
		WorkspacePath: "./../workspace",
	}
	workspace.Reset()
	workspace.Put("test.json", []byte("foobar"))
	artifact := processor.Artifact{
		Path: workspace.Path("test.json"),
		Size: 100,
	}
	CopyAction(artifact, config.Action{
		Type: "copy",
		Params: map[string]interface{}{
			"destination": "./../workspace/foobar.json",
		},
	})
	assert.FileExists(t, "./../workspace/foobar.json")
	assert.FileExists(t, "./../workspace/test.json")
}
