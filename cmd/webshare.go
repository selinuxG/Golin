/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

// webshareCmd represents the webshare command
var webshareCmd = &cobra.Command{
	Use:   "webshare",
	Short: "通过web形式共享目录",
	Long:  `基于http形式共享指定目录`,
	Run: func(cmd *cobra.Command, args []string) {

		port, err := cmd.Flags().GetString("port")
		if err != nil {
			fmt.Println(err)
			return
		}

		path, err := cmd.Flags().GetString("path")
		if err != nil {
			fmt.Println(err)
			return
		}
		if path != "./" {
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Printf("\x1b[%dm错误🤷‍ %s共享目录不存在！ \x1b[0m\n", 31, path)
				os.Exit(3)
			}
		}
		fmt.Printf("\x1b[%dm✔‍ 完成,已启动webshare 端口:%s 共享目录:%s \x1b[0m\n", 34, port, path)

		http.Handle("/", http.FileServer(http.Dir(path))) //把当前文件目录作为共享目录
		//如果是windos自动打开
		if runtime.GOOS == "windows" {
			url := "http://127.0.0.1:" + port + "/"
			cmd := exec.Command("cmd", "/C", "start "+url)
			cmd.Run()

		}
		http.ListenAndServe(":"+port, nil)
	},
}

func init() {
	rootCmd.AddCommand(webshareCmd)
	webshareCmd.Flags().StringP("port", "p", "11111", "启动端口默认是11111")
	webshareCmd.Flags().StringP("path", "a", "./", "共享目录默认是当前目录")
}
