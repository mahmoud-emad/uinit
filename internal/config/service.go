package config

// Config represents the configuration loaded from YAML.
type Config struct {
	Services []Service `yaml:"services"`
}

// Service represents a service defined by the user.
type Service struct {
	Name string `yaml:"name"`
	Cmd  string `yaml:"cmd"`
}
