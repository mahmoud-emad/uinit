package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "uinit",
	Short: "uinit is a small process supervisor written in Go",
	Long:  "A fast and flexible process supervisor built in Go.",
	// Let Execute report the error, instead of cobra printing it again
	// along with the whole usage text.
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, red.Render("Error:"), err)
		os.Exit(1)
	}
}
