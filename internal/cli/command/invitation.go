package command

import (
	"fmt"
	"os"

	"github.com/epistax1s/gomer/internal/cli"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var InvitationCmd = &cobra.Command{
	Use:   "invitation",
	Short: "Manage invitations",
	Long:  `Create and manage invitations for new users.`,
}

var listInvitationsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all invitations",
	Long:  `List all invitations in the system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli := cmd.Context().Value("cli").(*cli.CLI)

		invitations, err := cli.InvitationService.GetInvitesByCreator(1)
		if err != nil {
			return fmt.Errorf("failed to list invitations: %w", err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"Code", "Created At", "Used By", "Used At"})
		for _, invite := range invitations {
			link := fmt.Sprintf("https://t.me/%s?start=%s", cli.Config.Bot.Username, invite.Code)

			usedBy := "N/A"
			usedAt := "N/A"
			if invite.Used {
				usedBy = invite.UsedBy.Name
				usedAt = *invite.UsedAt
			}

			table.Append([]string{
				link,
				invite.CreatedAt,
				usedBy,
				usedAt,
			})
		}

		table.Render()
		return nil
	},
}

var createInvitationCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new invitation",
	Long:  `Create a new invitation code for a user to join the system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli := cmd.Context().Value("cli").(*cli.CLI)

		code, err := cli.InvitationService.GenerateInvite(1)
		if err != nil {
			return fmt.Errorf("failed to create invitation: %w", err)
		}

		fmt.Printf("Successfully created invitation with code: %s", code)

		return nil
	},
}

func init() {
	InvitationCmd.AddCommand(listInvitationsCmd)
	InvitationCmd.AddCommand(createInvitationCmd)
}
