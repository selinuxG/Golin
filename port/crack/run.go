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

func Run(host, port string, Timeout, chanCount int, mode string) {

	if chanCount <= 100 {
		chanCount = 100
	}
	ch = make(chan struct{}, chanCount)
	Timeout = 1 //因为已经探测过存活了，所以设置成1/s 超时

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //确保所有的goroutine都已经退出
	newport, _ := strconv.Atoi(port)

	for _, user := range Userlist(mode) {
		for _, passwd := range Passwdlist() {
			ch <- struct{}{}
			wg.Add(1)
			switch mode {
			case "ssh":
				go SSH(ctx, cancel, host, user, passwd, newport, Timeout)
			case "mysql":
				go mySql(ctx, cancel, host, user, passwd, newport, Timeout)
			case "redis":
				go rediscon(ctx, cancel, host, user, passwd, newport, Timeout)
			case "pgsql":
				go pgsql(ctx, cancel, host, user, passwd, newport, Timeout)
			case "sqlserver":
				go sqlservercon(ctx, cancel, host, user, passwd, newport, Timeout)
			case "ftp":
				go ftpcon(ctx, cancel, host, user, passwd, newport, Timeout)
			case "smb":
				go smbcon(ctx, cancel, host, user, passwd, newport, Timeout)
			case "telnet":
				go telnetcon(ctx, cancel, host, user, passwd, newport, Timeout)
			case "tomcat":
				go tomcat(ctx, cancel, host, user, passwd, newport, Timeout)
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
		color.GreenString("%s", "弱口令"),
		mode,
		fmt.Sprintf("%s	%s", user, passwd),
	)

}
