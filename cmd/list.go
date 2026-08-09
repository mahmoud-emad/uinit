package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured services",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return listConfiguredServices(args[0])
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listConfiguredServices(filepath string) error {
	// We should call the uinit daemon to list all services
	// cfg, err := config.NewMiniInit(filepath)
	// if err != nil {
	// 	return err
	// }

	// printLoadedServices(cfg)

	return nil
}
