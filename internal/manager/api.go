package manager

import (
	"encoding/json"
	"fmt"
	"net"
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
