package supervisor

import (
	"log"
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

func (c *UinitClient) List() error {
	defer c.conn.Close()

	rsp, err := c.sendRequest("LIST", "")
	if err != nil {
		return err
	}
	log.Println("rsp: ", rsp)
	return nil
}

func (c *UinitClient) Init() error {
	defer c.conn.Close()

	rsp, err := c.sendRequest("LIST", "")
	if err != nil {
		return err
	}
	log.Println("rsp: ", rsp)
	return nil
}
