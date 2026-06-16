package cmd

import (
	"fmt"
	"os"

	"github.com/auXiaoYuan/helios-cli/cmd/ai"
	"github.com/auXiaoYuan/helios-cli/cmd/backend"
	"github.com/auXiaoYuan/helios-cli/cmd/finance"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "helios-cli",
	Short: "helios-cli is a unified CLI for backend, ai and finance domains",
	Long: `helios-cli aggregates commands from three domains:
  - backend: mocktest, devcode
  - ai:      sealtoken, deepresearch
  - finance: buy, sell`,
}

// NewRootCmd builds a fresh root command tree. It is exposed for testing.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "helios-cli",
		Short: "helios-cli is a unified CLI for backend, ai and finance domains",
	}
	root.AddCommand(backend.NewBackendCmd())
	root.AddCommand(ai.NewAICmd())
	root.AddCommand(finance.NewFinanceCmd())
	return root
}

// Execute runs the root command.
func Execute() {
	rootCmd.AddCommand(backend.NewBackendCmd())
	rootCmd.AddCommand(ai.NewAICmd())
	rootCmd.AddCommand(finance.NewFinanceCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
