package port

import (
	"fmt"
	"github.com/spf13/cobra"
	"golin/port/crack"
	"os"
	"sync"
)

var (
	port       string //接受cli要扫描的端口
	iplist     []string
	portlist   []string
	NoPing     bool                      //是否禁止ping监测
	Carck      bool                      //是否进行弱口令扫描
	Poc        bool                      //是否进行poc扫描
	ch         = make(chan struct{}, 30) //控制并发数
	wg         = sync.WaitGroup{}
	chancount  int    //并发数量
	Timeout    int    //超时等待时常
	random     bool   //打乱顺序
	infolist   []INFO //成功的主机列表
	allcount   uint32 //IP*PORT的总数量
	donecount  uint32 //线程技术的数量
	outputMux  sync.Mutex
	userfile   string //user字典路径
	passwdfile string //passwd字典路径

)

type INFO struct {
	Host     string //主机
	Port     string //开放端口
	Protocol string //协议
}

func ParseFlags(cmd *cobra.Command, args []string) {
	ipFile, _ := cmd.Flags().GetString("ipfile") //读取文件
	parseFileIP(ipFile)

	ip, _ := cmd.Flags().GetString("ip")
	if ip == "" {
		conNETLocal() //当为指定IP时通过读取本地网卡并且过滤虚拟网卡进行扫描
		if len(iplist) == 0 {
			fmt.Printf("[-] 未指定扫描主机!\n")
			os.Exit(1)
		}
	} else {
		parseIP(ip)
	}

	port, _ = cmd.Flags().GetString("port")
	parsePort(port)

	excludeport, _ := cmd.Flags().GetString("exclude") //去重端口以及排查过滤端口
	excludePort(excludeport)

	NoPing, _ = cmd.Flags().GetBool("noping")

	chancount, _ = cmd.Flags().GetInt("chan") //并发数量
	ch = make(chan struct{}, chancount)

	Timeout, _ = cmd.Flags().GetInt("time") //超时等待时常

	nocrack, _ := cmd.Flags().GetBool("nocrack") //弱口令扫描
	Carck = !nocrack

	nopoc, _ := cmd.Flags().GetBool("nopoc") //poc扫描
	Poc = !nopoc

	random, _ = cmd.Flags().GetBool("random") //打乱顺序

	userfile, _ = cmd.Flags().GetString("userfile")
	passwdfile, _ = cmd.Flags().GetString("passwdfile")
	crack.Checkdistfile(userfile, passwdfile) //先读取是否有自定义的字典文件

	scanPort()

}
