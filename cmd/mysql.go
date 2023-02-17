/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"golin/run"
)

// mysqlCmd represents the mysql command
var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "运行采集Myql功能",
	Long:  `基于Mysql远程通过多线程连接执行指定sql语句并记录,连接等待为10秒左右,连不上则断开。`,
	Run:   run.Mysql,
}

func init() {
	rootCmd.AddCommand(mysqlCmd)
	mysqlCmd.Flags().StringP("ip", "i", "mysql.txt", "此参数是指定待远程采集的IP文件位置")
	mysqlCmd.Flags().StringP("spript", "s", run.Split, "此参数是指定IP文件中的分隔字符")
	mysqlCmd.Flags().StringP("value", "v", "", "此参数是指定执行单个主机")

}
