package cmd

import (
	"fmt"
	"os"

	"github.com/auXiaoYuan/helios-cli/cmd/ai"
	"github.com/auXiaoYuan/helios-cli/cmd/backend"
	"github.com/auXiaoYuan/helios-cli/cmd/finance"
	"github.com/spf13/cobra"
)

// 版本信息由 goreleaser 通过 ldflags 注入；本地 `go build` 时保留默认值。
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:     "helios-cli",
	Short:   "helios-cli is a unified CLI for backend, ai and finance domains",
	Version: fmt.Sprintf("%s (commit %s, built %s)", version, commit, date),
	Long: `helios-cli aggregates commands from three domains:
  - backend: mocktest, devcode
  - ai:      sealtoken, deepresearch
  - finance: buy, sell`,
}

// NewRootCmd builds a fresh root command tree. It is exposed for testing.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "helios-cli",
		Short:   "helios-cli is a unified CLI for backend, ai and finance domains",
		Version: fmt.Sprintf("%s (commit %s, built %s)", version, commit, date),
	}
	root.AddCommand(backend.NewBackendCmd())
	root.AddCommand(ai.NewAICmd())
	root.AddCommand(finance.NewFinanceCmd())
	return root
}

// Execute runs the root command.
func Execute() {
	// 启动时异步去 npm registry 查最新版本；不阻塞主命令。
	updateCh := startUpdateCheck(version)

	rootCmd.AddCommand(backend.NewBackendCmd())
	rootCmd.AddCommand(ai.NewAICmd())
	rootCmd.AddCommand(finance.NewFinanceCmd())
	err := rootCmd.Execute()

	// 命令跑完再尝试输出更新提示（最多等很短时间，拿不到就跳过）。
	reportUpdate(updateCh, version)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
