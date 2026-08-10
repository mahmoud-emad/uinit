package config

// Config represents the configuration loaded from YAML.
type Config struct {
	Processes []Process `yaml:"processes"`
}

// Process represents a process defined by the user.
type Process struct {
	Name string `yaml:"name"`
	Cmd  string `yaml:"cmd"`
}
