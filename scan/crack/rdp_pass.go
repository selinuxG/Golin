//go:build !rdp

package crack

import (
	"context"
)

// linux系统是不支持rdp弱口令扫描
func rdpcon(cancel context.CancelFunc, ip, user, passwd string, port, timeout int) {
}
