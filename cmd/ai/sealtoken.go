package ai

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewSealTokenCmd builds the `ai sealtoken` subcommand.
func NewSealTokenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sealtoken",
		Short: "Seal an AI access token",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "[ai.sealtoken] sealing AI access token...")
			return err
		},
	}
}
