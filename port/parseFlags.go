package port

import (
	"fmt"
	"github.com/spf13/cobra"
	"golin/global"
	"math/rand"
	"net"
	"os"
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
	allcount  int
	donecount int
	outputMux sync.Mutex
)

func scanPort() {

	//ping检测 true不进行检测 false检测

	var filteredIPList []string
	if !NoPing {
		fmt.Printf("[-] 开始探测存活主机.....\n")
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
					fmt.Printf("[*] 存活主机: %s 操作系统：%s \n", ip, pingOS)
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

	fmt.Printf("[-] 检测主机数量: %d 主机探测端口数量: %d 并发数: %d 端口连接超时:%d :)\n", len(iplist), len(removeDuplicates(portlist)), chancount, Timeout)

	allcount = len(iplist) * len(portlist)

	for _, ip := range iplist {
		for _, port := range portlist {
			ch <- struct{}{}
			wg.Add(1)
			go IsPortOpen(ip, port)
			continue
		}

	}
	wg.Wait()
	time.Sleep(time.Second * 1) //等待1秒是为了正常显示进度条
	fmt.Printf("\r+-----------------------------------------------------+\n")

	if save {
		saveXlsx(infolist)
	}

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
		fmt.Printf("\r| %-15s | %-5s | %-5s |%s \n", host, port, "open", parseprotocol)
		infolist = append(infolist, INFO{host, port, parseprotocol})
		outputMux.Unlock()

	}
}
