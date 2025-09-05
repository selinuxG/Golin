//go:build darwin

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// wifiCmd represents the wifi command for macOS
var wifiCmd = &cobra.Command{
	Use:   "wifi",
	Short: "WiFi 密码获取 (macOS 不支持)",
	Long:  `WiFi 密码获取功能在 macOS 上不可用`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("WiFi 密码获取功能在 macOS 上不可用")
	},
}

func init() {
	rootCmd.AddCommand(wifiCmd)
}
