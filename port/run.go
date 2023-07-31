package port

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sync"
)

var (
	iplist    = []string{}              //扫描的端口
	portlist  = []string{}              //扫描的端口
	NoPing    bool                      //是否禁止ping监测
	Carck     bool                      //是否进行弱口令扫描
	ch        = make(chan struct{}, 30) //控制并发数
	wg        = sync.WaitGroup{}
	chancount int    //并发数量
	Timeout   int    //超时等待时常
	random    bool   //打乱顺序
	save      bool   //是否保存
	infolist  []INFO //成功的主机列表
	allcount  int    //IP*PORT的总数量
	donecount int    //线程技术的数量
	outputMux sync.Mutex
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

	nocrack, _ := cmd.Flags().GetBool("nocrack") //弱口令扫描
	if nocrack {
		Carck = false
	} else {
		Carck = true
	}

	random, _ = cmd.Flags().GetBool("random")
	save, _ = cmd.Flags().GetBool("save")

	scanPort()

}
