package cmd

import (
	"github.com/spf13/cobra"
	"github.com/uinit/internal/client"
	"github.com/uinit/internal/config"
)

var startCmd = &cobra.Command{
	Use:   "start <process>",
	Short: "Start the loaded process",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return startProcess(args[0])
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startProcess(processName string) error {
	cli, err := client.NewClient(config.GetSockFile())
	if err != nil {
		return err
	}

	proc, err := cli.Start(processName)
	if err != nil {
		return err
	}

	printInspect(proc)
	return nil
}
