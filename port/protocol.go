package port

import (
	"bufio"
	"fmt"
	"golin/port/Protocol"
	"net"
	"strings"
	"time"
)

const (
	ftpCode = "220"
)

func parseProtocol(conn net.Conn, host, port string) string {
	if err := conn.SetReadDeadline(time.Now().Add(time.Duration(Timeout) * time.Second)); err != nil {
		fmt.Println("Error setting read deadline:", err)
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

	isWeb := Protocol.IsWeb(host, port)

	switch {
	case port == "445":
		return "SMB"

	case port == "135":
		return "RPC"

	case port == "3389":
		return "RDP"

	case port == "161":
		return "SNMP"

	case port == "137" || port == "138" || port == "139":
		return "NetBIOS"

	case port == "514":
		return "Syslog"

	case port == "1433":
		return "SqlServer"

	case port == "1521":
		return "Oracle"

	case isWeb != "":
		return isWeb

	case Protocol.IsSSHProtocol(line):
		return strings.ReplaceAll(strings.ReplaceAll(line, "\r", ""), "\n", "")

	case strings.HasPrefix(line, ftpCode):
		return "FTP"

	case strings.Contains(port, "637"): //只识别判断包含 637 端口的是否为Redis协议
		if Protocol.IsRedisProtocol(conn) {
			return "Redis"
		}

	case Protocol.IsTelnet(conn):
		return "Telnet"

	case strings.Contains(port, "3306"): //只识别判断包含 3306 端口的是否为MySQL协议
		if Protocol.IsMySqlProtocol(host, port) {
			return "MySql"
		}

	case strings.Contains(port, "5432"): //只识别判断包含 5432 端口的是否为PostgreSQL协议
		if Protocol.IsPgsqlProtocol(host, port) {
			return "PostgreSQL"
		}

	default:
		return ""

	}
	return ""
}
