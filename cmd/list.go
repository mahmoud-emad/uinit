package cmd

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"github.com/uinit/internal/supervisor"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured services",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listServices()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listServices() error {
	client, err := supervisor.NewClient()
	if err != nil {
		return err
	}

	rsp, err := client.List()
	if err != nil {
		return err
	}

	if !rsp.OK {
		return fmt.Errorf("response error %s", rsp.Message)
	}

	printServices(rsp.Data)

	return nil
}

// Fixed widths for the columns that never grow, plus the cap on how much
// of a command line is worth showing.
const (
	statusWidth = 10
	pidWidth    = 7
	uptimeWidth = 7
	maxCmdWidth = 40
	gutter      = "  "
)

func printServices(services []supervisor.ServiceInfo) {
	if len(services) == 0 {
		fmt.Println(style("No services loaded.", dim))
		return
	}

	nameWidth := len("SERVICE")
	cmdWidth := len("COMMAND")

	for _, service := range services {
		nameWidth = max(nameWidth, utf8.RuneCountInString(service.Name))
		cmdWidth = max(cmdWidth, utf8.RuneCountInString(service.Cmd))
	}

	cmdWidth = min(cmdWidth, maxCmdWidth)

	header := row(
		pad("SERVICE", nameWidth),
		pad("STATUS", statusWidth),
		pad("PID", pidWidth),
		pad("UPTIME", uptimeWidth),
		"COMMAND",
	)

	width := nameWidth + statusWidth + pidWidth + uptimeWidth + cmdWidth + 5*len(gutter)

	fmt.Println()
	fmt.Println(style(header, bold))
	fmt.Println(style(strings.Repeat("─", width), dim))

	running := 0

	for _, service := range services {
		if service.Status == supervisor.Running || service.Status == supervisor.Started {
			running++
		}

		fmt.Println(row(
			pad(service.Name, nameWidth),
			formatStatus(service.Status, statusWidth),
			style(pad(formatPID(service.PID), pidWidth), dim),
			pad(formatAge(service), uptimeWidth),
			style(truncate(service.Cmd, cmdWidth), dim),
		))
	}

	fmt.Println()
	fmt.Println(style(fmt.Sprintf("%s%d services · %d running", gutter, len(services), running), dim))
	fmt.Println()
}

// row indents the line and separates each cell by the gutter.
func row(cells ...string) string {
	return gutter + strings.Join(cells, gutter)
}
