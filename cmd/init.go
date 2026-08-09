package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/uinit/internal/supervisor"
)

var initCmd = &cobra.Command{
	Use:   "init <config>",
	Short: "Start the configured services",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return initServices(args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initServices(filepath string) error {
	supervisor, err := supervisor.NewSupervisor(filepath)
	if err != nil {
		log.Fatal(err)
	}

	if err := supervisor.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}
