/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"golin/global"
	"golin/run"
)

// linuxCmd represents the linux command
var linuxCmd = &cobra.Command{
	Use:   "linux",
	Short: "运行采集Linux采集模式",
	Long:  `基于SSH协议远程登录功能,通过多线程的方法批量进行采集`,
	Run:   run.Linux,
}

func init() {
	rootCmd.AddCommand(linuxCmd)
	linuxCmd.Flags().StringP("ip", "i", global.CmdLinuxPath, "此参数是指定待远程采集的IP文件位置")
	linuxCmd.Flags().StringP("cmd", "c", "", "此参数是指定待自定义执行的命令文件")
	linuxCmd.Flags().StringP("spript", "s", global.Split, "此参数是指定IP文件中的分隔字符")
	linuxCmd.Flags().StringP("value", "v", "", "此参数是指定执行单个主机")
	linuxCmd.Flags().StringP("cmdvalue", "C", "", "此参数是自定义执行命令（比-c优先级高）")
	linuxCmd.Flags().BoolP("echo", "e", false, "此参数是控制控制台是否输出结果,默认不进行输出")
	linuxCmd.Flags().BoolP("localhost", "l", false, "此参数是控制本机采集的模式")
}
