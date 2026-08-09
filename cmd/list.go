package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/uinit/internal/supervisor"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured services",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listServices()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listServices() error {
	client, err := supervisor.NewClient()
	if err != nil {
		return err
	}

	services := client.List()
	log.Printf("Services: %v", services)
	return nil
}
