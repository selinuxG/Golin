/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"golin/global"
	"golin/run"

	"github.com/spf13/cobra"
)

// redisCmd represents the redis command
var redisCmd = &cobra.Command{
	Use:   "redis",
	Short: "运行采集Redis功能",
	Long:  `基于Redis的远程登陆功能,通过多线程进行采集,基于info字段中的值判断,写入待采集文件主机时用户名为空即可。`,
	Run:   run.Redis,
}

func init() {
	rootCmd.AddCommand(redisCmd)
	redisCmd.Flags().StringP("ip", "i", global.CmdRedisPath, "此参数是指定待远程采集的IP文件位置")
	redisCmd.Flags().StringP("spript", "s", global.Split, "此参数是指定IP文件中的分隔字符")
	redisCmd.Flags().StringP("value", "v", "", "此参数是指定执行单个设备")

}
