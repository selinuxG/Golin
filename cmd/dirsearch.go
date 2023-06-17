package cmd

import (
	"github.com/spf13/cobra"
	"golin/dirscan"
)

// networkCmd represents the network command

var dirsearch = &cobra.Command{
	Use:   "dirsearch",
	Short: "进行网页目录扫描",
	Run:   dirscan.ParseFlags,
}

func init() {
	rootCmd.AddCommand(dirsearch)
	dirsearch.Flags().StringP("url", "u", "", "指定扫描url")
	dirsearch.Flags().StringP("proxy", "p", "", "指定代理")
	dirsearch.Flags().IntP("chan", "c", 30, "并发数量")
	dirsearch.Flags().IntP("timeout", "t", 3, "超时等待时常/s")
	dirsearch.Flags().StringP("file", "f", "", "此参数是指定读取的字典")
	dirsearch.Flags().StringP("code", "", "200", "此参数是指定状态码,多个按照,分割")
	dirsearch.Flags().IntP("wait", "", 0, "此参数是每次访问后等待多长时间")
	dirsearch.Flags().Bool("cally", false, "此参数是运行爬虫")
}
