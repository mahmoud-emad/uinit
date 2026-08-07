package main

import (
	"fmt"
	"os"

	"github.com/miniinit/internal/config"
)

func printUsage() {
	fmt.Println("Usage: miniinit services.yaml")
	return
}
func main() {
	args := os.Args
	if len(args) != 2 {
		printUsage()
		return
	}

	filePath := args[1]
	config, err := config.NewMiniInit(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", config)
}
