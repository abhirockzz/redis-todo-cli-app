package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "todo", Short: "Yet another TODO app. Uses Go and Redis", Version: "0.1.0"}

// Execute is the entry point
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("cannot start todo app - ", err)
	}
}
