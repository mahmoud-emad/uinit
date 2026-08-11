package cmd

import (
	"github.com/spf13/cobra"
	"github.com/uinit/internal/client"
	"github.com/uinit/internal/config"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <process>",
	Short: "Show the process details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return inspectProcess(args[0])
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

func inspectProcess(processName string) error {
	cli, err := client.NewClient(config.GetSockFile())
	if err != nil {
		return err
	}

	proc, err := cli.Inspect(processName)
	if err != nil {
		return err
	}

	printInspect(proc)
	return nil
}
