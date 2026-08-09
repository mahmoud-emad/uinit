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

	var req Request
	if err := json.Unmarshal(line, &req); err != nil {
		return s.sendResponse(conn, Response{
			OK:      false,
			Message: fmt.Sprintf("invalid request: %v", err),
		})
	}

	log.Printf("request: action=%q service=%q", req.Action, req.Service)

	switch req.Action {
	case "LIST":
		services := s.list()
		return s.sendResponse(conn, Response{
			OK:   true,
			Data: services,
		})

	case "INSPECT":
		serviceInfo, err := s.inspect(req.Service)
		if err != nil {
			log.Printf(
				"request failed: action=%q service=%q error=%v",
				req.Action,
				req.Service,
				err,
			)
			return s.sendResponse(conn, Response{
				OK:      false,
				Message: err.Error(),
			})
		}

		return s.sendResponse(conn, Response{
			OK:   true,
			Data: serviceInfo,
		})

	default:
		log.Printf(
			"unknown action: action=%q service=%q",
			req.Action,
			req.Service,
		)
		return s.sendResponse(conn, Response{
			OK:      false,
			Message: fmt.Sprintf("unknown action: %q", req.Action),
		})
	}
}
