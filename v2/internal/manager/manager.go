package manager

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/uinit/internal/config"
	"github.com/uinit/internal/process"
)

type ProcessManager struct {
	cfg       config.Config
	processes []*process.Process
}

// Return a new process manager instance with the config configured in the `yaml` config file.
func NewManager(configFilePath string) (*ProcessManager, error) {
	cfg, err := config.Load(configFilePath)
	if err != nil {
		return nil, err
	}

	pm := ProcessManager{
		cfg: cfg,
	}

	pm.loadProcesses(cfg)

	pm.startProcesses()

	return &pm, nil
}

func (pm *ProcessManager) loadProcesses(cfg config.Config) {
	for _, cfgProcess := range cfg.Processes {
		logPath := filepath.Join(
			config.GetLogDir(),
			cfgProcess.Name+".log",
		)

		pm.processes = append(pm.processes, &process.Process{
			Name:    cfgProcess.Name,
			Cmd:     cfgProcess.Cmd,
			Status:  process.Loaded,
			LogPath: logPath,
		})
	}
}

func (pm *ProcessManager) startProcesses() {
	for _, p := range pm.processes {
		if err := p.Start(); err != nil {
			log.Printf(
				"failed to start process %q: %v",
				p.Name,
				err,
			)
		}
	}
}

func (pm *ProcessManager) stopProcesses() error {
	if len(pm.processes) == 0 {
		return fmt.Errorf("there are no registered processes to start")
	}

	for _, process := range pm.processes {
		if err := process.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (pm *ProcessManager) getProcessByName(processName string) (*process.Process, error) {
	for _, process := range pm.processes {
		if process.Name != processName {
			continue
		}
		return process, nil
	}
	return nil, fmt.Errorf("cannot find process with name %s", processName)
}
