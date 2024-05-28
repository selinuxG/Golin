//go:build linux

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// redisCmd represents the redis command
var keylogger = &cobra.Command{
	Use:   "keylogger",
	Short: "键盘记录器",
	Long:  `实时记录键盘输入信息到日志文件中`,
	Run:   KeyLoggerCmd,
}

func init() {
	rootCmd.AddCommand(keylogger)
}

func KeyLoggerCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("[-] 仅支持Windows操作系统\n")
	return
}
