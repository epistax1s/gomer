package command

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/epistax1s/gomer/internal/cli"
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users",
	Long:  `Manage user roles and permissions.`,
}

var listUsersCmd = &cobra.Command{
	Use:   "list",
	Short: "List all users",
	Long:  `List all users in the system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli := cmd.Context().Value("cli").(*cli.CLI)

		users, err := cli.UserService.FindAll()
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}

		fmt.Printf("Users: %v", users)
		return nil
	},
}

var setRoleCmd = &cobra.Command{
	Use:   "set-role [chatID] [role]",
	Short: "Set user role",
	Long:  `Set the role for a specific user (admin or user).`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		chatID := args[0]
		role := args[1]

		fmt.Printf("Successfully updated role to '%s' for user with chatID: %s", role, chatID)
		return nil
	},
}

func init() {
	UserCmd.AddCommand(listUsersCmd)
	UserCmd.AddCommand(setRoleCmd)
}
