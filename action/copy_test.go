package action

import (
	"testing"

	"github.com/dantleech/artago/artifact"
	"github.com/dantleech/artago/config"
	"github.com/dantleech/artago/test"
	"github.com/stretchr/testify/assert"
)

func SetupWorkspace() test.Workspace {
	workspace := test.Workspace{
		WorkspacePath: "./../workspace",
	}
	workspace.Reset()
	workspace.Put("test.json", []byte("foobar"))
	return workspace
}

func TestCopyFile(t *testing.T) {
	workspace := SetupWorkspace()

	artifact := artifact.Artifact{
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

func TestCopyExpressionDestination(t *testing.T) {
	workspace := SetupWorkspace()

	artifact := artifact.Artifact{
		BuildId: "1234",
		Path:    workspace.Path("test.json"),
		Size:    100,
	}

	CopyAction(artifact, config.Action{
		Type: "copy",
		Params: map[string]interface{}{
			"destination": "./../workspace/%artifact.BuildId%/foobar.json",
		},
	})

	assert.FileExists(t, "./../workspace/1234/foobar.json")
	assert.FileExists(t, "./../workspace/test.json")
}
