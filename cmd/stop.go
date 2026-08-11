package cmd

import (
	"github.com/spf13/cobra"
	"github.com/uinit/internal/client"
	"github.com/uinit/internal/config"
)

var stopCmd = &cobra.Command{
	Use:   "stop <process>",
	Short: "stop the loaded process",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return stopProcess(args[0])
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stopProcess(processName string) error {
	cli, err := client.NewClient(config.GetSockFile())
	if err != nil {
		return err
	}

	proc, err := cli.Stop(processName)
	if err != nil {
		return err
	}

	printInspect(proc)
	return nil
}
