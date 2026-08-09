package supervisor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	if err := sup.startServices(); err != nil {
		return sup, err
	}
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

func (s *Supervisor) startServices() error {
	if len(s.services) == 0 {
		return fmt.Errorf("there are no registered services to start.")
	}

	for idx, service := range s.services {
		service, err := s.startService(service)
		if err != nil {
			return err
		}

		s.services[idx] = service
	}
	return nil
}

func (s *Supervisor) startService(service ManagedService) (ManagedService, error) {
	parts := strings.Fields(service.Config.Cmd)
	if len(parts) == 0 {
		return service, fmt.Errorf("empty command for service %q", service.Config.Name)
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	service.Runtime.StartedAt = time.Now()
	service.Runtime.Status = Starting

	if err := cmd.Start(); err != nil {
		service.Runtime.Status = Exited
		return service, fmt.Errorf("failed to start %q: %w", service.Config.Name, err)
	}

	service.Runtime.PID = cmd.Process.Pid
	service.Runtime.Status = Running

	// Monitor the process without blocking the supervisor.
	go s.monitorService(service.Config.Name, cmd)

	return service, nil
}

func (s *Supervisor) monitorService(name string, cmd *exec.Cmd) {
	err := cmd.Wait()

	for idx := range s.services {
		if s.services[idx].Config.Name != name {
			continue
		}

		service := &s.services[idx]

		service.Runtime.StoppedAt = time.Now()
		service.Runtime.ExitCode = cmd.ProcessState.ExitCode()

		if err != nil {
			service.Runtime.Status = Failed
		} else {
			service.Runtime.Status = Exited
		}

		break
	}
}

func (s *Supervisor) list() []ManagedService {
	return s.services
}
