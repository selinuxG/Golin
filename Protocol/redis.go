package Protocol

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// IsRedisProtocol Redis 协议识别
func IsRedisProtocol(conn net.Conn) bool {
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))

	// 发送一个 PING 命令
	_, err := fmt.Fprintf(conn, "*1\r\n$4\r\nPING\r\n")
	if err != nil {
		return false
	}

	// 读取返回的数据
	response, err := readResponse(conn)
	if err == nil {
		if strings.Contains(response, "NOAUTH") || strings.Contains(response, "PONG") {
			return true
		}

	}

	return false

}

// readResponse 返回第一行字符串以及错误
func readResponse(conn net.Conn) (string, error) {
	// 创建一个新的 reader
	reader := bufio.NewReader(conn)

	// 读取一行数据
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return line, nil
}
