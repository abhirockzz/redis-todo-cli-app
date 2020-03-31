package cmd

import (
	"github.com/abhirockzz/redis-todo-cli-go/db"
	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{Use: "delete", Short: "delete a todo", Run: Delete}

func init() {
	delCmd.Flags().String("id", "", "id for the todo you want to delete")
	delCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(delCmd)
}

// Delete - todo delete --id <id>
func Delete(cmd *cobra.Command, args []string) {
	id := cmd.Flag("id").Value.String()
	db.DeleteTodo(id)
}
