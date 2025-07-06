package main

import (
	"context"
	"fmt"
	"os"

	"github.com/epistax1s/gomer/internal/cli"
	"github.com/epistax1s/gomer/internal/cli/command"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gocli",
	Short: "Gomer CLI - management tool for Gomer",
	Long:  `A CLI tool for managing Gomer users, invitations, and roles.`,
}

func init() {
	rootCmd.AddCommand(command.UserCmd)
	rootCmd.AddCommand(command.InvitationCmd)
}

func main() {
	cli := cli.NewCLI()
	ctx := context.WithValue(context.Background(), "cli", cli)
	
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
