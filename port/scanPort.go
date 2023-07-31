package port

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

func scanPort() {
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

	allcount = len(iplist) * len(portlist)

	for _, ip := range iplist {
		for _, port := range portlist {
			ch <- struct{}{}
			wg.Add(1)
			outputMux.Lock()
			go IsPortOpen(ip, port)
			outputMux.Unlock()
		}
	}

	wg.Wait()
	time.Sleep(time.Second * 1) //等待1秒是为了正常显示进度条
	fmt.Printf("\r+------------------------------+\n")

	if save {
		if len(infolist) > 0 || len(iplist) > 0 {
			saveXlsx(infolist, iplist)
		}
	}

	fmt.Printf("[*] 扫描主机: %v 存活端口: %v \n",
		color.GreenString("%d", len(iplist)),
		color.GreenString("%d", len(infolist)),
	)

}
