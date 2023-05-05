package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "golin",
	Short: "A brief description of your application",
	Long:  `[-] 此工具是基于Golang多线程的模式开发，目的是批量执行各种设备的命令并记录;但一定要确保分隔字符必须是5个否则会直接跳过此设备。默认是~分隔,1:名称 2:连接地址 3:连接用户 4:连接密码 5:连接端口`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	cmd.Help()
	//},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
