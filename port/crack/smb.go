package crack

import (
	"context"
	"github.com/stacktitan/smb/smb"
	"sync"
)

func smbcon(ctx context.Context, cancel context.CancelFunc, ip, user, passwd string, port, timeout int, ch <-chan struct{}, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		<-ch
	}()
	select {
	case <-ctx.Done():
		return
	default:
	}

	options := smb.Options{
		Host:        ip,
		Port:        port,
		User:        user,
		Domain:      "",
		Workstation: "",
		Password:    passwd,
	}
	debug := false
	session, err := smb.NewSession(options, debug)
	if err == nil {
		defer session.Close()
		if session.IsAuthenticated {
			end(ip, user, passwd, port, "SMB")
			cancel()
		}
	}
	return
}
