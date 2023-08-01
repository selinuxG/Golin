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
	"25":    "邮件传输协议:SMTP",
	"53":    "域名解析:DNS",
	"110":   "邮件传输协议:POP3",
	"135":   "RPC 服务",
	"137":   "NetBIOS 名称服务",
	"138":   "NetBIOS 数据报服务",
	"139":   "NetBIOS 会话服务",
	"161":   "网络管理协议:SNMP",
	"162":   "网络管理报警和事件提醒:SNMP-trap",
	"143":   "邮件传输协议:IMAP",
	"445":   "Microsoft 的 SMB 协议",
	"465":   "带有 SSL 安全的 SMTP：SMTPS",
	"514":   "系统日志服务:syslog",
	"587":   "邮件提交协议（MSA）:Submission",
	"993":   "带有 SSL 安全的 IMAP：IMAPS",
	"995":   "带有 SSL 安全的 POP3：POP3S",
	"1080":  "SOCKS 代理",
	"1194":  "开放VPN",
	"1433":  "数据库:SQLServer",
	"1521":  "数据库:Oracle",
	"1723":  "对端到对端协议:PPTP",
	"2049":  "网络文件系统:NFS",
	"3389":  "远程桌面协议:RDP",
	"5601":  "系统管理:Kibana",
	"5900":  "虚拟网络计算:VNC",
	"5901":  "虚拟网络计算:VNC",
	"6000":  "X11",
	"6443":  "容器编排系统:Kubernetes",
	"9000":  "分布式计算框架:Hadoop",
	"5672":  "消息队列:RabbitMq",
	"9200":  "数据库:ElasticSearch",
	"27017": "数据库:MondoDB",
}

// parseProtocol 协议/组件分析：有的基于默认端口去对应服务
func parseProtocol(conn net.Conn, host, port string) string {

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
		return strings.ReplaceAll(strings.ReplaceAll(line, "\r", ""), "\n", "")

	case strings.HasPrefix(line, "220"):
		return "FTP"

	case Protocol.IsRedisProtocol(conn):
		return "数据库:Redis"

	case Protocol.IsTelnet(conn):
		return "Telnet"

	case Protocol.IsPgsqlProtocol(host, port):
		return "数据库:PostgreSQL"

	case Protocol.IsMongoDBProtocol(conn):
		return "数据库:MongoDB"

	default:
		isWeb := Protocol.IsWeb(host, port, Timeout)
		if isWeb != "" {
			return fmt.Sprintf(" %-6s| %s", "WEB应用", isWeb)
		}
	}

	isMySQL, version := Protocol.IsMySqlProtocol(host, port)
	if isMySQL {
		return fmt.Sprintf("数据库:MySQL:%s", version)
	}

	return ""
}
