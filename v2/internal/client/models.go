package client

import (
	"time"

	"github.com/uinit/internal/process"
)

type ProcessInfo struct {
	Name string
	Cmd  string

	LogPath string

	Status    process.Status
	PID       int
	StartedAt time.Time
	StoppedAt time.Time
	ExitCode  int
}
