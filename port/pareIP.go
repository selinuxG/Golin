package port

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// parseIP解析IP地址范围 支持：192.168.1.1-100、192.168.1.1/24、192.168.1.1、baidu.com
func parseIP(ip string) {

	if strings.Count(ip, ",") == 0 {
		ip = ip + ","
	}

	for _, p := range strings.Split(ip, ",") {

		if p == "" {
			continue
		}

		switch {
		case strings.Contains(p, "-"): //起始-结束ip
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
