package supervisor

import (
	"time"

	"github.com/uinit/internal/config"
)

// ServiceStatus represents the current state of a managed service.
type ServiceStatus int

const (
	Loaded ServiceStatus = iota
	Starting
	Started
	Running
	Failed
	Exited
)

func (s ServiceStatus) String() string {
	switch s {
	case Loaded:
		return "loaded"
	case Starting:
		return "starting"
	case Started:
		return "started"
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

// ServiceInfo is the client facing view of a managed service.
type ServiceInfo struct {
	Name      string        `json:"name"`
	Cmd       string        `json:"cmd"`
	Status    ServiceStatus `json:"status"`
	PID       int           `json:"pid"`
	LoadedAt  time.Time     `json:"loaded_at"`
	StartedAt time.Time     `json:"started_at"`
	StoppedAt time.Time     `json:"stopped_at"`
	ExitCode  int           `json:"exit_code"`
}

func (s ManagedService) Info() ServiceInfo {
	return ServiceInfo{
		Name:      s.Config.Name,
		Cmd:       s.Config.Cmd,
		Status:    s.Runtime.Status,
		PID:       s.Runtime.PID,
		LoadedAt:  s.Runtime.LoadedAt,
		StartedAt: s.Runtime.StartedAt,
		StoppedAt: s.Runtime.StoppedAt,
		ExitCode:  s.Runtime.ExitCode,
	}
}
