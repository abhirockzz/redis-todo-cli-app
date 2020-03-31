package cmd

import (
	"log"

	"github.com/abhirockzz/redis-todo-cli-go/db"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{Use: "update", Short: "update todo description, status or both", Run: Update}

func init() {
	updateCmd.Flags().String("id", "", "id of the todo you want to update")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().String("description", "", "new description")
	updateCmd.Flags().String("status", "", "new status: completed, pending, in-progress")

	rootCmd.AddCommand(updateCmd)
}

// Update - todo update --id <id> --status <new status> --description <new description>
func Update(cmd *cobra.Command, args []string) {
	id := cmd.Flag("id").Value.String()
	desc := cmd.Flag("description").Value.String()

	status := cmd.Flag("status").Value.String()

	if desc == "" && status == "" {
		log.Fatalf("either description or status is required")
	}

	if status == "completed" || status == "pending" || status == "in-progress" || status == "" {
		db.UpdateTodo(id, desc, status)
	} else {
		log.Fatalf("provide valid status - completed, pending or in-progress")
	}
}
