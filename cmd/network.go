/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"golin/network"
)

// networkCmd represents the network command

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "运行抓包功能",
	Long: `基于google发布的库进行抓取流量
过滤源 IP 地址为 192.168.1.1 的数据包：src host 192.168.1.1  过滤目标 IP 地址为 192.168.1.1 的数据包：dst host 192.168.1.1
过滤源或目标 IP 地址为 192.168.1.1 的数据包：host 192.168.1.1  过滤源端口为 80 的数据包：src port 80
过滤目标端口为 80 的数据包：dst port 80  过滤 TCP 数据包：tcp  过滤 UDP 数据包：udp  组合多个条件进行过滤：tcp and dst port 80 and src host 192.168.1.1`,
	Run: network.Networkrun,
}

func init() {

	rootCmd.AddCommand(networkCmd)
	networkCmd.Flags().StringP("interface", "i", "\\Device\\NPF_{B3C9B1B3-FFC0-4D59-8667-4B0BF6D354CE}", "指定监听网卡")
	networkCmd.Flags().StringP("bfp", "b", "tcp", "指定监听网卡")
	networkCmd.Flags().BoolP("interfaceall", "a", false, "显示本地所有网卡及IP地址")

}
