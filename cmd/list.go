package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uinit/internal/client"
	"github.com/uinit/internal/manager"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured processes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listProcesses()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listProcesses() error {
	cli, err := client.NewClient(manager.SocketPath)
	if err != nil {
		return err
	}

	rsp, err := cli.List()
	if err != nil {
		return err
	}

	if !rsp.OK {
		return fmt.Errorf("response error %s", rsp.Message)
	}

	printList(rsp.Data)
	return nil
}
