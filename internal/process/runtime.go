package process

import (
	"os/exec"
	"time"
)

type Process struct {
	cmd *exec.Cmd
}

type Status int

const (
	Loaded Status = iota
	Starting
	Running
	Failed
	Exited
)

func (s Status) String() string {
	switch s {
	case Loaded:
		return "loaded"
	case Starting:
		return "starting"
	case Running:
		return "running"
	case Failed:
		return "failed"
	case Exited:
		return "exited"
	default:
		return "unknown"
	}
}

type Runtime struct {
	Status    Status
	PID       int
	StartedAt time.Time
	StoppedAt time.Time
	ExitCode  int
}
