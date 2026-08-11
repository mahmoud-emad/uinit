package process

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

const startSettle = 300 * time.Millisecond

type Process struct {
	Name string
	Cmd  string

	LogPath string

	Status    Status
	PID       int
	StartedAt time.Time
	StoppedAt time.Time
	ExitCode  int

	cmd           *exec.Cmd
	stopRequested bool

	mu sync.RWMutex
}

type Status int

const (
	Loaded   Status = iota // configured but never started
	Starting               // Start() has been called
	Running                // process is alive
	Stopped                // uinit intentionally stopped it
	Failed                 // process failed to start or terminated with an error
	Exited                 // process terminated normally
)

func (s Status) String() string {
	switch s {
	case Loaded:
		return "loaded"
	case Starting:
		return "starting"
	case Running:
		return "running"
	case Stopped:
		return "stopped"
	case Failed:
		return "failed"
	case Exited:
		return "exited"
	default:
		return "unknown"
	}
}

func (p *Process) Start() error {
	p.mu.Lock()

	if p.Status == Running || p.Status == Starting {
		status := p.Status
		p.mu.Unlock()

		return fmt.Errorf(
			"process %q is in %s state",
			p.Name,
			strings.ToLower(status.String()),
		)
	}

	p.Status = Starting
	p.StartedAt = time.Now()
	p.StoppedAt = time.Time{}
	p.ExitCode = 0
	p.PID = 0
	p.stopRequested = false

	p.mu.Unlock()

	logFile, err := os.OpenFile(
		p.LogPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		p.mu.Lock()
		p.Status = Failed
		p.mu.Unlock()

		return fmt.Errorf(
			"open log file for %q: %w",
			p.Name,
			err,
		)
	}

	parts := strings.Fields(p.Cmd)
	if len(parts) == 0 {
		_ = logFile.Close()

		p.mu.Lock()
		p.Status = Failed
		p.mu.Unlock()

		return fmt.Errorf(
			"empty command for process %q",
			p.Name,
		)
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		_ = logFile.Close()

		p.mu.Lock()
		p.Status = Failed
		p.mu.Unlock()

		return fmt.Errorf(
			"start process %q: %w",
			p.Name,
			err,
		)
	}

	p.mu.Lock()

	p.cmd = cmd
	p.PID = cmd.Process.Pid
	p.Status = Running

	p.mu.Unlock()

	go p.monitor(cmd, logFile)

	time.Sleep(startSettle)

	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.Status != Running {
		return fmt.Errorf(
			"process %q %s on startup (exit code %d), see %s",
			p.Name,
			p.Status,
			p.ExitCode,
			p.LogPath,
		)
	}

	return nil
}

func (p *Process) monitor(
	cmd *exec.Cmd,
	logFile *os.File,
) {
	err := cmd.Wait()

	_ = logFile.Close()

	p.mu.Lock()
	defer p.mu.Unlock()

	p.StoppedAt = time.Now()
	p.ExitCode = cmd.ProcessState.ExitCode()

	if p.stopRequested {
		p.Status = Stopped
	} else if err != nil {
		p.Status = Failed
	} else {
		p.Status = Exited
	}

	p.cmd = nil
}

func (p *Process) Stop() error {
	p.mu.Lock()

	if p.Status == Failed ||
		p.Status == Exited ||
		p.Status == Stopped ||
		p.Status == Loaded {
		status := p.Status
		p.mu.Unlock()

		return fmt.Errorf(
			"process %q is already down (%s)",
			p.Name,
			status,
		)
	}

	if p.Status == Starting {
		p.mu.Unlock()

		return fmt.Errorf(
			"process %q is still starting",
			p.Name,
		)
	}

	pid := p.PID
	p.stopRequested = true

	p.mu.Unlock()

	// Never hold our mutex while communicating with the OS.
	if err := p.terminate(pid); err != nil {
		p.mu.Lock()
		p.stopRequested = false
		p.mu.Unlock()

		return fmt.Errorf(
			"failed to stop process %q: %w",
			p.Name,
			err,
		)
	}

	// monitor() will acquire the mutex and transition
	// the process to Stopped after Wait() returns.
	time.Sleep(startSettle)

	p.mu.RLock()
	defer p.mu.RUnlock()

	return nil
}

// terminate sends SIGTERM to ask the process to gracefully exit.
func (p *Process) terminate(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	if err := proc.Signal(syscall.SIGTERM); err != nil {
		return err
	}

	return nil
}
