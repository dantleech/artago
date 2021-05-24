package main

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
    config := LoadConfig("testdata/config/example.yml")
    assert.Equal(t, "127.0.0.1:9999", config.Address)
    spew.Dump(config)
}
