package port

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	linuxcount     int                       //linux 主机数量
	windowscount   int                       //windows 主机数量
	pingwg         = sync.WaitGroup{}        //ping的并发数
	pingch         = make(chan struct{}, 50) //ping的缓冲区数量
	filteredIPList []string                  //存放失败主机列表
)

func checkPing() {
	if !NoPing {
		SanPing()
		pingwg.Wait()
		// 删除ping失败的主机
		for _, ip := range filteredIPList {
			for i := 0; i < len(iplist); i++ {
				if iplist[i] == ip {
					iplist = append(iplist[:i], iplist[i+1:]...)
					break
				}
			}
		}
	}

	if random { //打乱主机顺序
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(iplist), func(i, j int) {
			iplist[i], iplist[j] = iplist[j], iplist[i]
		})
	}

	if !NoPing && len(iplist) == 0 {
		fmt.Printf("%s\n", color.RedString("%s", "[-] 通过尝试PING探测存活主机为0！可通过--noping跳过PING尝试"))
		return
	}

	fmt.Println("+------------------------------+")
	fmt.Printf("[*] Linux设备:%v Windows设备:%v 未识别:%v 共计存活:%v\n[*] 开始扫描端口:%v 并发数:%v 共计尝试:%v 端口连接超时:%v\n",
		color.GreenString("%d", linuxcount),
		color.GreenString("%d", windowscount),
		color.RedString("%d", len(iplist)-linuxcount-windowscount),
		color.GreenString("%d", len(iplist)),
		color.GreenString("%d", len(portlist)),
		color.GreenString("%d", chancount),
		color.GreenString("%d", len(iplist)*len(portlist)),
		color.GreenString("%d", Timeout),
	)
	fmt.Println("+------------------------------+")

}

func SanPing() {
	//fmt.Printf("%s\n", "下发PING任务...\n+------------------------------+")
	pingch = make(chan struct{}, chancount)
	for _, ip := range iplist {
		pingch <- struct{}{}
		pingwg.Add(1)
		ip := ip
		go func() {
			defer func() {
				pingwg.Done()
				<-pingch
			}()
			yesPing, pingOS, timems := NetWorkPing(ip) //是否ping通、ttl值
			if !yesPing {
				outputMux.Lock()
				filteredIPList = append(filteredIPList, ip) //ping不通放入待删除切片中不进行检测
				outputMux.Unlock()
			} else {
				outputMux.Lock()
				fmt.Printf("|%-5s| %-15s|%-7s|%-4s|%sms\n", color.GreenString("%s", "存活主机"), ip, pingOS, isPublicIP(net.ParseIP(ip)), timems)
				switch pingOS {
				case "linux":
					linuxcount += 1
				case "Windows":
					windowscount += 1
				}
				outputMux.Unlock()
			}
		}()
	}

}

// NetWorkPing 检查ping 返回是否可ping通以及操作系统
func NetWorkPing(ip string) (bool, string, string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "2", "-w", "1", ip)
	}
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("ping", "-c", "2", "-t", "1", ip)
	}
	if runtime.GOOS == "linux" {
		cmd = exec.Command("ping", "-c", "2", "-W", "1", ip)
	}
	output, err := cmd.Output()
	if err != nil {
		return false, "", ""
	}
	outttl := strings.ToLower(string(output)) //所有大写转换为小写

	// time
	re := regexp.MustCompile(`=(\d+)ms`)
	timems := ""
	timeStr := re.FindStringSubmatch(outttl)
	if len(timeStr) > 1 {
		inttime, _ := strconv.Atoi(timeStr[1])
		if inttime > 10 {
			timems = fmt.Sprintf("%s", color.RedString("%s", timeStr[1]))
		} else {
			timems = fmt.Sprintf("%s", color.GreenString("%s", timeStr[1]))
		}
	}

	if strings.Contains(outttl, "ttl") {
		// Extract TTL value
		re := regexp.MustCompile(`ttl=(\d+)`)
		ttlStr := re.FindStringSubmatch(outttl)

		if len(ttlStr) > 1 {
			ttl, _ := strconv.Atoi(ttlStr[1])
			switch {
			case ttl <= 64:
				return true, "linux", timems
			case ttl <= 128:
				return true, "Windows", timems
			default:
				return true, "Unknown", timems
			}
		}
	}
	return false, "", ""
}

// isPublicIP 检查是局域网还是互联网
func isPublicIP(IP net.IP) string {
	private := "互联网"
	// 定义私有网络的范围
	privateNetworks := []string{
		"10.0.0.0/8",     // 10.0.0.0 - 10.255.255.255
		"172.16.0.0/12",  // 172.16.0.0 - 172.31.255.255
		"192.168.0.0/16", // 192.168.0.0 - 192.168.255.255
	}

	for _, privateNet := range privateNetworks {
		_, ipnet, _ := net.ParseCIDR(privateNet)
		if ipnet.Contains(IP) {
			private = "局域网"
		}
	}

	return private
}
