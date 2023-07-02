package Protocol

import (
	"bufio"
	"io"
	"net"
	"strings"
	"time"
)

var (
	tel = []string{"password", "login", "huawei", "telnet", "copyright", "username", "password"}
)

func IsTelnet(conn net.Conn) bool {
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	reader := bufio.NewReader(conn)
	data, _ := io.ReadAll(reader)

	for _, v := range tel {
		if strings.Contains(strings.ToLower(string(data)), v) {
			return true
		}
	}

	return false
}
