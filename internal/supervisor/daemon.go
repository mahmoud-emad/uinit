package supervisor

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)

func (s *Supervisor) Run() error {
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}

	log.Printf("Listening on Linux socket: %s", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		if err := s.handleConnection(conn); err != nil {
			log.Printf("connection error: %v", err)
			continue
		}
	}
}

func (s *Supervisor) handleConnection(conn net.Conn) error {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	line, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}

	req := Request{}
	err = json.Unmarshal(line, &req)
	if err != nil {
		return err
	}

	if req.Action == "LIST" {
		err := s.list()
		if err != nil {
			return err
		}
	}
	return nil
}
