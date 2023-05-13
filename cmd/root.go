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
	Long:  `[-] 此工具是基于Golang多线程的模式开发，目的是批量执行各种设备的命令并记录;`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	cmd.Help()
	//},
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
\_____/ \_____/ |_____| |_| |_|  \_| ` + global.Version + ":" + global.Releasenotes)
}
