package cmd

import (
	"github.com/spf13/cobra"
	"golin/run/xshell"
)

var XshellCmd = &cobra.Command{
	Use:   "xshell",
	Short: "获取本机xshell配置信息",
	Long:  `基于Software\PremiumSoft注册表读取账号密码等信息`,
	Run:   xshell.Run,
}

func init() {
	rootCmd.AddCommand(XshellCmd)
}
