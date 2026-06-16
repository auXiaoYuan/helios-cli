package finance

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewSellCmd builds the `finance sell` subcommand.
func NewSellCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sell",
		Short: "Place a sell order",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "[finance.sell] placing sell order...")
			return err
		},
	}
}
