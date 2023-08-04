package port

import (
	"fmt"
	"github.com/fatih/color"
	"golin/global"
	"golin/port/crack"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

func scanPort() {
	checkPing()

	allcount = uint32(len(iplist) * len(portlist))

	for _, ip := range iplist {
		for _, port := range portlist {
			ch <- struct{}{}
			wg.Add(1)
			go IsPortOpen(ip, port) //扫描端口是否存活
		}
	}

	wg.Wait()

	fmt.Printf("\r+------------------------------+\n")
	fmt.Printf("[*] 存活主机:%v 存活端口:%v ssh:%v rdp:%v web服务:%v 数据库:%v \n",
		color.GreenString("%d", len(iplist)),
		color.GreenString("%d", len(infolist)),
		color.GreenString("%d", protocolExistsAndCount("ssh")),
		color.GreenString("%d", protocolExistsAndCount("rdp")),
		color.GreenString("%d", protocolExistsAndCount("WEB应用")),
		color.GreenString("%d", protocolExistsAndCount("数据库")),
	)

	if save {
		if len(infolist) > 0 || len(iplist) > 0 {
			saveXlsx(infolist, iplist)
		}
	}

}

// IsPortOpen 判断端口是否开放并进行xss、poc、弱口令扫描
func IsPortOpen(host, port string) {

	defer func() {
		wg.Done()
		<-ch
		atomic.AddUint32(&donecount, 1)
		global.Percent(&outputMux, donecount, allcount)
	}()

	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(Timeout)*time.Second)
	if err != nil {
		return
	}

	parseprotocol := parseProtocol(conn, host, port, Xss, Poc) //识别协议、xss、poc扫描
	fmt.Printf("\033[2K\r")                                    // 擦除整行
	fmt.Printf("\r| %-2s | %-15s | %-4s |%-50s \n",
		fmt.Sprintf("%s", color.GreenString("%s", "✓")),
		host,
		port,
		parseprotocol,
	)
	outputMux.Lock()
	infolist = append(infolist, INFO{host, port, parseprotocol})
	outputMux.Unlock()

	if Carck {
		protocol := strings.ToLower(parseprotocol)
		//支持遍历字典扫描的类型
		protocols := []string{"ssh", "mysql", "redis", "pgsql", "sqlserver", "ftp", "smb", "telnet", "tomcat", "rdp", "oracle"}
		for _, proto := range protocols {
			if strings.Contains(protocol, proto) { //不区分大小写
				crack.Run(host, port, Timeout, chancount, proto)
				break
			}
		}

		//mongodb模式只进行验证未授权访问
		if strings.Contains(protocol, "mongodb") {
			crack.Mongodbcon(host, port)
		}
	}

}

// protocolExistsAndCount 接受一个协议特征返回总数
func protocolExistsAndCount(protocol string) (count int) {
	count = 0
	for _, info := range infolist {
		if strings.Contains(strings.ToLower(info.Protocol), strings.ToLower(protocol)) {
			count++
		}
	}
	return count
}
