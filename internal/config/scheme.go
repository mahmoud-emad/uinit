package config

// / The MiniInit struct defines the config scheme defined in the yaml files.
type MiniInit struct {
	Services []MiniService `yaml:"services"`
}

// / MiniService struct defines the service described in the yaml files.
type MiniService struct {
	Name string `yaml:"name"`
	Cmd  string `yaml:"cmd"`
}
