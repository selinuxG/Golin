//go:build windows

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golin/run/windows"
	"math/rand"
)

// windowsCmd represents the execl command
var windowsCmd = &cobra.Command{
	Use:   "windows",
	Short: "读取安全策略生成html",
	Long:  `读取安全策略生成html`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("大概20s,做个脑筋急转弯吧朋友？")
		haha := []string{`你能做，我能做，大家都做；一个人能做，两个人不能一起做。这是做什么？`, `什么事每人每天都必须认真的做?`, `哪一个月有二十八天?`, `你能以最快速度，把冰变成水吗？`, `冬天，宝宝怕冷，到了屋里也不肯脱帽。可是他见了一个人乖乖地脱下帽，那人是谁？`, `老王一天要刮四五十次脸，脸上却仍有胡子。这是什么原因？`, `有一个字，人人见了都会念错。这是什么字？`}
		go func() {
			randomIndex := rand.Intn(len(haha))
			randomValue := haha[randomIndex]
			fmt.Println(randomValue)
		}()
		windows.Windows()
	},
}

func init() {
	rootCmd.AddCommand(windowsCmd)
}
