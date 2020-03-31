package cmd

import (
	"log"
	"os"

	"github.com/abhirockzz/redis-todo-cli-go/db"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{Use: "list", Short: "list all todos", Run: List}

func init() {
	listCmd.Flags().String("status", "", "completed, pending or in-progress")
	rootCmd.AddCommand(listCmd)
}

// List - todo list --status <completed or pending>
func List(cmd *cobra.Command, args []string) {
	status := cmd.Flag("status").Value.String()

	var todos []db.Todo
	if status == "completed" || status == "pending" || status == "in-progress" || status == "" {
		todos = db.ListTodos(status)
	} else {
		log.Fatalf("provide valid status - completed, pending or in-progress")
	}

	todoTable := [][]string{}

	for _, todo := range todos {
		todoTable = append(todoTable, []string{todo.ID, todo.Desc, todo.Status})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Description", "Status"})

	for _, v := range todoTable {
		table.Append(v)
	}
	table.Render()
}
