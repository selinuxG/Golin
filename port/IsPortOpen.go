package port

import (
	"fmt"
	"github.com/fatih/color"
	"golin/global"
	"golin/port/crack"
	"net"
	"strings"
	"time"
)

// IsPortOpen 判断端口是否开放
func IsPortOpen(host, port string) {

	defer func() {
		wg.Done()
		<-ch
		donecount += 1
		global.Percent(&outputMux, donecount, allcount)
	}()

	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(Timeout)*time.Second)
	if err != nil {
		return
	}

	outputMux.Lock()
	parseprotocol := parseProtocol(conn, host, port) //识别协议
	fmt.Printf("\r| %-2s | %-15s | %-4s |%s \n", fmt.Sprintf("%s", color.GreenString("%s", "✓")), host, port, parseprotocol)
	infolist = append(infolist, INFO{host, port, parseprotocol})
	outputMux.Unlock()

	if Carck {
		protocol := strings.ToLower(parseprotocol)
		switch {
		case strings.Contains(protocol, "ssh"):
			crack.Run(host, port, Timeout, chancount, "ssh")
		case strings.Contains(protocol, "mysql"):
			crack.Run(host, port, Timeout, chancount, "mysql")
		case strings.Contains(protocol, "redis"):
			crack.Run(host, port, Timeout, chancount, "redis")
		case strings.Contains(protocol, "pgsql"):
			crack.Run(host, port, Timeout, chancount, "pgsql")
		case strings.Contains(protocol, "sqlserver"):
			crack.Run(host, port, Timeout, chancount, "sqlserver")
		case strings.Contains(protocol, "ftp"):
			crack.Run(host, port, Timeout, chancount, "ftp")
		case strings.Contains(protocol, "smb"):
			crack.Run(host, port, Timeout, chancount, "smb")
		case strings.Contains(protocol, "ftp"):
			crack.Run(host, port, Timeout, chancount, "ftp")
		case strings.Contains(protocol, "telnet"):
			crack.Run(host, port, Timeout, chancount, "telnet")
		case strings.Contains(protocol, "tomcat"):
			crack.Run(host, port, Timeout, chancount, "tomcat")
		}
	}
}
