// Package manager runs the configured processes and serves the daemon
// socket the client talks to.
package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/uinit/internal/config"
	"github.com/uinit/internal/process"
)

const (
	// SocketPath is where the daemon listens, and where the client dials.
	SocketPath = "/tmp/uinit.sock"

	logDirPath = "/tmp/uinit/logs"
)

// ManagedProcess is a configured process plus the state the daemon keeps
// for it.
type ManagedProcess struct {
	Config  config.Process
	Runtime process.Runtime
	LogFile string
}

// ProcessInfo is the client facing view of a managed process.
type ProcessInfo struct {
	Name      string         `json:"name"`
	Cmd       string         `json:"cmd"`
	Status    process.Status `json:"status"`
	PID       int            `json:"pid"`
	StartedAt time.Time      `json:"started_at"`
	StoppedAt time.Time      `json:"stopped_at"`
	ExitCode  int            `json:"exit_code"`
	LogFile   string         `json:"log_file"`
}

func (p ManagedProcess) Info() ProcessInfo {
	return ProcessInfo{
		Name:      p.Config.Name,
		Cmd:       p.Config.Cmd,
		Status:    p.Runtime.Status,
		PID:       p.Runtime.PID,
		StartedAt: p.Runtime.StartedAt,
		StoppedAt: p.Runtime.StoppedAt,
		ExitCode:  p.Runtime.ExitCode,
		LogFile:   p.LogFile,
	}
}

type ProcessManager struct {
	socketPath string
	logDirPath string
	configFile string

	// The monitors write to processes while connections read from it, so
	// both sides go through mu.
	mu        sync.Mutex
	processes []ManagedProcess
}

func NewProcessManager(configFile string) (*ProcessManager, error) {
	cfg, err := config.NewConfig(configFile)
	if err != nil {
		return nil, err
	}

	// Socket directory.
	if err := os.MkdirAll(filepath.Dir(SocketPath), 0755); err != nil {
		return nil, err
	}

	// Log directory.
	if err := os.MkdirAll(logDirPath, 0755); err != nil {
		return nil, fmt.Errorf(
			"create log directory: %w",
			err,
		)
	}

	pm := &ProcessManager{
		socketPath: SocketPath,
		logDirPath: logDirPath,
		configFile: configFile,
	}

	pm.registerProcesses(cfg)

	if err := pm.startProcesses(); err != nil {
		return nil, err
	}

	return pm, nil
}

func (pm *ProcessManager) registerProcesses(cfg config.Config) {
	for _, cfgProcess := range cfg.Processes {
		pm.processes = append(pm.processes, ManagedProcess{
			Config: cfgProcess,
			Runtime: process.Runtime{
				Status:    process.Loaded,
				StartedAt: time.Now(),
			},
		})
	}
}

func (pm *ProcessManager) startProcesses() error {
	if len(pm.processes) == 0 {
		return fmt.Errorf("there are no registered processes to start")
	}

	for i := range pm.processes {
		if err := pm.startProcess(i); err != nil {
			return err
		}
	}

	return nil
}

func (pm *ProcessManager) startProcess(index int) error {
	proc := &pm.processes[index]

	proc.Runtime.Status = process.Starting
	proc.Runtime.StartedAt = time.Now()

	logPath := filepath.Join(pm.logDirPath, proc.Config.Name+".log")

	logFile, err := os.OpenFile(
		logPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		proc.Runtime.Status = process.Failed

		return fmt.Errorf(
			"open log file for %q: %w",
			proc.Config.Name,
			err,
		)
	}

	p, err := process.Start(proc.Config.Cmd, logFile)
	if err != nil {
		_ = logFile.Close()

		proc.Runtime.Status = process.Failed

		return fmt.Errorf(
			"failed to start process %q: %w",
			proc.Config.Name,
			err,
		)
	}

	proc.Runtime.PID = p.PID()
	proc.Runtime.Status = process.Running
	proc.LogFile = logPath

	go pm.monitorProcess(index, p, logFile)

	return nil
}

func (pm *ProcessManager) monitorProcess(
	index int,
	p *process.Process,
	logFile *os.File,
) {
	defer func() { _ = logFile.Close() }()

	err := p.Wait()

	pm.mu.Lock()
	defer pm.mu.Unlock()

	proc := &pm.processes[index]

	proc.Runtime.StoppedAt = time.Now()
	proc.Runtime.ExitCode = p.ExitCode()

	if err != nil {
		proc.Runtime.Status = process.Failed
		return
	}

	proc.Runtime.Status = process.Exited
}
