package crack

import (
	"context"
	"github.com/stacktitan/smb/smb"
)

func smbcon(cancel context.CancelFunc, ip, user, passwd string, port, timeout int) {
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
