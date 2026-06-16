package finance

import "github.com/spf13/cobra"

// NewFinanceCmd builds the finance domain command with its subcommands.
func NewFinanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finance",
		Short: "Finance domain commands",
		Long:  "Commands related to the finance domain, including buy and sell.",
	}
	cmd.AddCommand(NewBuyCmd())
	cmd.AddCommand(NewSellCmd())
	return cmd
}
