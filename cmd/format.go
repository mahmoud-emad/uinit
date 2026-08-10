package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/uinit/internal/config"
	"github.com/uinit/internal/manager"
	"github.com/uinit/internal/process"
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
	printField("Processes", fmt.Sprintf("%d", len(cfg.Processes)))
	fmt.Println()

	width := 0
	for _, p := range cfg.Processes {
		width = max(width, lipgloss.Width(p.Name))
	}

	name := lipgloss.NewStyle().Width(width + 2)
	for _, p := range cfg.Processes {
		fmt.Println(name.Render(p.Name) + dim.Render(p.Cmd))
	}

	fmt.Println()
	fmt.Println(dim.Render("Press Ctrl+C to stop."))
	fmt.Println()
}

// printList prints one or more processes as a table.
func printList(processes []manager.ProcessInfo) {
	if len(processes) == 0 {
		fmt.Println(dim.Render("No processes loaded."))
		return
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(false).BorderBottom(false).BorderLeft(false).BorderRight(false).
		BorderColumn(false).
		BorderStyle(dim).
		Headers("PROCESS", "STATUS", "PID", "UPTIME", "COMMAND").
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == table.HeaderRow {
				return cell.Bold(true)
			}
			return cell
		})

	running := 0

	for _, p := range processes {
		if p.Status == process.Running {
			running++
		}

		t.Row(p.Name, formatStatus(p.Status), formatPID(p.PID), formatAge(p), dim.Render(p.Cmd))
	}

	fmt.Println()
	fmt.Println(t)
	fmt.Println()
	fmt.Println(dim.Render(fmt.Sprintf("%d processes · %d running", len(processes), running)))
	fmt.Println()
}

// printLogs dumps what a process has written so far. The daemon only reports
// where the log file is, the reading happens here.
func printLogs(processes []manager.ProcessInfo) error {
	if len(processes) == 0 {
		fmt.Println(dim.Render("No processes found."))
		return nil
	}

	p := processes[0]

	if p.LogFile == "" {
		fmt.Println(dim.Render("No log file for " + p.Name + "."))
		return nil
	}

	logs, err := os.ReadFile(p.LogFile)
	if err != nil {
		return err
	}

	if len(logs) == 0 {
		fmt.Println(dim.Render("No logs yet for " + p.Name + "."))
		return nil
	}

	fmt.Print(string(logs))
	return nil
}

// printInspect prints one or more processes in full, one block each.
func printInspect(processes []manager.ProcessInfo) {
	if len(processes) == 0 {
		fmt.Println(dim.Render("No processes found."))
		return
	}

	for _, p := range processes {
		fmt.Println()
		fmt.Println("Process: " + bold.Render(p.Name))
		fmt.Println(rule())

		printField("Status", formatStatus(p.Status))
		printField("PID", formatPID(p.PID))
		printField("Command", p.Cmd)

		fmt.Println()
		printField("Started", formatTime(p.StartedAt))
		printField("Stopped", formatTime(p.StoppedAt))
		printField("Exit code", formatExitCode(p))
		printField("Duration", formatAge(p))

		// Error goes here once ProcessInfo carries it.
		printField("Logs", bold.Render(p.LogFile))

		fmt.Println()
	}
}

func printField(name, value string) {
	fmt.Println(label.Render(name+":") + value)
}

func rule() string {
	return dim.Render(strings.Repeat("─", ruleWidth))
}

func formatStatus(status process.Status) string {
	switch status {
	case process.Running:
		return green.Render("● " + status.String())
	case process.Starting:
		return yellow.Render("◌ " + status.String())
	case process.Failed:
		return red.Render("✕ " + status.String())
	case process.Exited:
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
func formatExitCode(p manager.ProcessInfo) string {
	if p.StoppedAt.IsZero() {
		return "-"
	}

	return fmt.Sprintf("%d", p.ExitCode)
}

// formatAge reports how long a process has been up, or how long it ran
// before it stopped.
func formatAge(p manager.ProcessInfo) string {
	if p.StartedAt.IsZero() {
		return "-"
	}

	end := time.Now()
	if !p.StoppedAt.IsZero() {
		end = p.StoppedAt
	}

	return formatDuration(end.Sub(p.StartedAt))
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
