package client

import (
	"fmt"
	"os"
)

func (c *UinitClient) List() ([]ProcessInfo, error) {
	defer func() { _ = c.conn.Close() }()

	rsp, err := c.sendRequest("LIST", "")
	if err != nil {
		return nil, err
	}

	var processes []ProcessInfo

	if err := c.decodeResponse(rsp.Data, &processes); err != nil {
		return nil, err
	}
	return processes, nil
}

func (c *UinitClient) Logs(processName string) (string, error) {
	defer func() { _ = c.conn.Close() }()

	rsp, err := c.sendRequest("LOGS", processName)
	if err != nil {
		return "", err
	}

	proc := ProcessInfo{}

	if err := c.decodeResponse(rsp.Data, &proc); err != nil {
		return "", err
	}

	if proc.LogPath == "" {
		return "", fmt.Errorf("No log file for %s .", processName)
	}

	logs, err := os.ReadFile(proc.LogPath)
	if err != nil {
		return "", err
	}

	if len(logs) == 0 {
		return "", fmt.Errorf("No logs yet for %s .", processName)
	}

	return string(logs), nil
}
