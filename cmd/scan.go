package cmd

import (
	"github.com/spf13/cobra"
	"golin/scan"
)

// ScanCmd represents the port command
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "资产测绘、协议识别、漏洞扫描",
	Run:   scan.ParseFlags,
}

func init() {
	rootCmd.AddCommand(ScanCmd)
	ScanCmd.Flags().StringP("ip", "i", "", "此参数是扫描的IP地址,格式支持192.168.1.1,192.168.1.1/24,192.168.1-10")
	ScanCmd.Flags().StringP("ipfile", "", "ip.txt", "此参数是扫描的IP文件,一行一个")
	ScanCmd.Flags().StringP("port", "p", "0", "此参数是指定的端口,不支持则默认端口,格式支持1,2,3,2-20")
	ScanCmd.Flags().StringP("exclude", "e", "", "此参数排除扫描的端口,格式支持:1,2,3")
	ScanCmd.Flags().StringP("excludeip", "", "noip.txt", "此参数排除扫描的IP")
	ScanCmd.Flags().Bool("noping", false, "此参数是禁止ping检测")
	ScanCmd.Flags().IntP("chan", "c", 100, "并发数量")
	ScanCmd.Flags().IntP("time", "t", 5, "超时等待时常/s")
	ScanCmd.Flags().IntP("done", "", 10, "端口整体扫描最长用时/m")
	ScanCmd.Flags().Bool("random", false, "打乱主机顺序")
	ScanCmd.Flags().Bool("img", false, "此参数进行保存WEB截图")
	ScanCmd.Flags().Bool("nocrack", false, "此参数是不进行弱口令扫描")
	ScanCmd.Flags().Bool("nopoc", false, "此参数是不进行poc漏洞扫描")
	ScanCmd.Flags().StringP("userfile", "", "", "此参数是自定义用户字典文件")
	ScanCmd.Flags().StringP("passwdfile", "", "", "此参数是自定义密码字典文件")
	ScanCmd.Flags().StringP("fofa", "", "", "此参数是调用fofa数据进行扫描")
	ScanCmd.Flags().IntP("fofasize", "", 100, "获取多少条fofa数据")
	ScanCmd.Flags().Bool("web", false, "此参数是仅扫描常用web端口")
	ScanCmd.Flags().Bool("dbs", false, "此参数是仅扫描常用数据库端口")
	ScanCmd.Flags().Bool("risk", false, "此参数是仅扫描高危端口")

}
