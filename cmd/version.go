package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// 这些变量在编译时通过 -ldflags 注入
var (
	Version   = "dev"     // tag 名或分支名
	Commit    = "unknown" // git commit hash (短)
	BuildTime = "unknown" // 构建时间
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func printVersion() {
	fmt.Printf("kubeclean %s (commit: %s, built: %s)\n", Version, Commit, BuildTime)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
