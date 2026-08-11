package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uinit/internal/client"
	"github.com/uinit/internal/config"
)

var logsCmd = &cobra.Command{
	Use:   "logs <process>",
	Short: "Show the logs of a process",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return logsProcess(args[0])
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}

func logsProcess(processName string) error {
	cli, err := client.NewClient(config.GetSockFile())
	if err != nil {
		return err
	}

	logs, err := cli.Logs(processName)
	if err != nil {
		return err
	}

	printLogs(logs)
	return nil
}

func printLogs(logs string) {
	fmt.Println(logs)
}
