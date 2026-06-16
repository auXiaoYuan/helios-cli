package backend

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewDevCodeCmd builds the `backend devcode` subcommand.
func NewDevCodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "devcode",
		Short: "Generate or manage backend development code",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "[backend.devcode] generating development code...")
			return err
		},
	}
}
