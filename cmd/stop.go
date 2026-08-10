package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uinit/internal/client"
	"github.com/uinit/internal/manager"
)

var stopCmd = &cobra.Command{
	Use:   "stop <process>",
	Short: "stop a loaded process",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return stopProcess(args[0])
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stopProcess(processName string) error {
	cli, err := client.NewClient(manager.SocketPath)
	if err != nil {
		return err
	}

	rsp, err := cli.Stop(processName)
	if err != nil {
		return err
	}

	if !rsp.OK {
		return fmt.Errorf("response error %s", rsp.Message)
	}

	return printStatus(rsp.Data)
}
