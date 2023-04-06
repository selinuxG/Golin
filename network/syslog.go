package network

import (
	"fmt"
	"strconv"
)

func Rsyslog(protocol, port string) {
	intport, err := strconv.Atoi(port)
	if err != nil {
		zlog.Warn("syslog端口转换失败！")
		return
	}

	if intport > 65535 || intport == 0 {
		zlog.Warn("非法端口！")
		return
	}
	if protocol == "tcp" || protocol == "udp" {
		fmt.Println(protocol)
		return

	}
	zlog.Warn("协议只能是tcp/udp！")
}
