package scan

import (
	"bufio"
	"fmt"
	Protocol2 "golin/Protocol"
	"net"
	"strings"
	"time"
)

// parseProtocol 协议/组件分析：有的基于默认端口去对应服务
func parseProtocol(conn net.Conn, host, port string, Poc bool) string {

	if protocol, ok := portProtocols[port]; ok {
		return protocol
	}

	if err := conn.SetReadDeadline(time.Now().Add(time.Duration(Timeout) * time.Second)); err != nil {
		return ""
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		line = ""
	}

	switch {
	case Protocol2.IsSSHProtocol(line):
		return Protocol2.IsSSHProtocolApp(line)

	case strings.HasPrefix(line, "220"):
		return "FTP"

	case Protocol2.IsRedisProtocol(conn):
		return "数据库|Redis"

	case Protocol2.IsTelnet(conn):
		return "Telnet"

	case Protocol2.IsPgsqlProtocol(host, port):
		return "数据库|PostgreSQL"

	case Protocol2.IsRsyncProtocol(line):
		return "rsync|" + line

	default:
		isWeb := Protocol2.IsWeb(host, port, Timeout, Poc)
		for _, v := range isWeb {
			if v != "" {
				return fmt.Sprintf("%-5s| %s", "WEB应用", v)
			}
		}

	}

	isMySQL, version := Protocol2.IsMySqlProtocol(host, port)
	if isMySQL {
		return fmt.Sprintf("数据库|MySQL:%s", version)
	}

	return defaultPort(port)
}

func defaultPort(port string) string {
	defMap := map[string]string{
		"3306":  "数据库|MySQL",
		"23":    "Telnet",
		"21":    "FTP",
		"80":    "WEB应用",
		"443":   "WEB应用",
		"5236":  "数据库|达梦",
		"61616": "ActiveMQ",
	}
	value, exists := defMap[port]
	if exists {
		return value
	}
	return ""
}
