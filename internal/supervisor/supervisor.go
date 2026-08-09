package supervisor

import (
	"fmt"
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
	_ = os.Remove(socketPath)

	dir := filepath.Dir(socketPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return Supervisor{}, err
	}

	sup := Supervisor{
		socketPath: socketPath,
		ConfigFile: configFile,
	}

	if err := sup.registerServices(cfg); err != nil {
		return Supervisor{}, err
	}

	return sup, nil
}

func (s *Supervisor) registerServices(config config.Config) error {
	for _, service := range config.Services {
		mservice := ManagedService{
			Config: service,
			Runtime: ServiceRuntime{
				Status:   Loaded,
				LoadedAt: time.Now(),
			},
		}

		services := append(s.services, mservice)
		fmt.Println(services)
	}
	return nil
}

func (s *Supervisor) list() error {
	fmt.Println("Listing services")
	for _, service := range s.services {
		fmt.Println("service: ", service)
	}

	return nil
}
