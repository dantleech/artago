package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	config := LoadConfig("../testdata/config/example.yml")
	assert.Equal(t, "127.0.0.1:9999", config.Address)
}

func TestLoadDefaultConfig(t *testing.T) {
	wd, _ := os.Getwd()
	os.Chdir("../testdata/config")
	config := LoadConfig("")
	assert.Equal(t, "127.0.0.1:7777", config.Address)
	os.Chdir(wd)
}
