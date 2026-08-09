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
	_, err := config.NewConfig(filepath)
	if err != nil {
		return err
	}

	supervisor.NewSupervisor()
	// printLoadedServices(cfg)

	return nil
}

// func printLoadedServices(cfg config.MiniInit) {
// 	fmt.Printf("\nLoaded %d services:\n", len(cfg.Services))

// 	for _, service := range cfg.Services {
// 		fmt.Printf("- %s:\n", service.Name)
// 		fmt.Printf("  command: %s\n\n", service.Cmd)
// 	}
// }
