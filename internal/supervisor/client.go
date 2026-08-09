package supervisor

import (
	"net"
)

type UinitClient struct {
	socketPath string
	conn       net.Conn
}

func NewClient() (UinitClient, error) {
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
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return conn, err
	}

	return conn, nil
}

func (c *UinitClient) List() (Response, error) {
	defer c.conn.Close()

	rsp, err := c.sendRequest("LIST", "")
	if err != nil {
		return Response{}, err
	}

	return rsp, nil
}
