package client

import (
	"fmt"
	"os"
)

func (c *UinitClient) List() ([]ProcessInfo, error) {
	defer func() { _ = c.conn.Close() }()

	return request[[]ProcessInfo](c, "LIST", "")
}

func (c *UinitClient) Inspect(processName string) (*ProcessInfo, error) {
	defer func() { _ = c.conn.Close() }()

	return request[*ProcessInfo](c, "INSPECT", processName)
}

func (c *UinitClient) Start(processName string) (*ProcessInfo, error) {
	defer func() { _ = c.conn.Close() }()

	return request[*ProcessInfo](c, "START", processName)
}

func (c *UinitClient) Stop(processName string) (*ProcessInfo, error) {
	defer func() { _ = c.conn.Close() }()

	return request[*ProcessInfo](c, "STOP", processName)
}

func (c *UinitClient) Logs(processName string) (string, error) {
	defer func() { _ = c.conn.Close() }()

	proc, err := request[*ProcessInfo](c, "LOGS", processName)
	if err != nil {
		return "", err
	}

	if proc.LogPath == "" {
		return "", fmt.Errorf("no log file for %q", processName)
	}

	logs, err := os.ReadFile(proc.LogPath)
	if err != nil {
		return "", err
	}

	if len(logs) == 0 {
		return "", fmt.Errorf("no logs yet for %q", processName)
	}

	return string(logs), nil
}
