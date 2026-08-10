package client

import (
	"net"
)

type UinitClient struct {
	socketPath string
	conn       net.Conn
}
