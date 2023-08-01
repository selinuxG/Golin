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

	if Carck {
		protocol := strings.ToLower(parseprotocol)
		//支持扫描的类型
		protocols := []string{"ssh", "mysql", "redis", "pgsql", "sqlserver", "ftp", "smb", "telnet", "tomcat", "rdp", "oracle", "mongodb"}

		for _, proto := range protocols {
			if strings.Contains(protocol, proto) { //不区分大小写
				crack.Run(host, port, Timeout, chancount, proto)
				break
			}
		}
	}

	outputMux.Unlock()
}
