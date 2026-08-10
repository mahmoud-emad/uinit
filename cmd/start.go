package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uinit/internal/client"
	"github.com/uinit/internal/manager"
)

var startCmd = &cobra.Command{
	Use:   "start <process>",
	Short: "Start a loaded process",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return startProcess(args[0])
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startProcess(processName string) error {
	cli, err := client.NewClient(manager.SocketPath)
	if err != nil {
		return err
	}

	rsp, err := cli.Start(processName)
	if err != nil {
		return err
	}

	if !rsp.OK {
		return fmt.Errorf("response error %s", rsp.Message)
	}

	return printStatus(rsp.Data)
}
