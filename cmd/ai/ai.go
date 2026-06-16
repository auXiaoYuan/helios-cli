package ai

import "github.com/spf13/cobra"

// NewAICmd builds the ai domain command with its subcommands.
func NewAICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI domain commands",
		Long:  "Commands related to the AI domain, including sealtoken and deepresearch.",
	}
	cmd.AddCommand(NewSealTokenCmd())
	cmd.AddCommand(NewDeepResearchCmd())
	return cmd
}
