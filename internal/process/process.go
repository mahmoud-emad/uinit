// Package process starts a single OS process and reports on it.
package process

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// Start runs command and sends everything it writes to out. Where out points
// is up to the caller, the process package only wires it up.
func Start(command string, out io.Writer) (*Process, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	cmd.Stdout = out
	cmd.Stderr = out

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start process: %w", err)
	}

	return &Process{
		cmd: cmd,
	}, nil
}

func (p *Process) PID() int {
	return p.cmd.Process.Pid
}

func (p *Process) Wait() error {
	return p.cmd.Wait()
}

func (p *Process) ExitCode() int {
	if p.cmd.ProcessState == nil {
		return -1
	}

	return p.cmd.ProcessState.ExitCode()
}
