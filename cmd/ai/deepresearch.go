package ai

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewDeepResearchCmd builds the `ai deepresearch` subcommand.
func NewDeepResearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deepresearch",
		Short: "Run AI deep research workflows",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "[ai.deepresearch] running deep research...")
			return err
		},
	}
}
