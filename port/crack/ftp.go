package crack

import (
	"context"
	"fmt"
	"github.com/jlaffaye/ftp"
	"time"
)

func ftpcon(cancel context.CancelFunc, ip, user, passwd string, port, timeout int) {
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
