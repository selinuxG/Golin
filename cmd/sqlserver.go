package cmd

import (
	"github.com/spf13/cobra"
	"golin/global"
	"golin/run"
)

// linuxCmd represents the linux command
var sqlServerCmd = &cobra.Command{
	Use:   "sqlserver",
	Short: "运行采集sqlserver采集模式",
	Long:  `基于远程登录功能,通过多线程的方法批量进行采集`,
	Run:   run.SqlServer,
}

func init() {
	rootCmd.AddCommand(sqlServerCmd)
	sqlServerCmd.Flags().StringP("ip", "i", global.CmdsqlServerPath, "此参数是指定待远程采集的IP文件位置")
	sqlServerCmd.Flags().StringP("spript", "s", global.Split, "此参数是指定IP文件中的分隔字符")
	sqlServerCmd.Flags().StringP("value", "v", "", "此参数是单次执行")
	sqlServerCmd.Flags().BoolP("echo", "e", false, "此参数是控制控制台是否输出结果,默认不进行输出")
}
