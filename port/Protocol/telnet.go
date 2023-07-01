package Protocol

import (
	"bufio"
	"io"
	"net"
	"strings"
)

var (
	tel = []string{"password", "login", "huawei", "telnet", "copyright", "username", "password"}
)

func IsTelnet(conn net.Conn) bool {
	//_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	reader := bufio.NewReader(conn)
	data, _ := io.ReadAll(reader)

	for _, v := range tel {
		if strings.Count(strings.ToLower(string(data)), v) > 0 {
			return true
		}
	}

	return false
}
