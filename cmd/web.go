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
	webCmd.Flags().StringP("port", "p", "1818", "此参数是指定运行端口")
	webCmd.Flags().StringP("ip", "i", "127.0.0.1", "此参数是指定运行网卡ip")
	webCmd.Flags().BoolP("save", "s", false, "此参数是指定除去返回的文件是否还固定留存一份,暂时不可用")

}
