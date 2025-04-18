package scan

import (
	"crypto/tls"
	"fmt"
	"github.com/fatih/color"
	"golin/global"
	"golin/poc"
	"golin/scan/crack"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const clearLine = "\033[2K\r"                                            //清除当前行
const portformatString = clearLine + "\r| %-2s | %-15s | %-4s |%-50s \n" //端口存活的占位符

func scanPort(donetime int) {
	defer func() {
		global.Percent(donecount, allcount) //输出100%的进度条
		echoCrack()                         //输出弱口令资产
		echoPoc()                           //输出漏洞资产
		end()                               //输出总体结果
		saveXlsx(infolist, iplist)          //结果保存文件
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
				if proto == "rdp" && checkTLSVersion(host, port) != nil {
					break
				}
				crack.Run(host, port, Timeout, chancount, proto)
				break
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

// checkTLSVersion 如果tls版本低于1.2则不进行rdp扫描
func checkTLSVersion(host, port string) error {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host, port), conf)
	if err != nil {
		return fmt.Errorf("failed to establish connection: %w", err)
	}
	defer conn.Close()

	state := conn.ConnectionState()

	// 检查协议版本
	if state.Version < tls.VersionTLS12 {
		return fmt.Errorf("insecure protocol version")
	}

	return nil
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

// end 运行结束是输出,输出一些统计信息
func end() {
	vulnerablehost, _ := calculateVulnerablePercentage(iplist, crack.MapCrackHost, poc.ListPocInfo)
	fmt.Printf("\r+------------------------------------------------------------+\n")
	fmt.Printf("\r[*] 漏洞主机:%v Linux:%v Windows:%v 存活端口:%v ssh:%v rdp:%v web:%v 数据库:%v 弱口令:%v 漏洞:%v \n",
		printRed("%v%v", vulnerablehost, printGreen("/"+strconv.Itoa(len(iplist)))),
		printGreen("%v", linuxcount),
		printGreen("%v", windowscount),
		printGreen("%v", len(infolist)),
		printGreen("%v", protocolExistsAndCount("ssh")),
		printGreen("%v", protocolExistsAndCount("rdp")),
		printGreen("%v", protocolExistsAndCount("WEB应用")),
		printGreen("%v", protocolExistsAndCount("数据库")),
		printGreen("%v", len(crack.MapCrackHost)),
		printGreen("%v", len(poc.ListPocInfo)),
	)
	if global.SaveIMG {
		couunt, err := global.CountDirFiles(global.SsaveIMGDIR)
		if err != nil {
			couunt = 0
		}
		fmt.Printf("[*] Web扫描截图保存目录：%v 当前共计截图数量：%v\n",
			printGreen("%v", global.SsaveIMGDIR), printGreen("%v", couunt))
	}
}

// calculateVulnerablePercentage 计算有漏洞的IP数量和百分比
func calculateVulnerablePercentage(iplist []string, crackHosts map[crack.HostPort]crack.SussCrack, pocInfos []poc.Flagcve) (int, float64) {
	uniqueIPs := make(map[string]struct{})

	// 正则表达式用于从URL中提取IP地址
	ipRegex := regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)

	// 从crack.MapCrackHost中收集所有Host
	for hostPort := range crackHosts {
		if ip := net.ParseIP(hostPort.Host); ip != nil {
			uniqueIPs[hostPort.Host] = struct{}{}
		}
	}

	// 从Flagcve的URL中提取IP地址
	for _, pocInfo := range pocInfos {
		if matches := ipRegex.FindStringSubmatch(pocInfo.Url); len(matches) > 0 {
			if ip := net.ParseIP(matches[0]); ip != nil {
				uniqueIPs[matches[0]] = struct{}{}
			}
		}
	}

	// 计算有漏洞的IP数量
	vulnerableCount := len(uniqueIPs)
	totalCount := len(iplist)

	// 计算百分比
	if totalCount == 0 {
		return vulnerableCount, 0.0
	}
	percentage := (float64(vulnerableCount) / float64(totalCount)) * 100.0
	return vulnerableCount, percentage
}
