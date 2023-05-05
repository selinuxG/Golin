package cmd

import (
	"github.com/spf13/cobra"
	"golin/global"
	"golin/run"
)

// mysqlCmd represents the mysql command
var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "运行采集MySql采集功能",
	Long:  `基于Mysql远程通过多线程连接执行指定sql语句并记录,连接等待为10秒左右,连不上则断开。`,
	Run:   run.Mysql,
}

func init() {
	rootCmd.AddCommand(mysqlCmd)
	mysqlCmd.Flags().StringP("ip", "i", global.CmdMysqlPath, "此参数是指定待远程采集的IP文件位置")
	mysqlCmd.Flags().StringP("spript", "s", global.Split, "此参数是指定IP文件中的分隔字符")
	mysqlCmd.Flags().StringP("value", "v", "", "此参数是指定执行单个主机")
	mysqlCmd.Flags().StringP("cmd", "c", "", "此参数是自定义执行sql语句")
	mysqlCmd.Flags().BoolP("echo", "e", false, "此参数指定是控制是否输出结果")
}
