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

const clearLine = "\033[2K\r"                                            //清除当前行
const portformatString = clearLine + "\r| %-2s | %-15s | %-4s |%-50s \n" //端口存活的占位符

func scanPort() {
	defer func() {
		end()                      //输出总体结果
		saveXlsx(infolist, iplist) //结果保存文件
	}()
	checkPing()

	allcount = uint32(len(iplist) * len(portlist))

	for _, ip := range iplist {
		for _, port := range portlist {
			ch <- struct{}{}
			wg.Add(1)
			go IsPortOpen(ip, port) //启动线程扫描端口是否存活
		}
	}

	wg.Wait()
	global.Percent(donecount, allcount)

}

// IsPortOpen 判断端口是否开放并进行xss、poc、弱口令扫描
func IsPortOpen(host, port string) {

	defer func() {
		wg.Done()
		<-ch
		atomic.AddUint32(&donecount, 1)
		global.Percent(donecount, allcount)
	}()

	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(Timeout)*time.Second)
	if err != nil {
		return
	}
	Protocol := parseProtocol(conn, host, port, Xss, Poc) //识别协议、xss、poc扫描
	thisINFO := INFO{host, port, Protocol}

	fmt.Printf(portformatString, printGreen("%v", "✓"), thisINFO.Host, thisINFO.Port, thisINFO.Protocol) //端口存活信息

	outputMux.Lock()
	infolist = append(infolist, INFO{host, port, thisINFO.Protocol})
	outputMux.Unlock()

	if Carck {
		protocol := strings.ToLower(thisINFO.Protocol)
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

// printGreen 绿色编码输出
func printGreen(format string, a ...interface{}) string {
	return color.GreenString(format, a...)
}

// end 运行结束是输出,输出一些统计信息
func end() {
	fmt.Printf("\r+------------------------------+\n")
	fmt.Printf("\r[*] 存活主机:%v 存活端口:%v ssh:%v rdp:%v web服务:%v 数据库:%v \n",
		printGreen("%v", len(iplist)),
		printGreen("%v", len(infolist)),
		printGreen("%v", protocolExistsAndCount("ssh")),
		printGreen("%v", protocolExistsAndCount("rdp")),
		printGreen("%v", protocolExistsAndCount("WEB应用")),
		printGreen("%v", protocolExistsAndCount("数据库")),
	)

}
