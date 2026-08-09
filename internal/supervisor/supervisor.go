package supervisor

import (
	"log"
	"net"
	"os"
	"path/filepath"
)

const socketPath = "/tmp/uinit.sock"

func NewSupervisor() error {
	_ = os.Remove(socketPath)

	dir := filepath.Dir(socketPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("failed to create directory: %v", err)
	}

	_, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Listening on Linux socket: %s", socketPath)
	return nil
}
