package port

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	linuxcount     int                       //linux 主机数量
	windowscount   int                       //windows 主机数量
	pingwg         = sync.WaitGroup{}        //ping的并发数
	pingch         = make(chan struct{}, 50) //ping的缓冲区数量
	filteredIPList []string                  //存放失败主机列表
)

func SanPing() {
	fmt.Printf("%s\n", "下发PING任务...\n+------------------------------+")
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
			yesPing, pingOS := NetWorkPing(ip) //是否ping通、ttl值
			if !yesPing {
				outputMux.Lock()
				filteredIPList = append(filteredIPList, ip) //ping不通放入待删除切片中不进行检测
				outputMux.Unlock()
			} else {
				outputMux.Lock()
				fmt.Printf("| %-15s|%-5s\n", ip, pingOS)
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
func NetWorkPing(ip string) (bool, string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "2", "-w", "1", ip)
	} else {
		cmd = exec.Command("ping", "-c", "2", "-W", "1", ip)
	}

	output, err := cmd.Output()
	if err != nil {
		return false, ""
	}
	outttl := strings.ToLower(string(output)) //所有大写转换为小写
	if strings.Contains(outttl, "ttl") {
		// Extract TTL value
		re := regexp.MustCompile(`ttl=(\d+)`)
		ttlStr := re.FindStringSubmatch(outttl)

		if len(ttlStr) > 1 {
			ttl, _ := strconv.Atoi(ttlStr[1])
			switch {
			case ttl <= 64:
				return true, "linux"
			case ttl <= 128:
				return true, "Windows"
			default:
				return true, "Unknown"
			}
		}
	}
	return false, ""
}
