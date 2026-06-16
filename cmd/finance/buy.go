package finance

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewBuyCmd builds the `finance buy` subcommand.
func NewBuyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "buy",
		Short: "Place a buy order",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "[finance.buy] placing buy order...")
			return err
		},
	}
}
