package Protocol

import (
	"fmt"
	"strings"
)

const sshPrefix = "SSH-"

// IsSSHProtocol SSH 协议识别
func IsSSHProtocol(line string) bool {
	return strings.HasPrefix(line, sshPrefix)
}

func IsSSHProtocolApp(str string) string {
	str = strings.ReplaceAll(strings.ReplaceAll(str, "\r", ""), "\n", "")
	if strings.Contains(str, "Comware") {
		return fmt.Sprintf("%-5s|%s", "H3C", str)
	}
	if strings.Contains(str, "Cisco") {
		return fmt.Sprintf("%-5s|%s", "Cisco", str)
	}
	return str
}
