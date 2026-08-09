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
