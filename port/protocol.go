package port

import (
	"bufio"
	"fmt"
	"golin/port/Protocol"
	"net"
	"strings"
	"time"
)

var portProtocols = map[string]string{
	"25":    "SMTP",
	"53":    "DNS",
	"110":   "POP3",
	"135":   "RPC服务",
	"137":   "NetBIOS名称服务",
	"138":   "NetBIOS数据报服务",
	"139":   "NetBIOS会话服务",
	"161":   "SNMP",
	"162":   "SNMP-trap",
	"143":   "IMAP",
	"445":   "SMB",
	"465":   "SMTPS",
	"514":   "syslog",
	"587":   "Submission",
	"993":   "IMAPS",
	"995":   "POP3S",
	"1433":  "数据库|SqlServer",
	"1521":  "数据库|Oracle",
	"1723":  "PPTP",
	"2049":  "NFS",
	"3389":  "RDP",
	"5900":  "VNC",
	"5901":  "VNC",
	"6000":  "X11",
	"5672":  "RabbitMq",
	"27017": "数据库|MongoDB",
	"2181":  "ZooKeeper",
}

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
	case Protocol.IsSSHProtocol(line):
		return Protocol.IsSSHProtocolApp(line)

	case strings.HasPrefix(line, "220"):
		return "FTP"

	case Protocol.IsRedisProtocol(conn):
		return "数据库|Redis"

	case Protocol.IsTelnet(conn):
		return "Telnet"

	case Protocol.IsPgsqlProtocol(host, port):
		return "数据库|PostgreSQL"

	default:
		isWeb := Protocol.IsWeb(host, port, Timeout, Poc)
		for _, v := range isWeb {
			if v != "" {
				return fmt.Sprintf("%-5s| %s", "WEB应用", v)
			}
		}

	}

	isMySQL, version := Protocol.IsMySqlProtocol(host, port)
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
		"61616": "ActiveMQ",
	}
	value, exists := defMap[port]
	if exists {
		return value
	}
	return ""
}
