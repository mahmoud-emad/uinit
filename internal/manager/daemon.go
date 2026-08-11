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

	"github.com/uinit/internal/config"
)

func (pm *ProcessManager) Run() error {
	daemonLogPath := config.GetDaemonLogFile()

	logFile, err := os.OpenFile(
		daemonLogPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("open daemon log: %w", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime)

	listener, err := net.Listen("unix", config.GetSockFile())
	if err != nil {
		if errors.Is(err, syscall.EADDRINUSE) {
			return fmt.Errorf("uinit daemon is already running")
		}

		return err
	}

	defer func() {
		_ = listener.Close()
		_ = os.Remove(config.GetSockFile())
	}()

	log.Printf("uinit daemon started")
	log.Printf("listening on Unix socket: %s", config.GetSockFile())

	pm.handleSignals(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}

			log.Printf("accept error: %v", err)
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
		if err := pm.stopProcesses(); err != nil {
			log.Printf("shutdown error: %v", err)
		}

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
