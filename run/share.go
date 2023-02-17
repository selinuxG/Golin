package run

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	succpath = "采集完成目录"        //保存采集目录
	pem      = 755             //创建文件、目录时的权限
	Split    = "~"             //默认分隔符
	DeFfile  = "Golin运行记录.log" //程序运行记录文件
)

var (
	count        int                                                    //总数量,多少行文件就是多少
	wg           sync.WaitGroup                                         //线程
	errhost      []string                                               //失败主机列表
	runcmd       = Linux_cmd()                                          //运行的linux默认cmd命令
	denynametype = []string{"\\", "\\/", "*", "?", "\"", "<", ">", "|"} //windos下不允许创建名称的特殊符号。
)

// Rangefile 遍历文件并创建线程 path=模式目录 spr=按照什么分割 runtype运行类型
func Rangefile(path string, spr string, runtype string) {
	fire, _ := ioutil.ReadFile(path)
	lines := strings.Split(string(fire), "\n")
	wg.Add(len(lines))
	//count += len(lines)
	for i := 0; i < len(lines); i++ {
		firecount := strings.Count(lines[i], spr)
		if firecount != 4 {
			wg.Done()
			continue
		}

		linedata := lines[i]
		Name := strings.Split(string(linedata), spr)[0]
		Host := strings.Split(string(linedata), spr)[1]
		User := strings.Split(string(linedata), spr)[2]
		Passwrod := strings.Split(string(linedata), spr)[3]
		Port1 := strings.Split(string(linedata), spr)[4]
		//windos中换行符可能存在为/r/n,之前分割/n,还留存/r,清除它
		Porttmp := strings.Replace(Port1, "\r", "", -1)
		Port, err := strconv.Atoi(Porttmp)
		if err != nil {
			wg.Done()
			errhost = append(errhost, Host)
			continue
		}
		//判断host是不是正确的IP地址格式
		address := net.ParseIP(Host)
		if address == nil {
			wg.Done()
			continue
		}
		//判断端口范围是否是1-65535
		if Port == 0 || Port > 65535 {
			wg.Done()
			continue
		}
		//总数量+1
		count += 1
		//如果是Windows先判断保存文件是否存在特殊字符,是的话不执行直接记录为失败主机
		if runtime.GOOS == "windows" {
			if InSlice(denynametype, Name) {
				wg.Done()
				errhost = append(errhost, Host)
				continue
			}
		}
		fmt.Printf("\u001B[%dm✔‍ 开启线程 %s_%s \x1b[0m\n", 34, Name, Host)
		switch runtype {
		case "Linux":
			go Runssh(Name, Host, User, Passwrod, Port, runcmd)
		case "Mysql":
			go RunMysql(Name, User, Passwrod, Host, strconv.Itoa(Port))
		case "Redis":
			go Runredis(Name, Host, Passwrod, strconv.Itoa(Port))
		}
	}
}

// Onlyonerun 只允许一次的模式
func Onlyonerun(value string, spr string, runtype string) {
	firecount := strings.Count(value, spr)
	if firecount != 4 {
		fmt.Printf("\x1b[%dm错误🤷‍ 格式不正确！ \x1b[0m\n", 31)
		return
	}
	Name := strings.Split(value, spr)[0]
	Host := strings.Split(value, spr)[1]
	User := strings.Split(value, spr)[2]
	Passwrod := strings.Split(value, spr)[3]
	Port1 := strings.Split(value, spr)[4]
	Porttmp := strings.Replace(Port1, "\r", "", -1)
	Port, err := strconv.Atoi(Porttmp)
	if err != nil {
		fmt.Printf("\x1b[%dm错误‍ 端口格式转换失败,退出 \x1b[0m\n", 31)
		os.Exit(3)
	}
	address := net.ParseIP(Host)
	if address == nil {
		fmt.Printf("\x1b[%dm不是正确的IP地址,退出 \x1b[0m\n", 31)
		os.Exit(3)
	}
	//判断端口范围是否是1-65535
	if Port == 0 || Port > 65535 {
		fmt.Printf("\x1b[%dm不是正确的端口范围,退出 \x1b[0m\n", 31)
		os.Exit(3)
	}
	//如果是Windows先判断保存文件是否存在特殊字符,是的话不执行直接记录为失败主机
	if runtime.GOOS == "windows" {
		if InSlice(denynametype, Name) {
			fmt.Printf("\x1b[%dm错误:保存文件包含特殊字符,无法保存,请修改在执行。\x1b[0m\n", 31)
			os.Exit(3)
		}
	}
	switch runtype {
	case "Linux":
		wg.Add(1)
		fmt.Printf("\x1b[%dm✔‍ 开启单主机执行:Linux模式,开始采集%s！ \x1b[0m\n", 34, Host)
		go Runssh(Name, Host, User, Passwrod, Port, runcmd)
	case "Mysql":
		wg.Add(1)
		fmt.Printf("\x1b[%dm✔‍ 开启单主机执行:Mysql模式,开始采集%s！ \x1b[0m\n", 34, Host)
		go RunMysql(Name, User, Passwrod, Host, strconv.Itoa(Port))
	case "Redis":
		wg.Add(1)
		fmt.Printf("\x1b[%dm✔‍ 开启单主机执行:Redis模式,开始采集%s！ \x1b[0m\n", 34, Host)
		go Runredis(Name, Host, Passwrod, strconv.Itoa(Port))
	}
}

// Checkfile 判断某个模式下的默认文件是否存在
func Checkfile(name string, data string, pems int, path string) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		datanew := []byte(string(data))
		ioutil.WriteFile(path, datanew, fs.FileMode(pems))
		fmt.Printf("\x1b[%dm错误🤷‍ %s文件不存在！ \x1b[0m\n", 31, name)
		fmt.Printf("\x1b[%dm提示🤦‍ 已自动创建符合格式的%s,请补充后在执行吧！ \x1b[0m\n", 34, name)
		os.Exit(3)
	}
}

// Deffile 程序退出前运行的函数，用于生成日志
func Deffile(moude string, count int, success int, errhost []string) {
	path := DeFfile
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Create(path)
	}
	file, _ := os.OpenFile(DeFfile, os.O_WRONLY|os.O_APPEND, pem)
	write := bufio.NewWriter(file)
	write.WriteString("执行模式为:" + moude + "\n完成时间:" + Nowtime() + "\n采集总数量为:" + strconv.Itoa(count) + "\n成功数量为:" + strconv.Itoa(success) + "\n失败数量为:" + strconv.Itoa(count-success) + "\n")
	if count-success > 0 {
		for _, v := range errhost {
			write.WriteString("失败主机:" + v + "\n")
		}
	}
	write.WriteString("<------------------------------------------>\n")
	write.Flush()
	defer file.Close()
	return
}

// Nowtime 获取当前时间
func Nowtime() string {
	timeObj := time.Now()
	year := timeObj.Year()
	month := timeObj.Month()
	day := timeObj.Day()
	hour := timeObj.Hour()
	minute := timeObj.Minute()
	second := timeObj.Second()
	timenow := fmt.Sprintf("%d-%d-%d %d:%d:%d", year, month, day, hour, minute, second)
	return timenow
}

// InSlice 判断字符串是否在 不允许命名的slice中。
func InSlice(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
