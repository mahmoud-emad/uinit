package supervisor

import (
	"time"

	"github.com/uinit/internal/config"
)

// ServiceStatus represents the current state of a managed service.
type ServiceStatus int

const (
	Loaded ServiceStatus = iota
	Started
	Running
	Exited
)

func (s ServiceStatus) String() string {
	switch s {
	case Loaded:
		return "loaded"
	case Started:
		return "started"
	case Running:
		return "running"
	case Exited:
		return "exited"
	default:
		return "unknown"
	}
}

// ServiceRuntime represents the runtime state of a service.
type ServiceRuntime struct {
	Status    ServiceStatus
	PID       int
	LoadedAt  time.Time
	StartedAt time.Time
	StoppedAt time.Time
	ExitCode  int
}

type ManagedService struct {
	Config  config.Service
	Runtime ServiceRuntime
}
