package cmd

import (
	"github.com/spf13/cobra"
	"golin/port"
)

// portCmd represents the port command
var portCmd = &cobra.Command{
	Use:   "port",
	Short: "资产测绘、协议识别、漏洞扫描",
	Run:   port.ParseFlags,
}

func init() {
	rootCmd.AddCommand(portCmd)
	portCmd.Flags().StringP("ip", "i", "", "此参数是扫描的IP地址,格式支持192.168.1.1,192.168.1.1/24,192.168.1-10")
	portCmd.Flags().StringP("ipfile", "", "ip.txt", "此参数是扫描的IP文件,一行一个")
	portCmd.Flags().StringP("port", "p", "0", "此参数是指定的端口,不支持则默认端口,格式支持1,2,3,2-20")
	portCmd.Flags().StringP("exclude", "e", "", "此参数排除扫描的端口,格式支持:1,2,3")
	portCmd.Flags().StringP("excludeip", "", "noip.txt", "此参数排除扫描的IP")
	portCmd.Flags().Bool("noping", false, "此参数是禁止ping检测")
	portCmd.Flags().IntP("chan", "c", 100, "并发数量")
	portCmd.Flags().IntP("time", "t", 5, "超时等待时常/s")
	portCmd.Flags().Bool("random", false, "打乱主机顺序")
	portCmd.Flags().Bool("nocrack", false, "此参数是不进行弱口令扫描")
	portCmd.Flags().Bool("nopoc", false, "此参数是不进行poc漏洞扫描")
	portCmd.Flags().StringP("userfile", "", "", "此参数是自定义用户字典文件")
	portCmd.Flags().StringP("passwdfile", "", "", "此参数是自定义密码字典文件")
	portCmd.Flags().StringP("fofa", "", "", "此参数是调用fofa数据进行扫描")
}
