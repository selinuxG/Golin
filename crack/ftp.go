package crack

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"time"
)

func ftpcon(ip, user, passwd string, port, timeout int) {
	defer func() {
		wg.Done()
		<-ch
	}()
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", ip, port), ftp.DialWithTimeout(time.Duration(timeout)*time.Second))
	if err == nil {
		err = c.Login(user, passwd)
		if err == nil {
			end(ip, user, passwd, port)
			_ = c.Quit()
		}
	}
	return
}
