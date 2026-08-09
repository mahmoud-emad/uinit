package supervisor

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
)

func (s *Supervisor) Run() error {
	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		if errors.Is(err, syscall.EADDRINUSE) {
			return fmt.Errorf("uinit daemon is already running")
		}

		return err
	}

	defer func() {
		listener.Close()
		os.Remove(s.socketPath)
	}()

	log.Printf("Listening on Unix socket: %s", s.socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		if err := s.handleConnection(conn); err != nil {
			log.Printf("connection error: %v", err)
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
		services := s.list()
		rsp := Response{
			OK:   true,
			Data: services,
		}

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
	return nil
}
