package crack

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"strconv"
	"sync"
)

// ConnectionFunc 定义一个函数类型
type ConnectionFunc func(ctx context.Context, cancel context.CancelFunc, host, user, passwd string, newport, timeout int, ch <-chan struct{}, wg *sync.WaitGroup)

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
	ch := make(chan struct{}, chanCount)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //确保所有的goroutine都已经退出
	newport, _ := strconv.Atoi(port)

	for _, user := range Userlist(mode) {
		for _, passwd := range Passwdlist() {
			fmt.Printf("\033[2K\r") // 擦除整行
			fmt.Printf("\r%s", color.MagentaString("\r[...] 正在进行弱口令扫描 -> %s", fmt.Sprintf("%s://%s:%s?user=%s?passwd=%s", mode, host, port, user, passwd)))

			ch <- struct{}{}
			wg.Add(1)
			if connFunc, ok := connectionFuncs[mode]; ok {
				go connFunc(ctx, cancel, host, user, passwd, newport, Timeout, ch, &wg)
			} else {
				wg.Done()
				<-ch
			}

		}
	}
	wg.Wait()
}

func end(host, user, passwd string, port int, mode string) {
	fmt.Printf("\033[2K\r") // 擦除整行
	fmt.Printf("\r| %-2s | %-15s | %-4d |%-6s|%-4s|%-50s \n",
		fmt.Sprintf("%s", color.GreenString("%s", "✓")),
		host,
		port,
		color.RedString("%s", "弱口令"),
		mode,
		fmt.Sprintf("%s", color.RedString(fmt.Sprintf("%s	%s", user, passwd))),
	)
	fmt.Printf("\033[2K\r") // 擦除整行
}
func done(ch <-chan struct{}, wg *sync.WaitGroup) {
	<-ch
	wg.Done()
}
