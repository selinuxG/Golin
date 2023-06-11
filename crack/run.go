package crack

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"sync"
	"time"
)

var (
	ch        = make(chan struct{}, 100)
	wg        = sync.WaitGroup{}
	starttime = time.Time{}
)

func Run(cmd *cobra.Command, args []string) {
	starttime = time.Now()
	info := parseFlags(cmd)
	fmt.Printf("[*] 运行暴力破解模式:%s,主机:%d个,尝试用户%d个,尝试密码%d个,共计%d种可能！ 线程数:%d\n", info.Mode, len(info.IP), len(info.User), len(info.Passwd), len(info.User)*len(info.Passwd), info.Chan)
	ch = make(chan struct{}, info.Chan)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //
	// 确保所有的goroutine都已经退出
	for _, ip := range info.IP {
		if info.NoPing {
			if !NetWorkStatus(ip) { //进行ping
				continue
			}
		}
		for _, user := range info.User {
			for _, passwd := range info.Passwd {
				ch <- struct{}{}
				wg.Add(1)
				switch info.Mode {
				case "ssh":
					go SSH(ctx, cancel, ip, user, passwd, info.Prot)
				case "mysql":
					go mySql(ctx, cancel, ip, user, passwd, info.Prot)
				case "redis":
					go rediscon(ip, user, passwd, info.Prot) //因为redis特征可以一个用户多密码，所以让字典跑完
				case "pgsql":
					go pgsql(ctx, cancel, ip, user, passwd, info.Prot)
				case "sqlserver":
					go sqlservercon(ctx, cancel, ip, user, passwd, info.Prot)
				case "oracle":
					go oracle(ctx, cancel, ip, user, passwd, info.Prot)
				}
			}
		}
	}
	wg.Wait()
}

// end 结束
func end(ip, user, passwd string, port int) {
	fmt.Printf("[*] 破解完成！IP:%s  用户:%s  密码:%s  端口:%d  \n", ip, user, passwd, port)
}
