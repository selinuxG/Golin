package port

import (
	"fmt"
	"golin/global"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// parseIP解析IP地址范围 支持：192.168.1.1-100、192.168.1.1/24、192.168.1.1、baidu.com、http://www.baidu.com
func parseIP(ip string) {
	for _, p := range strings.Split(ip, ",") {

		//识别端口
		rePort := regexp.MustCompile(`:\d{1,5}`)
		matchPort := rePort.FindString(p)

		//1、匹配CIDR子网掩码的地址 2、匹配IP 3、匹配第一个/之前的数据
		reCIDR := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/\d{1,2}\b`)
		reIP := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
		matchCIDR := reCIDR.FindString(p)
		if matchCIDR != "" {
			p = matchCIDR
		} else {
			matchIP := reIP.FindString(p)
			if matchIP != "" {
				p = matchIP
			} else {
				replacer := strings.NewReplacer("https://", "", "http://", "")
				p = replacer.Replace(p)
				index := strings.Index(p, "/")
				if index != -1 {
					p = p[:index]
				}
			}
		}
		p += matchPort

		checkPort := strings.Split(p, ":") //快速扫描
		if len(checkPort) == 2 {
			p = checkPort[0]
			portlist = []string{checkPort[1]}
		}

		switch {
		case strings.Contains(p, "-"): //起始-结束ip
			matched := reIP.MatchString(p) //域名可能会包含-符号，则会处罚此规则则直接跳过了，此bug目前先留存
			if !matched {
				continue
			}
			ipa := strings.Split(strings.Split(p, "-")[0], ".")

			start := strings.Split(ipa[3], "-")[0]
			startNum, _ := strconv.Atoi(start)

			end := strings.Split(p, "-")[1]
			endNum, _ := strconv.Atoi(end)

			ipa = ipa[:len(ipa)-1]
			ipStart := strings.Join(ipa, ".")

			if startNum > endNum {
				fmt.Printf("[-] 起始范围大于结束范围！\n")
				continue
			}

			for i := startNum; i <= endNum; i++ {
				nowIP := fmt.Sprintf("%s.%d", ipStart, i)
				iplist = append(iplist, nowIP)
			}

		case strings.Contains(p, "/"):
			_, ipnet, _ := net.ParseCIDR(p)
			for ipl := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ipl); incrementIP(ipl) {
				if len(ipl.To4()) == net.IPv4len {
					lastByte := ipl.To4()[3]
					if lastByte == 0 || lastByte == 255 { //起初0开头以及255结尾
						continue
					}
				}
				iplist = append(iplist, ipl.String())
			}

		default:
			if net.ParseIP(p) == nil {
				_, err := net.ResolveIPAddr("ip", p) //ip失败时判断是不是域名
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
			iplist = append(iplist, p)
		}
	}
}

// incrementIP 将ip地址增加1
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// parseFileIP解析扫描文件
func parseFileIP(path string) {
	if global.PathExists(path) {
		data, _ := os.ReadFile(path)
		for _, v := range strings.Split(string(data), "\n") {
			if v != "" {
				if strings.Contains(v, "-") {
					continue
				}
				replacer := strings.NewReplacer("\r", "", " ", "", "https://", "", "http://", "")
				v = replacer.Replace(v)
				if len(strings.Split(v, ":")) == 2 {
					ip := strings.Split(v, ":")[0]
					nowPort := strings.Split(v, ":")[1]
					portlist = append(portlist, nowPort)
					parseIP(ip)
					continue
				}
				parseIP(v)
			}
		}
	}
	return
}

// conNETLocal 获取当前可以获取的IP网段
func conNETLocal() {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, i := range interfaces {
		// 跳过虚拟网卡
		if strings.Contains(i.Name, "VirtualBox") || strings.Contains(i.Name, "VMware") {
			continue
		}

		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return
		}

		addresses, err := byName.Addrs()
		if err != nil {
			return
		}

		for _, v := range addresses {
			// 检查 IP 地址是否是 CIDR 表示法
			if ipNet, ok := v.(*net.IPNet); ok {
				// 过滤掉 IPv6 地址和 /32 子网
				if ipNet.IP.To4() != nil && ipNet.Mask.String() != "ffffffff" {
					// 获取 IP 地址的第一个字节，用于过滤掉以 169 和 127 开头的 IP 地址
					firstOctet := ipNet.IP.To4()[0]
					if firstOctet != 169 && firstOctet != 127 {
						// 对于每个网络，生成所有可能的 IP 地址
						for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); incrementIP(ip) {
							iplist = append(iplist, ip.String())
						}
					}

				}
			}
		}
	}

}
