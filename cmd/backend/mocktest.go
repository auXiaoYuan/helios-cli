package backend

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewMockTestCmd builds the `backend mocktest` subcommand.
func NewMockTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mocktest",
		Short: "Run mock tests for backend services",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "[backend.mocktest] running mock tests...")
			return err
		},
	}
}
