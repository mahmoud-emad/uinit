package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/uinit/internal/supervisor"
)

// ANSI escapes used to style the output.
const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	dim    = "\033[2m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
)

// colored reports whether the output supports ANSI styling. Piping the
// output to a file or another program keeps it plain.
var colored = supportsColor()

func supportsColor() bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}

	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}

func style(text string, codes ...string) string {
	if !colored || len(codes) == 0 {
		return text
	}

	return strings.Join(codes, "") + text + reset
}

// pad right aligns text into width, counting runes so that the multi byte
// status glyphs still line up.
func pad(text string, width int) string {
	if n := utf8.RuneCountInString(text); n < width {
		return text + strings.Repeat(" ", width-n)
	}

	return text
}

func truncate(text string, width int) string {
	if utf8.RuneCountInString(text) <= width {
		return text
	}

	return string([]rune(text)[:width-1]) + "…"
}

func statusStyle(status supervisor.ServiceStatus) (string, string) {
	switch status {
	case supervisor.Running:
		return "● running", green
	case supervisor.Started:
		return "● started", green
	case supervisor.Starting:
		return "◌ starting", yellow
	case supervisor.Exited:
		return "○ exited", dim
	case supervisor.Failed:
		return "✕ failed", red
	default:
		return "· loaded", cyan
	}
}

func formatStatus(status supervisor.ServiceStatus, width int) string {
	label, color := statusStyle(status)

	return style(pad(label, width), color)
}

func formatPID(pid int) string {
	if pid == 0 {
		return "-"
	}

	return fmt.Sprintf("%d", pid)
}

func formatAge(service supervisor.ServiceInfo) string {
	if service.StartedAt.IsZero() {
		return "-"
	}

	end := time.Now()

	if !service.StoppedAt.IsZero() {
		end = service.StoppedAt
	}

	return formatDuration(end.Sub(service.StartedAt))
}

func formatDuration(d time.Duration) string {
	switch {
	case d < time.Second:
		return "<1s"
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}
