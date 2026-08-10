package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uinit/internal/client"
	"github.com/uinit/internal/config"
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
	cli, err := client.NewClient(config.GetSockFile())
	if err != nil {
		return err
	}

	rsp, err := cli.List()
	if err != nil {
		return err
	}

	printList(rsp)
	return nil
}

func printList(data []client.ProcessInfo) {
	fmt.Println(data)
}
