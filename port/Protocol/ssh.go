package Protocol

import "strings"

const sshPrefix = "SSH-"

// IsSSHProtocol SSH 协议识别
func IsSSHProtocol(line string) bool {
	return strings.HasPrefix(line, sshPrefix)
}
