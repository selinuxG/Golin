package port

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golin/global"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	iplist    = []string{}              //扫描的端口
	portlist  = []string{}              //扫描的端口
	NoPing    bool                      //是否禁止ping监测
	ch        = make(chan struct{}, 30) //控制并发数
	wg        = sync.WaitGroup{}
	chancount int    //并发数量
	Timeout   int    //超时等待时常
	random    bool   //打乱顺序
	save      bool   //是否保存
	infolist  []INFO //成功的主机列表
)

type INFO struct {
	Host     string //主机
	Port     string //开放端口
	Protocol string //协议
}

func ParseFlags(cmd *cobra.Command, args []string) {
	ip, _ := cmd.Flags().GetString("ip")
	if ip == "" {
		fmt.Printf("[-] 未指定扫描主机!通过 golin port -i 指定,支持：192.168.1.1,192.168.1.1/24,192.168.1.1-100\n")
		os.Exit(1)
	}
	parseIP(ip)

	port, _ := cmd.Flags().GetString("port")
	parsePort(port)

	excludeport, _ := cmd.Flags().GetString("exclude") //去重端口以及排查过滤端口
	excludePort(excludeport)

	NoPing, _ = cmd.Flags().GetBool("noping")

	chancount, _ = cmd.Flags().GetInt("chan") //并发数量
	ch = make(chan struct{}, chancount)

	Timeout, _ = cmd.Flags().GetInt("time") //超时等待时常
	if Timeout <= 0 {
		Timeout = 3
	}

	random, _ = cmd.Flags().GetBool("random")
	save, _ = cmd.Flags().GetBool("save")

	scanPort()

}

var (
	allcount     int //IP*PORT的总数量
	donecount    int //线程技术的数量
	outputMux    sync.Mutex
	linuxcount   int //linux 主机数量
	windowscount int //windows 主机数量
)

func scanPort() {

	var filteredIPList []string
	if !NoPing {
		fmt.Printf("%v\n", color.GreenString("%s", "[*] 开始探测存活主机......"))
		pingwg := sync.WaitGroup{}
		for _, ip := range iplist {
			pingwg.Add(1)
			ip := ip
			go func() {
				defer pingwg.Done()
				yesPing, pingOS := global.NetWorkStatus(ip)
				if !yesPing {
					outputMux.Lock()
					filteredIPList = append(filteredIPList, ip) //ping不通放入待删除切片中不进行检测
					outputMux.Unlock()
				} else {
					outputMux.Lock()
					fmt.Printf("[-] 存活主机: %s 操作系统：%s \n", ip, pingOS)
					switch pingOS {
					case "Linux/Unix":
						linuxcount += 1
					case "Windows":
						windowscount += 1
					}
					outputMux.Unlock()
				}
			}()
		}
		pingwg.Wait()
		// Remove filteredIPList elements from iplist
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
	fmt.Println("+-----------------------------------------------------+")
	fmt.Printf("[*] Linux设备:%v Windows设备:%v 未识别:%v 共计存活:%v\n[*] 开始扫描端口:%v 并发数:%v 端口连接超时:%v\n",
		color.GreenString("%d", linuxcount),
		color.GreenString("%d", windowscount),
		color.RedString("%d", len(iplist)-linuxcount-windowscount),
		color.GreenString("%d", len(iplist)),
		color.GreenString("%d", len(portlist)),
		color.GreenString("%d", chancount),
		color.GreenString("%d", Timeout),
	)

	allcount = len(iplist) * len(portlist)

	for _, ip := range iplist {
		for _, port := range portlist {
			ch <- struct{}{}
			wg.Add(1)
			go IsPortOpen(ip, port)
		}
	}

	wg.Wait()
	time.Sleep(time.Second * 1) //等待1秒是为了正常显示进度条
	fmt.Printf("\r+-----------------------------------------------------+\n")

	if save {
		if len(infolist) > 0 || len(iplist) > 0 {
			saveXlsx(infolist, iplist)
		}
	}

	fmt.Printf("[*] 存活主机: %v 存活端口: %v 数据库: %v Web: %v SSH: %v RDP: %v \n",
		color.GreenString("%d", len(iplist)),
		color.GreenString("%d", len(infolist)),
		color.GreenString("%d", countPortOccurrences("数据库")),
		color.GreenString("%d", countPortOccurrences("Web应用")),
		color.GreenString("%d", countPortOccurrences("SSH")),
		color.GreenString("%d", countPortOccurrences("RDP")),
	)

}

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

	if err == nil {
		outputMux.Lock()
		parseprotocol := parseProtocol(conn, host, port) //识别协议
		fmt.Printf("\r| %-15s | %-5s | %-2s |%s \n", host, port, "✓", parseprotocol)
		infolist = append(infolist, INFO{host, port, parseprotocol})
		outputMux.Unlock()

	}
}

// countPortOccurrences 接受协议特征，返回共计数量
func countPortOccurrences(protocol string) int {
	count := 0
	for _, info := range infolist {
		if strings.Contains(info.Protocol, protocol) {
			count += 1
		}
	}
	return count
}
