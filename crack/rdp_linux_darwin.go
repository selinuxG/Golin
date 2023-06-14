//go:build linux || darwin

package crack

import (
	"context"
)

func rdpcon(ctx context.Context, cancel context.CancelFunc, ip, user, passwd string, port, timeout int) {
}
