package cmd

import (
	"github.com/spf13/cobra"
	"golin/clientinfo"
)

// netstatCmd represents the netstat command
var netstatCmd = &cobra.Command{
	Use:   "netstat",
	Short: "获取CPU进程信息",
	Long:  `降序获取CPU进程信息,输出:PID、PATH、`,
	Run:   clientinfo.Netstat,
}

func init() {
	rootCmd.AddCommand(netstatCmd)
}
