//go:build darwin

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// windowsCmd represents the windows command for macOS
var windowsCmd = &cobra.Command{
	Use:   "windows",
	Short: "Windows 安全策略检查 (macOS 不支持)",
	Long:  `Windows 安全策略检查功能在 macOS 上不可用`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Windows 安全策略检查功能在 macOS 上不可用")
	},
}

func init() {
	rootCmd.AddCommand(windowsCmd)
}
