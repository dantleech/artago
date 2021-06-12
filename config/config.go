package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Action struct {
	Type   string                 `yaml:"type"`
	Params map[string]interface{} `yaml:"params"`
}

type Rule struct {
	Predicate string   `yaml:"rule"`
	Actions   []Action `yaml:"actions"`
}

type Config struct {
	Address       string `yaml:"address"`
	WorkspacePath string `yaml:"workspacePath"`
	Rules         []Rule `yaml:"rules"`
	PublicDir     string `yaml:"publicDir"`
}

func LoadConfig(path string) Config {
	if path != "" {
		return loadConfig(path)
	}

	for _, path := range [3]string{"artago.yml", "artago.yaml", "artago.yml.dist"} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			log.Printf("Using config file `%s`", path)
			return loadConfig(path)
		}
	}

	return loadConfig("")
}

func loadConfig(path string) Config {
	if path == "" {
		return Config{
			Address:       ":8080",
			WorkspacePath: "workspace",
			Rules:         []Rule{},
			PublicDir:     "public",
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal(fmt.Sprintf("Config file not found at: %s", path))
	}

	rawConfig, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	config := Config{}

	d := yaml.NewDecoder(rawConfig)
	d.SetStrict(true)
	e := d.Decode(&config)
	if e != nil {
		log.Fatal(e.Error())
	}

	return config
}
