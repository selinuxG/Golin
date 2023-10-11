package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"golin/global"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "golin",
	Short: "弱口令检测、 漏洞扫描、端口扫描（协议识别，组件识别）、等保模拟定级、自动化运维、等保工具",
	Long:  `主机存活探测、漏洞扫描、子域名扫描、端口扫描、各类服务数据库爆破、poc扫描、xss扫描、webtitle探测、web指纹识别、web敏感信息泄露、web目录浏览、web文件下载、等保安全风险问题风险自查等； 弱口令/未授权访问：40余种； WEB组件识别：300余种； 漏洞扫描：XSS、任意文件访问、任意命令执行、敏感信息泄露、默认账户密码...； 资产扫描：扫描存活主机->判断存活端口->识别协议/组件->基于组件协议进行弱口令、漏洞扫描->输出报告`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	golin := fmt.Sprintf(`

  ▄████  ▒█████   ██▓     ██▓ ███▄    █ 
 ██▒ ▀█▒▒██▒  ██▒▓██▒    ▓██▒ ██ ▀█   █ 
▒██░▄▄▄░▒██░  ██▒▒██░    ▒██▒▓██  ▀█ ██▒
░▓█  ██▓▒██   ██░▒██░    ░██░▓██▒  ▐▌██▒
░▒▓███▀▒░ ████▓▒░░██████▒░██░▒██░   ▓██░
 ░▒   ▒ ░ ▒░▒░▒░ ░ ▒░▓  ░░▓  ░ ▒░   ▒ ▒ 
  ░   ░   ░ ▒ ▒░ ░ ░ ▒  ░ ▒ ░░ ░░   ░ ▒░
░ ░   ░ ░ ░ ░ ▒    ░ ░    ▒ ░   ░   ░ ░ 
      ░     ░ ░      ░  ░ ░           ░ 
                                        
 `)

	fmt.Printf("%s\nVersion: %s\ngithub : %s\nauthor : %s\n\n",
		color.GreenString("%s", golin),
		color.GreenString("%s %s", global.Version, global.Releasenotes),
		color.BlueString("%s", "https://github.com/selinuxG/Golin"),
		color.MagentaString("%s", "gaoyeshang -> VX:SelinuxG"),
	)
}
