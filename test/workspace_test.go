package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkspace(t *testing.T) {
	workspace := Workspace{
		WorkspacePath: "./../workspace",
	}
	workspace.Put("test.json", []byte("foobar"))
	assert.Equal(t, []byte("foobar"), workspace.Get("test.json"))
}
