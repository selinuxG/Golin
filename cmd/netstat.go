package cmd

import (
	"github.com/spf13/cobra"
	"golin/clientinfo"
)

// netstatCmd represents the netstat command
var netstatCmd = &cobra.Command{
	Use:   "netstat",
	Short: "获取网络连接信息",
	Long:  `降序获取网络连接信息 比如：pid、local_address、remote_address、status、program_name`,
	Run:   clientinfo.Netstat,
}

func init() {
	rootCmd.AddCommand(netstatCmd)
}
