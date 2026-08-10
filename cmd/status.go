package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uinit/internal/client"
	"github.com/uinit/internal/manager"
)

var statusCmd = &cobra.Command{
	Use:   "status <process>",
	Short: "Show the logs of a process",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return statusProcess(args[0])
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func statusProcess(processName string) error {
	cli, err := client.NewClient(manager.SocketPath)
	if err != nil {
		return err
	}

	rsp, err := cli.Status(processName)
	if err != nil {
		return err
	}

	if !rsp.OK {
		return fmt.Errorf("response error %s", rsp.Message)
	}

	return printStatus(rsp.Data)
}
