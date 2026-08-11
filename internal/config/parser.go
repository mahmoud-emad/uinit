// Package config loads the process manager configuration from YAML.
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const RuntimeDir = "/tmp/uinit"

type ProcessConfig struct {
	Name string `yaml:"name"`
	Cmd  string `yaml:"cmd"`
}

type Config struct {
	Processes []ProcessConfig `yaml:"processes"`
}

func GetSockFile() string {
	return filepath.Join(RuntimeDir, "uinit.sock")
}

func GetDaemonLogFile() string {
	return filepath.Join(GetLogDir(), "daemon.log")
}

func GetLogDir() string {
	return filepath.Join(RuntimeDir, "logs")
}

func Load(filePath string) (Config, error) {
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
