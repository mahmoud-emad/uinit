// Package client talks to the uinit daemon over its Unix socket.
package client

import (
	"net"

	"github.com/uinit/internal/manager"
)

func NewClient(socketPath string) (UinitClient, error) {
	cli := UinitClient{
		socketPath: socketPath,
	}

	conn, err := cli.connect()
	if err != nil {
		return cli, err
	}

	cli.conn = conn
	return cli, nil
}

func (c *UinitClient) connect() (net.Conn, error) {
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return conn, err
	}

	return conn, nil
}

func (c *UinitClient) List() (manager.Response, error) {
	defer func() { _ = c.conn.Close() }()

	rsp, err := c.sendRequest("LIST", "")
	if err != nil {
		return manager.Response{}, err
	}

	return rsp, nil
}

func (c *UinitClient) Inspect(processName string) (manager.Response, error) {
	defer func() { _ = c.conn.Close() }()

	rsp, err := c.sendRequest("INSPECT", processName)
	if err != nil {
		return manager.Response{}, err
	}

	return rsp, nil
}

func (c *UinitClient) Logs(processName string) (manager.Response, error) {
	defer func() { _ = c.conn.Close() }()

	rsp, err := c.sendRequest("LOGS", processName)
	if err != nil {
		return manager.Response{}, err
	}

	return rsp, nil
}

func (c *UinitClient) Status(processName string) (manager.Response, error) {
	defer func() { _ = c.conn.Close() }()

	rsp, err := c.sendRequest("STATUS", processName)
	if err != nil {
		return manager.Response{}, err
	}

	return rsp, nil
}
