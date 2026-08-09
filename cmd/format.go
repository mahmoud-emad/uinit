package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/uinit/internal/config"
	"github.com/uinit/internal/supervisor"
)

// lipgloss drops the colors on its own when the output is piped or NO_COLOR
// is set, so nothing here needs to check for a terminal.
var (
	bold   = lipgloss.NewStyle().Bold(true)
	dim    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	red    = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	cyan   = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

	label = dim.Width(12)
	cell  = lipgloss.NewStyle().PaddingRight(2)
)

const ruleWidth = 52

// printStartup prints what the daemon is about to run.
func printStartup(configFile string, cfg config.Config) {
	fmt.Println()
	fmt.Println(bold.Render("uinit"))
	fmt.Println(rule())

	printField("Config", configFile)
	printField("Services", fmt.Sprintf("%d", len(cfg.Services)))
	fmt.Println()

	width := 0
	for _, s := range cfg.Services {
		width = max(width, lipgloss.Width(s.Name))
	}

	name := lipgloss.NewStyle().Width(width + 2)
	for _, s := range cfg.Services {
		fmt.Println(name.Render(s.Name) + dim.Render(s.Cmd))
	}

	fmt.Println()
	fmt.Println(dim.Render("Press Ctrl+C to stop."))
	fmt.Println()
}

// printList prints one or more services as a table.
func printList(services []supervisor.ServiceInfo) {
	if len(services) == 0 {
		fmt.Println(dim.Render("No services loaded."))
		return
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(false).BorderBottom(false).BorderLeft(false).BorderRight(false).
		BorderColumn(false).
		BorderStyle(dim).
		Headers("SERVICE", "STATUS", "PID", "UPTIME", "COMMAND").
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == table.HeaderRow {
				return cell.Bold(true)
			}
			return cell
		})

	running := 0

	for _, s := range services {
		if s.Status == supervisor.Running || s.Status == supervisor.Started {
			running++
		}

		t.Row(s.Name, formatStatus(s.Status), formatPID(s.PID), formatAge(s), dim.Render(s.Cmd))
	}

	fmt.Println()
	fmt.Println(t)
	fmt.Println()
	fmt.Println(dim.Render(fmt.Sprintf("%d services · %d running", len(services), running)))
	fmt.Println()
}

// printInspect prints one or more services in full, one block each.
func printInspect(services []supervisor.ServiceInfo) {
	if len(services) == 0 {
		fmt.Println(dim.Render("No services found."))
		return
	}

	for _, s := range services {
		fmt.Println()
		fmt.Println("Service: " + bold.Render(s.Name))
		fmt.Println(rule())

		printField("Status", formatStatus(s.Status))
		printField("PID", formatPID(s.PID))
		printField("Command", s.Cmd)

		fmt.Println()
		printField("Loaded", formatTime(s.LoadedAt))
		printField("Started", formatTime(s.StartedAt))
		printField("Stopped", formatTime(s.StoppedAt))
		printField("Exit code", formatExitCode(s))
		printField("Duration", formatAge(s))

		// Error and Logs go here once ServiceInfo carries them.

		fmt.Println()
	}
}

func printField(name, value string) {
	fmt.Println(label.Render(name+":") + value)
}

func rule() string {
	return dim.Render(strings.Repeat("─", ruleWidth))
}

func formatStatus(status supervisor.ServiceStatus) string {
	switch status {
	case supervisor.Running, supervisor.Started:
		return green.Render("● " + status.String())
	case supervisor.Starting:
		return yellow.Render("◌ " + status.String())
	case supervisor.Failed:
		return red.Render("✕ " + status.String())
	case supervisor.Exited:
		return dim.Render("○ " + status.String())
	default:
		return cyan.Render("· " + status.String())
	}
}

func formatPID(pid int) string {
	if pid == 0 {
		return "-"
	}

	return fmt.Sprintf("%d", pid)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}

	return t.Format(time.DateTime)
}

// formatExitCode only means something once the process is gone.
func formatExitCode(s supervisor.ServiceInfo) string {
	if s.StoppedAt.IsZero() {
		return "-"
	}

	return fmt.Sprintf("%d", s.ExitCode)
}

// formatAge reports how long a service has been up, or how long it ran
// before it stopped.
func formatAge(s supervisor.ServiceInfo) string {
	if s.StartedAt.IsZero() {
		return "-"
	}

	end := time.Now()
	if !s.StoppedAt.IsZero() {
		end = s.StoppedAt
	}

	return formatDuration(end.Sub(s.StartedAt))
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
