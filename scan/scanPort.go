package scan

import (
	"fmt"
	"github.com/fatih/color"
	"golin/global"
	"golin/poc"
	"golin/scan/crack"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

const clearLine = "\033[2K\r"                                            //清除当前行
const portformatString = clearLine + "\r| %-2s | %-15s | %-4s |%-50s \n" //端口存活的占位符

func scanPort(donetime int) {
	defer func() {
		global.Percent(donecount, allcount) //输出100%的进度条,并且补全任务信息
		echoCrack()                         //输出弱口令资产
		echoPoc()                           //输出漏洞资产
		endEcho()                           //输出总体结果
		saveXlsx(infolist, iplist)          //结果保存文件
		global.StartScreenshotWorkers(10)   //启动WEB截图
		// 补全任务信息
		global.Job.EndTime = time.Now().Format("2006-01-02 15:04:05")
		global.Job.VulnerabilityCount = len(poc.ListPocInfo)
		global.Job.CrackCount = len(crack.MapCrackHost)
		endHtml() //输出html
	}()
	checkPing()

	allcount = uint32(len(iplist) * len(portlist))

	for _, ip := range iplist {
		for _, port := range portlist {
			ch <- struct{}{}
			wg.Add(1)
			go func(ip string, port string) {
				defer wg.Done()
				defer func() { <-ch }()
				done := make(chan struct{}, 1)
				go func() {
					IsPortOpen(ip, port)
					done <- struct{}{} //发送完成信号
				}()

				select {
				case <-done: // 正常完成
					return
				case <-time.After(time.Duration(donetime) * time.Minute): // 超时
					return
				}
			}(ip, port)
		}
	}
	wg.Wait()

}

// IsPortOpen 判断端口是否开放并进行xss、poc、弱口令扫描
func IsPortOpen(host, port string) {

	defer func() {
		//wg.Done()
		//<-ch
		atomic.AddUint32(&donecount, 1)
		global.Percent(donecount, allcount)
	}()

	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, time.Duration(Timeout)*time.Second)
	if err != nil {
		return
	}
	Protocol := parseProtocol(conn, host, port, Poc) //识别协议
	thisINFO := INFO{host, port, Protocol}
	outputMux.Lock()
	fmt.Printf(portformatString, printGreen("%v", "✓"), thisINFO.Host, thisINFO.Port, thisINFO.Protocol) //端口存活信息
	infolist = append(infolist, INFO{host, port, thisINFO.Protocol})
	outputMux.Unlock()
	// 永恒之蓝检测
	if port == "445" {
		if crack.MS17010Scan(host) {
			outputMux.Lock()
			poc.ListPocInfo = append(poc.ListPocInfo, poc.Flagcve{Url: host, Cve: "MSF17010", Flag: "永恒之蓝"})
			outputMux.Unlock()
		}
	}

	if Carck {
		protocol := strings.ToLower(thisINFO.Protocol)
		//支持遍历字典扫描的类型
		protocols := []string{"ssh", "mysql", "redis", "postgresql", "sqlserver", "ftp", "smb", "telnet", "tomcat", "rdp", "oracle"}
		for _, proto := range protocols {
			if strings.Contains(protocol, proto) { //不区分大小写
				if proto == "rdp" {
					if global.CrackRDP {
						crack.Run(host, port, Timeout, chancount, proto)
					}
					break // 不允许扫RDP就跳过
				} else {
					crack.Run(host, port, Timeout, chancount, proto)
					break
				}
			}
		}

		if strings.Contains(protocol, "mongodb") {
			crack.Mongodbcon(host, port)
		}

		if strings.Contains(protocol, "zookeeper") {
			poc.ZookeeperCon(host, port)
		}

		if strings.Contains(protocol, "Rsync") {
			crack.Rsync(host, port)
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

func printRed(format string, a ...interface{}) string {
	return color.RedString(format, a...)
}
