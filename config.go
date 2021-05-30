package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Action struct {
	Type   string      `yaml:"type"`
	Params interface{} `yaml:"params"`
}

type Rule struct {
	Predicate string   `yaml:"rule"`
	Actions   []Action `yaml:"actions"`
}

type Config struct {
	Address       string `yaml:"address"`
	WorkspacePath string `yaml:"workspacePath"`
	Rules         []Rule `yaml:"rules"`
}

func LoadConfig(path string) Config {
	if path != "" {
		return loadConfig(path)
	}

	for _, path := range [2]string{"artag.yml", "artag.yml.dist"} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
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
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal(fmt.Sprintf("Config file not found at: %s", path))
	}

	rawConfig, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	config := Config{}
	e := yaml.Unmarshal([]byte(rawConfig), &config)
	if e != nil {
		log.Fatal(e.Error())
	}

	return config
}
