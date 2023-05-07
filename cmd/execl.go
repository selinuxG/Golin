package cmd

import (
	"github.com/spf13/cobra"
	"golin/run"
)

// execlCmd represents the execl command
var execlCmd = &cobra.Command{
	Use:   "execl",
	Short: "golin读取xlsx文件",
	Long:  `通过读取xlsx文件生成golin可读取的格式文件`,
	Run:   run.Execl,
}

func init() {
	rootCmd.AddCommand(execlCmd)
	execlCmd.Flags().StringP("file", "f", "", "此参数是指定读取的文件")
	execlCmd.Flags().StringP("name", "n", "", "此参数是指定名称代表的列")
	execlCmd.Flags().StringP("ip", "i", "", "此参数是指定ip代表的列")
	execlCmd.Flags().StringP("port", "P", "", "此参数是指定端口代表的列")
	execlCmd.Flags().StringP("user", "u", "", "此参数是指定用户所代表的列")
	execlCmd.Flags().StringP("passwd", "p", "", "此参数是指定密码所代表的列")
	execlCmd.Flags().StringP("sheet", "s", "Sheet1", "此参数是指定sheet名称")
	execlCmd.Flags().StringP("sava", "o", "linux_xlsx.txt", "此参数是指定保存的文件")
}
