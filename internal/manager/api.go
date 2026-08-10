package manager

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/uinit/internal/process"
)

type Request struct {
	Action  string `json:"action"`
	Process string `json:"process,omitempty"`
}

type Response struct {
	OK      bool          `json:"ok"`
	Message string        `json:"message,omitempty"`
	Data    []ProcessInfo `json:"data,omitempty"`
}

func (pm *ProcessManager) sendResponse(conn net.Conn, rsp Response) error {
	encoded, err := json.Marshal(rsp)
	if err != nil {
		return err
	}

	encoded = append(encoded, '\n')
	_, err = conn.Write(encoded)
	if err != nil {
		return err
	}
	return nil
}

func (pm *ProcessManager) handleRequest(req Request) Response {
	switch req.Action {
	case "LIST":
		return Response{
			OK:   true,
			Data: pm.List(),
		}

	// Both answer with the same view, the client either prints it or reads
	// the log file it points at.
	case "INSPECT", "LOGS", "STATUS":
		info, err := pm.Inspect(req.Process)
		if err != nil {
			return errorResponse(err)
		}

		return Response{
			OK:   true,
			Data: []ProcessInfo{info},
		}
	case "START":
		info, err := pm.Start(req.Process)
		if err != nil {
			return errorResponse(err)
		}

		return Response{
			OK:   true,
			Data: []ProcessInfo{info},
		}

	default:
		return Response{
			OK: false,
			Message: fmt.Sprintf(
				"unknown action: %q",
				req.Action,
			),
		}
	}
}

func errorResponse(err error) Response {
	return Response{
		OK:      false,
		Message: err.Error(),
	}
}

func (pm *ProcessManager) List() []ProcessInfo {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	processes := make([]ProcessInfo, 0, len(pm.processes))

	for _, p := range pm.processes {
		processes = append(processes, p.Info())
	}

	return processes
}

func (pm *ProcessManager) Inspect(name string) (ProcessInfo, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for _, p := range pm.processes {
		if p.Config.Name != name {
			continue
		}

		return p.Info(), nil
	}

	return ProcessInfo{}, fmt.Errorf(
		"process %q not found",
		name,
	)
}

// startSettle is how long Start waits before reporting back. A successful
// fork/exec only means the kernel accepted the command, a process that dies
// on startup takes a moment to be reaped, and reporting "running" before then
// is a lie the next LIST contradicts.
const startSettle = 300 * time.Millisecond

func (pm *ProcessManager) Start(name string) (ProcessInfo, error) {
	pm.mu.Lock()

	idx := -1

	for i := range pm.processes {
		if pm.processes[i].Config.Name == name {
			idx = i
			break
		}
	}

	if idx == -1 {
		pm.mu.Unlock()

		return ProcessInfo{}, fmt.Errorf(
			"process %q not found",
			name,
		)
	}

	p := &pm.processes[idx]

	if p.Runtime.Status == process.Running {
		pm.mu.Unlock()

		return ProcessInfo{}, fmt.Errorf(
			"process %q is already running",
			name,
		)
	}

	if p.Runtime.Status == process.Starting {
		pm.mu.Unlock()

		return ProcessInfo{}, fmt.Errorf(
			"process %q is starting",
			name,
		)
	}

	if err := pm.startProcess(p); err != nil {
		pm.mu.Unlock()

		return ProcessInfo{}, err
	}

	pm.mu.Unlock()

	// The monitor takes pm.mu to record an early exit, so waiting while
	// holding it would guarantee we still see the stale "running".
	time.Sleep(startSettle)

	pm.mu.Lock()
	defer pm.mu.Unlock()

	// idx stays valid, processes are only appended at registration.
	info := pm.processes[idx].Info()

	if info.Status != process.Running {
		return ProcessInfo{}, fmt.Errorf(
			"process %q %s on startup (exit code %d), see %s",
			name,
			info.Status,
			info.ExitCode,
			info.LogFile,
		)
	}

	return info, nil
}
