package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/uinit/internal/client"
	"github.com/uinit/internal/process"
)

// newTable returns a writer that aligns whatever is written to it on tabs, so
// the columns below only have to be separated, never padded.
func newTable() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
}

// row writes one tab separated line. A tabwriter only buffers, so the write
// itself cannot fail in a way worth reporting; Flush is where that is checked.
func row(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

func printList(processes []client.ProcessInfo) {
	if len(processes) == 0 {
		fmt.Println("No processes loaded.")
		return
	}

	w := newTable()

	row(w, "PROCESS\tSTATUS\tPID\tUPTIME\tCOMMAND\n")
	for _, p := range processes {
		row(
			w, "%s\t%s\t%s\t%s\t%s\n",
			p.Name, status(p.Status), pid(p), uptime(p), p.Cmd,
		)
	}

	flush(w)
}

func printInspect(p *client.ProcessInfo) {
	w := newTable()

	row(w, "PROCESS\t%s\n", p.Name)
	row(w, "COMMAND\t%s\n", p.Cmd)
	row(w, "STATUS\t%s\n", status(p.Status))
	row(w, "PID\t%s\n", pid(*p))
	row(w, "UPTIME\t%s\n", uptime(*p))
	// An exit code only means something once the process is done running.
	if p.Status == process.Exited || p.Status == process.Failed {
		row(w, "EXIT CODE\t%d\n", p.ExitCode)
	}
	row(w, "LOG FILE\t%s\n", p.LogPath)

	flush(w)
}

func flush(w *tabwriter.Writer) {
	if err := w.Flush(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

// status marks the live processes with a filled dot, so a long list can be
// skimmed without reading every word.
func status(s process.Status) string {
	if s == process.Running {
		return "● " + s.String()
	}
	return "○ " + s.String()
}

// pid is only meaningful once a process has actually run.
func pid(p client.ProcessInfo) string {
	if p.PID == 0 {
		return "-"
	}
	return strconv.Itoa(p.PID)
}

// uptime is how long a running process has been up, or how long a finished one
// stayed up before it went away.
func uptime(p client.ProcessInfo) string {
	if p.StartedAt.IsZero() {
		return "-"
	}

	end := p.StoppedAt
	if p.Status == process.Running {
		end = time.Now()
	}
	if end.Before(p.StartedAt) {
		return "-"
	}

	return end.Sub(p.StartedAt).Truncate(time.Second).String()
}
