package cmd

import (
	"fmt"
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
	fmt.Println(`
 _____   _____   _       _   __   _  
/  ___| /  _  \ | |     | | |  \ | | 
| |     | | | | | |     | | |   \| | 
| |  _  | | | | | |     | | | |\   | 
| |_| | | |_| | | |___  | | | | \  | 
\_____/ \_____/ |_____| |_| |_|  \_| ` + "\nhttps://github.com/selinuxG/Golin-cli " + global.Version + ":" + global.Releasenotes)
	//fmt.Printf("\n")
}
