// Package client talks to the uinit daemon over its Unix socket.
package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/uinit/internal/manager"
)

type UinitClient struct {
	socketPath string
	conn       net.Conn
}

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

func request[T any](
	c *UinitClient,
	action string,
	processName string,
) (T, error) {
	var result T

	rsp, err := c.sendRequest(action, processName)
	if err != nil {
		return result, err
	}

	if !rsp.OK {
		return result, fmt.Errorf("%s", rsp.Message)
	}

	if err := c.decodeResponse(rsp.Data, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (c *UinitClient) sendRequest(
	action string,
	processName string,
) (manager.Response, error) {
	req := manager.Request{
		Action:  action,
		Process: processName,
	}

	rsp := manager.Response{}

	encoded, err := json.Marshal(req)
	if err != nil {
		return rsp, err
	}

	encoded = append(encoded, '\n')

	_, err = c.conn.Write(encoded)
	if err != nil {
		return rsp, err
	}

	reader := bufio.NewReader(c.conn)

	line, err := reader.ReadBytes('\n')
	if err != nil {
		return rsp, err
	}

	if err := json.Unmarshal(line, &rsp); err != nil {
		return rsp, err
	}

	return rsp, nil
}

func (c *UinitClient) decodeResponse(data []byte, target any) error {
	return json.Unmarshal(data, target)
}
