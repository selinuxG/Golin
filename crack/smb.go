package crack

import (
	"fmt"
	"github.com/stacktitan/smb/smb"
)

func smbcon(ip, user, passwd string, port, timeout int) {
	defer func() {
		wg.Done()
		<-ch
	}()

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
			end(ip, user, passwd, port)
		}
	} else {
		fmt.Println(err)
	}
	return
}
