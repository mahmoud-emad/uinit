package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func NewMiniInit(filePath string) (MiniInit, error) {
	mini := MiniInit{}
	file, err := os.Open(filePath)
	if err != nil {
		return mini, err
	}

	defer file.Close()

	buf := make([]byte, 200)
	n, err := file.Read(buf)
	if err != nil {
		return mini, err
	}

	if err := yaml.Unmarshal([]byte(buf[:n]), &mini); err != nil {
		return mini, err
	}
	return mini, nil
}
