//go:build linux

package crack

import (
	"context"
)

// linux系统是不支持rdp弱口令扫描
func rdpcon(ctx context.Context, cancel context.CancelFunc, ip, user, passwd string, port, timeout int) {
	defer func() {
		wg.Done()
		<-ch
	}()
	select {
	case <-ctx.Done():
		return
	default:
		cancel()
	}
}
