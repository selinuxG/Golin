package crack

import (
	"context"
	"fmt"
	"github.com/jlaffaye/ftp"
	"sync"
	"time"
)

func ftpcon(ctx context.Context, cancel context.CancelFunc, ip, user, passwd string, port, timeout int, ch <-chan struct{}, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		<-ch
	}()
	select {
	case <-ctx.Done():
		return
	default:
	}

	c, err := ftp.Dial(fmt.Sprintf("%s:%d", ip, port), ftp.DialWithTimeout(time.Duration(timeout)*time.Second))
	if err == nil {
		err = c.Login(user, passwd)
		if err == nil {
			end(ip, user, passwd, port, "FTP")
			_ = c.Quit()
			cancel()
		}
	}
	return
}
