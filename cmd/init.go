package cmd

import (
	"github.com/spf13/cobra"
	"github.com/uinit/internal/manager"
)

var initCmd = &cobra.Command{
	Use:   "init <config-file-path.yaml>",
	Short: "Start the configured processes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit(args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(configFilePath string) error {
	m, err := manager.NewManager(configFilePath)
	if err != nil {
		return err
	}

	return m.Run()
}
