package manager

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

func (pm *ProcessManager) Run() error {
	listener, err := net.Listen("unix", pm.socketPath)
	if err != nil {
		if errors.Is(err, syscall.EADDRINUSE) {
			return fmt.Errorf("uinit daemon is already running")
		}

		return err
	}

	defer func() {
		_ = listener.Close()
		_ = os.Remove(pm.socketPath)
	}()

	pm.handleSignals(listener)

	log.Printf("Listening on Unix socket: %s", pm.socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}

			return err
		}

		go pm.handleConnection(conn)
	}
}

func (pm *ProcessManager) handleSignals(listener net.Listener) {
	sigChan := make(chan os.Signal, 1)

	// - On Ctrl+C the daemon closes the listener and exits, but never signals the
	// children. ping and friends survive as orphans reparented to init. That's a
	// real bug, not a style issue.
	signal.Notify(
		sigChan,
		os.Interrupt,
		syscall.SIGTERM,
	)

	go func() {
		<-sigChan

		log.Println("Shutting down uinit...")

		_ = listener.Close()
	}()
}

func (pm *ProcessManager) handleConnection(conn net.Conn) {
	defer func() { _ = conn.Close() }()

	reader := bufio.NewReader(conn)

	line, err := reader.ReadBytes('\n')
	if err != nil {
		log.Printf("failed to read request: %v", err)
		return
	}

	var req Request

	if err := json.Unmarshal(line, &req); err != nil {
		_ = pm.sendResponse(conn, Response{
			OK:      false,
			Message: fmt.Sprintf("invalid request: %v", err),
		})
		return
	}

	log.Printf(
		"request: action=%q process=%q",
		req.Action,
		req.Process,
	)

	rsp := pm.handleRequest(req)

	if err := pm.sendResponse(conn, rsp); err != nil {
		log.Printf("failed to send response: %v", err)
	}
}
