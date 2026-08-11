package process

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Process struct {
	Name string
	Cmd  string

	LogPath string

	Status    Status
	PID       int
	StartedAt time.Time
	StoppedAt time.Time
	ExitCode  int

	// The process command it is running
	cmd *exec.Cmd
	mu  sync.RWMutex
}

type Status int

const (
	Loaded   Status = iota // configured but never started
	Starting               // `Start()` has been called
	Running                // process is alive
	Stopped                // uinit intentionally stopped it
	Failed                 // process terminated normally by itself
	Exited                 // process terminated unexpectedly / non-zero exit
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
	if p.Status == Running || p.Status == Starting {
		return fmt.Errorf("process is in %s state", strings.ToLower(p.Status.String()))
	}

	p.Status = Starting
	p.StartedAt = time.Now()
	p.StoppedAt = time.Time{}
	p.ExitCode = 0

	logFile, err := os.OpenFile(
		p.LogPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		p.Status = Failed

		return fmt.Errorf(
			"open log file for %q: %w",
			p.Name,
			err,
		)
	}

	parts := strings.Fields(p.Cmd)
	if len(parts) == 0 {
		logFile.Close()
		p.Status = Failed

		return fmt.Errorf("empty command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		logFile.Close()
		p.Status = Failed

		return fmt.Errorf("start process: %w", err)
	}

	p.cmd = cmd
	p.PID = cmd.Process.Pid
	p.Status = Running

	go p.monitor(cmd, logFile)

	return nil
}

func (p *Process) monitor(cmd *exec.Cmd, logFile *os.File) {
	err := cmd.Wait()

	logFile.Close()

	p.mu.Lock()
	defer p.mu.Unlock()

	p.StoppedAt = time.Now()
	p.ExitCode = cmd.ProcessState.ExitCode()

	if err != nil {
		p.Status = Failed
		return
	}

	p.Status = Exited
}

func (p *Process) Stop() error {
	return nil
}
func (p *Process) Inspect() error {
	return nil
}
func (p *Process) Logs() error {
	return nil
}
