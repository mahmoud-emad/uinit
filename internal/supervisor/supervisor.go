package supervisor

import (
	"os"
	"path/filepath"
	"time"

	"github.com/uinit/internal/config"
)

const socketPath = "/tmp/uinit.sock"

type Supervisor struct {
	socketPath string
	services   []ManagedService
	ConfigFile string
}

func NewSupervisor(configFile string) (Supervisor, error) {
	cfg, err := config.NewConfig(configFile)
	if err != nil {
		return Supervisor{}, err
	}

	dir := filepath.Dir(socketPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return Supervisor{}, err
	}

	sup := Supervisor{
		socketPath: socketPath,
		ConfigFile: configFile,
	}

	sup.registerServices(cfg)
	return sup, nil
}

func (s *Supervisor) registerServices(cfg config.Config) {
	for _, service := range cfg.Services {
		managedService := ManagedService{
			Config: service,
			Runtime: ServiceRuntime{
				Status:   Loaded,
				LoadedAt: time.Now(),
			},
		}

		s.services = append(s.services, managedService)
	}
}

func (s *Supervisor) list() []ManagedService {
	return s.services
}
