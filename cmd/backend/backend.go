package backend

import "github.com/spf13/cobra"

// NewBackendCmd builds the backend domain command with its subcommands.
func NewBackendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backend",
		Short: "Backend domain commands",
		Long:  "Commands related to the backend domain, including mocktest and devcode.",
	}
	cmd.AddCommand(NewMockTestCmd())
	cmd.AddCommand(NewDevCodeCmd())
	return cmd
}
