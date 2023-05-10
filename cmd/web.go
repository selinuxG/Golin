package cmd

import (
	"github.com/spf13/cobra"
	"golin/web"
)

// execlCmd represents the execl command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "web形式运行",
	Long:  `通过gin实现资产任务下发`,
	Run:   web.Start,
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().StringP("port", "p", "1818", "指定运行端口")
	webCmd.Flags().StringP("ip", "i", "127.0.0.1", "指定运行网卡ip")
	webCmd.Flags().BoolP("save", "s", false, "是否额外保存文件")

}
