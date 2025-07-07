package command

import (
	"fmt"
	"os"

	"github.com/epistax1s/gomer/internal/cli"
	"github.com/epistax1s/gomer/internal/model"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var UserCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"u"},
	Short:   "Manage users",
	Long:    `Manage user roles and permissions.`,
}

var listUsersCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List all users",
	Long:    `List all users in the system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli := cmd.Context().Value("cli").(*cli.CLI)

		users, err := cli.UserService.FindAll()
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"ID", "ChatID", "Name", "Username", "Department", "Role", "Status", "EE"})
		for _, user := range users {
			ee := ""
			if user.EE {
				ee = "+"
			}
			table.Append([]string{
				fmt.Sprintf("%d", user.ID),
				fmt.Sprintf("%d", user.ChatID),
				user.Name,
				user.Username,
				user.Department.Name,
				user.Role,
				user.Status,
				ee,
			})
		}
		table.Render()
		return nil
	},
}

var setRoleCmd = &cobra.Command{
	Use:     "set-role [id] [role]",
	Aliases: []string{"r"},
	Short:   "Set user role",
	Long:    `Set the role for a specific user (ADMIN or USER).`,
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cli := cmd.Context().Value("cli").(*cli.CLI)

		idStr := args[0]
		role := args[1]

		if role != "ADMIN" && role != "USER" {
			return fmt.Errorf("invalid role: %s. Allowed values: ADMIN, USER", role)
		}

		var id int64
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil {
			return fmt.Errorf("invalid id: %s", idStr)
		}

		user, _ := cli.UserService.FindByID(id)
		if user == nil {
			return fmt.Errorf("failed to find user: %w", err)
		}

		user.Role = role
		err = cli.UserService.Save(user)
		if err != nil {
			return fmt.Errorf("failed to update user role: %w", err)
		}

		fmt.Printf("Successfully updated role to '%s' for user with id: %d\n", role, id)
		return nil
	},
}

var blockUserCmd = &cobra.Command{
	Use:     "block [id]",
	Aliases: []string{"b"},
	Short:   "Block user (set status to deleted)",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cli := cmd.Context().Value("cli").(*cli.CLI)

		idStr := args[0]
		var id int64
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil {
			return fmt.Errorf("invalid id: %s", idStr)
		}

		user, _ := cli.UserService.FindByID(id)
		if user == nil {
			return fmt.Errorf("user not found: %d", id)
		}

		user.Status = model.UserStatusDeleted
		err = cli.UserService.Save(user)
		if err != nil {
			return fmt.Errorf("failed to block user: %w", err)
		}

		fmt.Printf("User with id %d blocked (status set to 'deleted')\n", id)
		return nil
	},
}

var unblockUserCmd = &cobra.Command{
	Use:     "unblock [id]",
	Aliases: []string{"u"},
	Short:   "Unblock user (set status to limbo if currently deleted)",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cli := cmd.Context().Value("cli").(*cli.CLI)

		idStr := args[0]
		var id int64
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil {
			return fmt.Errorf("invalid id: %s", idStr)
		}

		user, _ := cli.UserService.FindByID(id)
		if user == nil {
			return fmt.Errorf("user not found: %d", id)
		}

		if user.Status != model.UserStatusDeleted {
			return fmt.Errorf("user is not in 'deleted' status, cannot unblock")
		}

		user.Status = model.UserStatusLimbo
		err = cli.UserService.Save(user)
		if err != nil {
			return fmt.Errorf("failed to unblock user: %w", err)
		}

		fmt.Printf("User with id %d unblocked (status set to 'limbo')\n", id)
		return nil
	},
}

func init() {
	UserCmd.AddCommand(listUsersCmd)
	UserCmd.AddCommand(setRoleCmd)
	UserCmd.AddCommand(blockUserCmd)
	UserCmd.AddCommand(unblockUserCmd)
}
