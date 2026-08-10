package manager

import (
	"fmt"
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

	// Create a copy of the user config file so the cli can read from it

	pm.loadProcesses(cfg)
	if err := pm.startProcesses(); err != nil {
		return nil, err
	}

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

func (pm *ProcessManager) startProcesses() error {
	if len(pm.processes) == 0 {
		return fmt.Errorf("there are no registered processes to start")
	}

	for _, process := range pm.processes {
		if err := process.Start(); err != nil {
			return err
		}
	}
	return nil
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
