package supervisor

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
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

	// Handle graceful shutdown.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down uinit...")
		listener.Close()
	}()

	log.Printf("Listening on Unix socket: %s", s.socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			// Listener was closed because of shutdown.
			return nil
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
