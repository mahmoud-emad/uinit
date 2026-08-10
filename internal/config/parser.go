// Package config loads the process list from a YAML file.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func NewConfig(filePath string) (Config, error) {
	cfg := Config{}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
