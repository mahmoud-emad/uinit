package cmd

import (
	"encoding/json"
	"fmt"
	"time"

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

	rsp, err := client.List()
	if err != nil {
		return err
	}

	if !rsp.OK {
		return fmt.Errorf("response error %s", rsp.Message)
	}

	printServices(rsp.Data)

	return nil
}

func printServices(data interface{}) {
	servicesJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("failed to format services:", err)
		return
	}

	var services []supervisor.ManagedService

	if err := json.Unmarshal(servicesJSON, &services); err != nil {
		fmt.Println("failed to parse services:", err)
		return
	}

	if len(services) == 0 {
		fmt.Println("No services loaded.")
		return
	}

	fmt.Printf("%-20s %-12s %-25s\n", "SERVICE", "STATUS", "LOADED AT")
	fmt.Println("------------------------------------------------------------")

	for _, service := range services {
		fmt.Printf(
			"%-20s %-12s %-25s\n",
			service.Config.Name,
			service.Runtime.Status,
			service.Runtime.LoadedAt.Format(time.RFC3339),
		)
	}
}
