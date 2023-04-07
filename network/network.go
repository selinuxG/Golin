package network

import (
	"fmt"
	"github.com/spf13/cobra"
	"golin/config"
)

var (
	zlog = config.Log
)

func Networkrun(cmd *cobra.Command, args []string) {
	syslog, err := cmd.Flags().GetBool("syslog")
	if err != nil {
		fmt.Println(err)
		return
	}
	if syslog {
		if len(args) == 2 {
			Rsyslog(args[0], args[1])
			return
		}
		zlog.Warn("参数错误,接受2个参数,第一个是协议（tcp/udp）,第二个是端口")
		return
	}
}
