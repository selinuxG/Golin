package crack

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"sync"
)

var (
	ch = make(chan struct{}, 100)
	wg = sync.WaitGroup{}
)

// ConnectionFunc 定义一个函数类型
type ConnectionFunc func(ctx context.Context, cancel context.CancelFunc, host, user, passwd string, newport, Timeout int)

// connectionFuncs 创建一个映射，将字符串映射到对应的函数
var connectionFuncs = map[string]ConnectionFunc{
	"ssh":       SSH,
	"mysql":     mySql,
	"redis":     rediscon,
	"pgsql":     pgsql,
	"sqlserver": sqlservercon,
	"ftp":       ftpcon,
	"smb":       smbcon,
	"telnet":    telnetcon,
	"tomcat":    tomcat,
	"rdp":       rdpcon,
	"oracle":    oraclecon,
}

func Run(host, port string, Timeout, chanCount int, mode string) {
	if chanCount < 300 {
		chanCount = 300
	}
	ch = make(chan struct{}, chanCount)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //确保所有的goroutine都已经退出
	newport, _ := strconv.Atoi(port)

	for _, user := range Userlist(mode) {
		for _, passwd := range Passwdlist() {
			ch <- struct{}{}
			wg.Add(1)
			// 如果运行模式在connectionFuncs中有key值则进行弱口令扫描
			if connFunc, ok := connectionFuncs[mode]; ok {
				go connFunc(ctx, cancel, host, user, passwd, newport, Timeout)
			} else {
				wg.Done()
				<-ch
			}
		}
	}
	wg.Wait()
}

func end(host, user, passwd string, port int, mode string) {
	fmt.Printf("\r| %-2s | %-15s | %-4d |%-6s|%-4s|%s \n",
		fmt.Sprintf("%s", color.GreenString("%s", "✓")),
		host,
		port,
		color.RedString("%s", "弱口令"),
		mode,
		fmt.Sprintf("%s	%s", user, passwd),
	)

}
