package cmd

import (
	"github.com/spf13/cobra"
	"golin/port"
)

// portCmd represents the port command
var portCmd = &cobra.Command{
	Use:   "port",
	Short: "端口扫描",
	Run:   port.ParseFlags,
}

func init() {
	rootCmd.AddCommand(portCmd)
	portCmd.Flags().StringP("ip", "i", "", "此参数是扫描的IP地址,格式支持192.168.1.1,192.168.1.1/24,192.168.1-10")
	portCmd.Flags().StringP("port", "p", "0", "此参数是指定的端口,不支持则默认端口,格式支持1,2,3,2-20")
	portCmd.Flags().StringP("exclude", "e", "", "此参数排除扫描的端口,格式支持:1,2,3")
	portCmd.Flags().Bool("noping", false, "此参数是禁止ping检测")
	portCmd.Flags().IntP("chan", "c", 30, "并发数量")
	portCmd.Flags().IntP("time", "t", 3, "超时等待时常/s")
	portCmd.Flags().Bool("random", false, "打乱主机顺序")
	portCmd.Flags().Bool("save", false, "保存结果")

}
