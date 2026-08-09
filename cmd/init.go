package cmd

import (
	"github.com/spf13/cobra"
	"github.com/uinit/internal/config"
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
	// Read the config up front so the banner can show what is about to run.
	cfg, err := config.NewConfig(filepath)
	if err != nil {
		return err
	}

	printStartup(filepath, cfg)

	sup, err := supervisor.NewSupervisor(filepath)
	if err != nil {
		return err
	}

	return sup.Run()
}
