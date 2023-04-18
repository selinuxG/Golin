/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"golin/global"
	"golin/run"

	"github.com/spf13/cobra"
)

// routeCmd represents the route command
var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "运行采集网络设备功能",
	Long:  `基于SSH的功能进行采集`,
	Run:   run.Route,
}

func init() {
	rootCmd.AddCommand(routeCmd)
	routeCmd.Flags().StringP("ip", "i", global.CmdRoutepath, "此参数是指定待远程采集的IP文件位置")
	routeCmd.Flags().BoolP("echo", "e", false, "此参数是控制控制台是否输出结果,默认不进行输出")
	routeCmd.Flags().StringP("cmd", "c", "", "此参数是指定待自定义执行的命令文件")
	routeCmd.Flags().StringP("spript", "s", global.Split, "此参数是指定IP文件中的分隔字符")
	routeCmd.Flags().StringP("cmdvalue", "C", "", "此参数是自定义执行命令（比-c优先级高）")
	routeCmd.Flags().BoolP("python", "p", false, "此参数是指定python位置，绝对路径，如'D:\\python3\\python.exe'")
}
