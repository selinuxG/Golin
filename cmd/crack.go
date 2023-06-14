package cmd

import (
	"github.com/spf13/cobra"
	"golin/crack"
)

var expCmd = &cobra.Command{
	Use:   "crack",
	Short: "暴力破解功能,目前支持：ssh、mysql、redis、postgresql、sqlserver、ftp",
}

// 破解 ssh
var ssh = &cobra.Command{
	Use:   "ssh",
	Short: "ssh暴力破解",
	Run:   crack.Run,
}

// 破解 mysql
var mysql = &cobra.Command{
	Use:   "mysql",
	Short: "mysql暴力破解",
	Run:   crack.Run,
}

// 破解 redis
var redis = &cobra.Command{
	Use:   "redis",
	Short: "redis暴力破解",
	Run:   crack.Run,
}

// 破解 pgsql
var pgsql = &cobra.Command{
	Use:   "pgsql",
	Short: "pgsql暴力破解",
	Run:   crack.Run,
}

// 破解 sqlserver
var sqlserver = &cobra.Command{
	Use:   "sqlserver",
	Short: "sqlserver暴力破解",
	Run:   crack.Run,
}

// 破解 ftp
var ftp = &cobra.Command{
	Use:   "ftp",
	Short: "ftp暴力破解",
	Run:   crack.Run,
}

// 破解 rdp
var rdp = &cobra.Command{
	Use:   "rdp",
	Short: "rdp",
	Run:   crack.Run,
}

// 破解 smb
var smb = &cobra.Command{
	Use:   "smb",
	Short: "smb",
	Run:   crack.Run,
}

func init() {
	rootCmd.AddCommand(expCmd)
	commands := []*cobra.Command{ssh, mysql, redis, pgsql, sqlserver, ftp, rdp, smb}
	for _, cmd := range commands {
		expCmd.AddCommand(cmd)
		cmd.Flags().StringP("ip", "i", "", "此参数是指定验证的IP")
		cmd.Flags().IntP("port", "P", 0, "此参数是指定验证的端口")
		cmd.Flags().StringP("user", "u", "", "此参数是指定用户文件")
		cmd.Flags().StringP("passwd", "p", "", "此参数是指定密码文件")
		cmd.Flags().StringP("fire", "f", "", "此参数是指定主机列表，格式IP:Port 一行一个")
		cmd.Flags().IntP("chan", "c", 30, "并发数量")
		cmd.Flags().IntP("time", "t", 3, "超时等待时常/s")
		cmd.Flags().Bool("noping", false, "此参数是指定不运行ping监测")
	}
}
