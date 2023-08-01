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
	Short: "A brief description of your application",
	Long:  `[-] 此工具是基于多线程模式开发，目的是进行快速准确的等保核查、端口扫描、组件识别、子域名扫描、目录扫描等功能;`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	golin := fmt.Sprintf(`
 _____   _____   _       _   __   _  
/  ___| /  _  \ | |     | | |  \ | | 
| |     | | | | | |     | | |   \| | 
| |  _  | | | | | |     | | | |\   | 
| |_| | | |_| | | |___  | | | | \  | 
\_____/ \_____/ |_____| |_| |_|  \_| `)

	fmt.Printf("%s\nVersion: %s\ngithub : %s\nauthor : %s\n\n",
		color.RedString("%s", golin),
		color.GreenString("%s %s", global.Version, global.Releasenotes),
		color.BlueString("%s", "https://github.com/selinuxG/Golin"),
		color.MagentaString("%s", "gaoyeshang"),
	)
}
