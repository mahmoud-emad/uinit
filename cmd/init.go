package cmd

import (
	"github.com/spf13/cobra"
	"github.com/uinit/internal/config"
	"github.com/uinit/internal/manager"
)

var initCmd = &cobra.Command{
	Use:   "init <config>",
	Short: "Start the configured processes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return initProcesses(args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initProcesses(filepath string) error {
	// Read the config up front so the banner can show what is about to run.
	cfg, err := config.NewConfig(filepath)
	if err != nil {
		return err
	}

	printStartup(filepath, cfg)

	pm, err := manager.NewProcessManager(filepath)
	if err != nil {
		return err
	}

	return pm.Run()
}
