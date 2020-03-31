package cmd

import (
	"github.com/abhirockzz/redis-todo-cli-go/db"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{Use: "create", Short: "create a todo with description", Run: Create}

func init() {
	createCmd.Flags().String("description", "", "create todo with description")
	createCmd.MarkFlagRequired("description")
	rootCmd.AddCommand(createCmd)
}

// Create - todo create --description <text>
func Create(cmd *cobra.Command, args []string) {
	desc := cmd.Flag("description").Value.String()
	db.CreateTodo(desc)
}
