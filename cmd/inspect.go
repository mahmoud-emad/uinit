package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uinit/internal/supervisor"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <service>",
	Short: "Inspect service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return inspectServices(args[0])
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

func inspectServices(serviceName string) error {
	client, err := supervisor.NewClient()
	if err != nil {
		return err
	}

	rsp, err := client.Inspect(serviceName)
	if err != nil {
		return err
	}

	if !rsp.OK {
		return fmt.Errorf("response error %s", rsp.Message)
	}

	printInspect(rsp.Data)
	return nil
}
