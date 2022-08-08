package redis

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	start     = time.Now()   //开始时间
	count     int            //总数量
	succcount int            //成功数量
	errcount  int            //失败数量
	wg        sync.WaitGroup //线程
)

func Run() {
	log.Println("------即将启动redis测评功能")
	fire, err := ioutil.ReadFile("redis.txt")
	if err != nil {
		fmt.Println("读取redis.txt文件错误--->")
		iptxt("redis.txt")
		//程序退出，状态码0表示成功,非0表示出错,程序会立刻终止，并且 defer 的函数不会被执行
		os.Exit(1)
	}
	lines := strings.Split(string(fire), "\n")
	wg.Add(len(lines))
	for i := 0; i < len(lines); i++ {
		firecount := strings.Count(lines[i], "~")
		if firecount != 4 {
			log.Println("格式错误")
			errcount++
			wg.Done()
			continue
		} else {

			a := lines[i]
			//fmt.Println(a)
			myname := strings.Split(string(a), "~")[0]
			myhost := strings.Split(string(a), "~")[1]
			//myuser := strings.Split(string(a), "~")[2]
			mypasswd := strings.Split(string(a), "~")[3]
			myport1 := strings.Split(string(a), "~")[4]
			go initClient(myname, myhost, mypasswd, myport1)
		}
	}
	wg.Wait()
	defer et()
}

func initClient(myname, myhost, mypasswd, myport1 string) {
	count++
	defer wg.Done()
	log.Println("开启线程--->", count, nowtime(), "---->", "Redis名称:", myname, "---->", myhost)
	Port := strings.Replace(myport1, "\r", "", -1)
	ctx := context.Background()
	adr := myhost + ":" + Port
	client := redis.NewClient(&redis.Options{
		Addr:            adr,
		Password:        mypasswd,
		DB:              0,
		DialTimeout:     1 * time.Second,
		MinRetryBackoff: 1 * time.Second,
		ReadTimeout:     1 * time.Second,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		errcount++
		errlog()
		errlog := nowtime() + "---->" + myname + "---->" + myhost + "------->连接失败\n"
		file, _ := os.OpenFile("采集数量汇总.log", os.O_WRONLY|os.O_APPEND, 0666)
		write := bufio.NewWriter(file)
		write.WriteString(errlog)
		write.Flush()
		defer file.Close()
		return
	}
	succcount++
	client.Get(ctx, "config").Val()
	ipaddr := client.ConfigGet(ctx, "bind").Val()
	lofile := client.ConfigGet(ctx, "logfile").Val()
	loglevel := client.ConfigGet(ctx, "loglevel").Val()
	pass := client.ConfigGet(ctx, "requirepass").Val()
	redistimout := client.ConfigGet(ctx, "timeout").Val()
	redisport := client.ConfigGet(ctx, "port").Val()
	redisdir := client.ConfigGet(ctx, "dir").Val()
	//confinfo := client.Info(ctx).Val()
	succlog()
	fire := "采集完成目录//" + myname + "--redis.log"
	succmylog(fire)
	file, _ := os.OpenFile(fire, os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString("-----基本信息------\n")
	write.WriteString(fmt.Sprintf("地址限制策略为:%s\n", ipaddr[1]))
	write.WriteString(fmt.Sprintf("日志存储为:%s  日志等级为:%s\n", lofile[1], loglevel[1]))
	write.WriteString(fmt.Sprintf("密码信息为:%s\n", pass[1]))
	write.WriteString(fmt.Sprintf("超时时间为:%s\n", redistimout[1]))
	write.WriteString(fmt.Sprintf("redis运行端口为:%s\n", redisport[1]))
	write.WriteString(fmt.Sprintf("redis运行位置为:%s\n", redisdir[1]))
	write.WriteString("\n-----info信息------\n")
	write.WriteString(client.Info(ctx).Val())
	write.Flush()
}

func auditlog() {

}

func iptxt(file string) {
	var port string = "6379"
	data := "redis名称~1.1.1.1~~password~" + port
	datanew := []byte(data)
	ioutil.WriteFile(file, datanew, 0600)
	pwd, _ := os.Getwd()
	echo := "创建读取redis信息文件成功-->目录-->" + pwd
	fmt.Println(echo)
}

func succmylog(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Create(path)
	}
}
func succlog() {
	path := "采集完成目录"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, 0666)
	}
}

func et() {
	path := "采集数量汇总.log"
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//计算运行结束时间
	elapsed := time.Since(start)
	fmt.Println(elapsed)
	write := bufio.NewWriter(file)
	errco := count - succcount
	//获取当前用户
	username, err := user.Current()
	if err != nil {
		fmt.Println(err.Error())
	}
	user := username.Username
	nn := "执行用户为:" + user + "\n执行模式为:redis" + "\n完成时间:" + nowtime() + "\n总数量为:" + strconv.Itoa(count) + "\n" + "成功数量为:" + strconv.Itoa(succcount) + "\n" + "失败数量为:" + strconv.Itoa(errco) + "\n" + "-----------------------------\n"
	write.WriteString(nn)
	write.Flush()
	defer file.Close()
}

//获取当前时间
func nowtime() string {
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
func errlog() {
	path := "采集数量汇总.log"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Create(path)
	}
}
