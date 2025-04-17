package cmd

import (
	"github.com/spf13/cobra"
	"golin/run/navicat"
)

var NavicatCmd = &cobra.Command{
	Use:   "navicat",
	Short: "获取本机navicat配置信息",
	Long:  `基于Software\PremiumSoft注册表读取账号密码等信息`,
	Run:   navicat.Run,
}

func init() {
	rootCmd.AddCommand(NavicatCmd)
}
