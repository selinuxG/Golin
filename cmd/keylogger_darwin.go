//go:build darwin

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// keyloggerCmd represents the keylogger command for macOS
var keyloggerCmd = &cobra.Command{
	Use:   "keylogger",
	Short: "键盘记录器 (macOS 不支持)",
	Long:  `键盘记录器功能在 macOS 上不可用`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("键盘记录器功能在 macOS 上不可用")
	},
}

func init() {
	rootCmd.AddCommand(keyloggerCmd)
}
