package crack

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"golin/global"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ListCrackHost []SussCrack

type SussCrack struct {
	Host   string
	User   string
	Passwd string
	Port   int
	Mode   string
}

type HostPort struct {
	Host string
	Port int
}

var MapCrackHost = make(map[HostPort]SussCrack) //使用 Host和Port作为键，SussCrack结构体作为值。对于每个 Host 和 Port 的组合，只会存储一次弱口令信息，避免并发遇到匿名用户输出。

// ConnectionFunc 定义一个函数类型
type ConnectionFunc func(cancel context.CancelFunc, host, user, passwd string, newport, timeout int)

// connectionFuncs 创建一个映射，将字符串映射到对应的函数
var connectionFuncs = map[string]ConnectionFunc{
	"ssh":        SSH,
	"mysql":      mySql,
	"redis":      rediscon,
	"postgresql": pgsql,
	"sqlserver":  sqlservercon,
	"ftp":        ftpcon,
	"smb":        smbcon,
	"telnet":     telnetcon,
	"tomcat":     tomcat,
	"rdp":        rdpcon,
	"oracle":     oraclecon,
}

func Run(host, port string, Timeout, chanCount int, mode string) {
	ch := make(chan struct{}, chanCount)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Timeout))
	defer cancel() //确保所有的goroutine都已经退出
	newport, _ := strconv.Atoi(port)

	for _, user := range Userlist(mode) {
		for _, passwd := range Passwdlist() {
			passwd = strings.Replace(passwd, "%user%", user, -1) //支持在密码中使用%user%来代替用户名
			fmt.Printf("\033[2K\r")                              // 擦除整行
			fmt.Printf("\r%s", color.MagentaString("\r[...] 正在进行弱口令扫描 -> %s", fmt.Sprintf("%s://%s:%s?user=%s?passwd=%s", mode, host, port, user, passwd)))

			ch <- struct{}{}
			wg.Add(1)
			if connFunc, ok := connectionFuncs[mode]; ok {
				go crackOnce(ctx, cancel, host, user, passwd, newport, Timeout, ch, &wg, connFunc, mode)
			} else {
				wg.Done()
				<-ch
			}

		}
	}
	wg.Wait()
}

func end(host, user, passwd string, port int, mode string) {
	global.PrintLock.Lock()
	defer global.PrintLock.Unlock()
	//ListCrackHost = append(ListCrackHost, SussCrack{host, user, passwd, port, mode})
	MapCrackHost[HostPort{Host: host, Port: port}] = SussCrack{host, user, passwd, port, mode}

}

func done(ch <-chan struct{}, wg *sync.WaitGroup) {
	<-ch
	wg.Done()
}

func crackOnce(ctx context.Context, cancel context.CancelFunc, host, user, passwd string, newport, timeout int, ch <-chan struct{}, wg *sync.WaitGroup, connFunc ConnectionFunc, key string) {
	defer done(ch, wg)

	hasDone := make(chan struct{}, 1)
	go func() {
		connFunc(cancel, host, user, passwd, newport, timeout)
		hasDone <- struct{}{}
	}()

	select {
	case <-hasDone:
		return
	case <-ctx.Done():
		if os.Getenv("crack") == "on" {
			fmt.Println(key, host, user, passwd, "crack timeout.....done")
		}
		return
	}
}
