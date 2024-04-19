package cmd

import (
	"github.com/spf13/cobra"
	"golin/clientinfo"
)

// cpuCmd represents the cpu command
var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "获取CPU进程信息",
	Long:  `降序获取CPU进程信息,输出:PID、PATH、`,
	Run:   clientinfo.CPU,
}

func init() {
	rootCmd.AddCommand(cpuCmd)
	cpuCmd.Flags().IntP("count", "c", 5, "降序输出多少个")
}
