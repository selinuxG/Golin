package Protocol

import "strings"

// IsRsyncProtocol Rsync 协议识别
func IsRsyncProtocol(line string) bool {
	return strings.HasPrefix(line, "@RSYNCD")
}
