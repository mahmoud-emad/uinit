package manager

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/uinit/internal/process"
)

type Request struct {
	Action  string `json:"action"`
	Process string `json:"process,omitempty"`
}

type Response struct {
	OK      bool            `json:"ok"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (pm *ProcessManager) sendResponse(conn net.Conn, rsp Response) error {
	encoded, err := json.Marshal(rsp)
	if err != nil {
		return err
	}

	encoded = append(encoded, '\n')

	_, err = conn.Write(encoded)
	return err
}

func (pm *ProcessManager) handleRequest(req Request) Response {
	switch req.Action {
	case "LIST":
		encoded, err := json.Marshal(pm.List())
		if err != nil {
			return errorResponse(err)
		}

		return Response{
			OK:   true,
			Data: encoded,
		}
	case "START":
		proc, err := pm.Start(req.Process)
		if err != nil {
			return errorResponse(err)
		}

		encoded, err := json.Marshal(proc)
		if err != nil {
			return errorResponse(err)
		}

		return Response{
			OK:   true,
			Data: encoded,
		}
	case "STOP":
		proc, err := pm.Stop(req.Process)
		if err != nil {
			return errorResponse(err)
		}

		encoded, err := json.Marshal(proc)
		if err != nil {
			return errorResponse(err)
		}

		return Response{
			OK:   true,
			Data: encoded,
		}
	case "LOGS", "INSPECT":
		logs, err := pm.Inspect(req.Process)
		if err != nil {
			return errorResponse(err)
		}

		encoded, err := json.Marshal(logs)
		if err != nil {
			return errorResponse(err)
		}

		return Response{
			OK:   true,
			Data: encoded,
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

func (pm *ProcessManager) List() []*process.Process {
	return pm.processes
}

func (pm *ProcessManager) Inspect(processName string) (*process.Process, error) {
	proc, err := pm.getProcessByName(processName)
	if err != nil {
		return nil, err
	}
	return proc, nil
}

func (pm *ProcessManager) Start(processName string) (*process.Process, error) {
	proc, err := pm.getProcessByName(processName)
	if err != nil {
		return nil, err
	}

	if err := proc.Start(); err != nil {
		return nil, err
	}
	return proc, nil
}

func (pm *ProcessManager) Stop(processName string) (*process.Process, error) {
	proc, err := pm.getProcessByName(processName)
	if err != nil {
		return nil, err
	}

	if err := proc.Stop(); err != nil {
		return nil, err
	}
	return proc, nil
}
