/*
Copyright © 2023 NAME HERE <selinuxg@163.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"golin/network"
)

// networkCmd represents the network command

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "运行网络相关功能,目前仅有syslog模拟器",
	Long:  "",
	Run:   network.Networkrun,
}

func init() {
	rootCmd.AddCommand(networkCmd)
	networkCmd.Flags().BoolP("syslog", "s", false, "模拟syslog接收端服务器")
}
