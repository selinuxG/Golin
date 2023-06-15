package crack

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"runtime"
	"sync"
)

var (
	ch = make(chan struct{}, 30)
	wg = sync.WaitGroup{}
)

func Run(cmd *cobra.Command, args []string) {
	info := parseFlags(cmd)
	checkWindow(info)
	fmt.Printf("[*] 运行弱口令检测模式:%s,主机:%d个,尝试用户%d个,尝试密码%d个,主机需单独%d次尝试,共计尝试%d次,超时等待:%d/s 线程数:%d 祝君好运:)\n", info.Mode, len(info.IP), len(info.User), len(info.Passwd), len(info.User)*len(info.Passwd), len(info.IP)*len(info.User)*len(info.Passwd), info.Timeout, info.Chan)
	ch = make(chan struct{}, info.Chan)

	newport := info.Prot
	for _, ip := range info.IP {
		value, ok := fireip[ip] //先从map中确认是否有自定义的端口
		if ok {
			newport = value
		}
		if !info.NoPing { //true为跳过ping
			if !NetWorkStatus(ip) {
				clearLine(ip, "", "", newport, "warning")
				continue
			}
			if !IsPortOpen(ip, info.Prot, info.Timeout, info.Mode) {
				clearLine(ip, "", "", newport, "warning")
				continue
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel() //确保所有的goroutine都已经退出

		for _, user := range info.User {
			for _, passwd := range info.Passwd {
				clearLine(ip, user, passwd, newport, "start")
				ch <- struct{}{}
				wg.Add(1)
				switch info.Mode {
				case "ssh":
					go SSH(ctx, cancel, ip, user, passwd, newport, info.Timeout)
				case "mysql":
					go mySql(ctx, cancel, ip, user, passwd, newport, info.Timeout)
				case "redis":
					go rediscon(ip, user, passwd, newport, info.Timeout) //因为redis特征可以一个用户多密码，所以让字典跑完
				case "pgsql":
					go pgsql(ctx, cancel, ip, user, passwd, newport, info.Timeout)
				case "sqlserver":
					go sqlservercon(ctx, cancel, ip, user, passwd, newport, info.Timeout)
				case "ftp":
					go ftpcon(ip, user, passwd, newport, info.Timeout) //因为ftp匿名账户的问题，所以让字典跑完
				case "rdp":
					go rdpcon(ctx, cancel, ip, user, passwd, newport, info.Timeout)
				case "smb":
					go smbcon(ip, user, passwd, newport, info.Timeout) //因为smb匿名账户的问题，所以让字典跑完
				case "telnet":
					go telnetcon(ctx, cancel, ip, user, passwd, newport, info.Timeout) //因为smb匿名账户的问题，所以让字典跑完
				case "tomcat":
					go tomcat(ctx, cancel, ip, user, passwd, newport, info.Timeout)
				}
			}
		}
		fmt.Print("\033[2K") // 擦除整行
		fmt.Print("\r")      // 光标移动到行首
	}
	wg.Wait()
}

// end 结束
func end(ip, user, passwd string, port int) {
	fmt.Print("\033[2K") // 擦除整行
	fmt.Printf("\r[*] 发现口令！IP:%s  用户:%s  密码:%s  端口:%d  \n", ip, user, passwd, port)
}

// clearLine 清除上行并显示新数据
func clearLine(host, user, passwd string, port int, mode string) {
	fmt.Print("\033[2K") // 擦除整行
	fmt.Print("\r")      // 光标移动到行首
	if mode == "start" {
		fmt.Printf("\r[-] 开始尝试主机：%s 端口：%d 用户：%s 密码：%s\r", host, port, user, passwd)
		return
	}
	fmt.Printf("[-] 探测主机：%s:%d 网络｜端口不可达 or 端口运行协议与探测模式不匹配 跳过！\n", host, port)

}

// checkWindow 检查特定模式是否仅运行在特地系统下运行
func checkWindow(info INFO) {
	if info.Mode == "rdp" {
		if runtime.GOOS != "windows" {
			fmt.Printf("[-] 当前操作系统为%s,此模式仅允许运行在Windows操作系统下! 拜拜！\n", runtime.GOOS)
			os.Exit(0)
		}
	}

}
