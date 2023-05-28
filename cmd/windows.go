package cmd

import (
	"github.com/spf13/cobra"
	"golin/run/windows"
	"runtime"
)

// windowsCmd represents the execl command
var windowsCmd = &cobra.Command{
	Use:   "windows",
	Short: "读取安全策略生成html",
	Long:  `读取安全策略生成html`,
	Run: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS == "windows" {
			windows.Windows()
		}
	},
}

func init() {
	rootCmd.AddCommand(windowsCmd)
}
