package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"golin/checkpass"
	"golin/dbcp"
	"golin/files_md5"
	"golin/oracle"
	"golin/osinfo"
	"golin/redis"
	"golin/windows"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/ssh"
)

var (
	//成功数量,正常执行完命令后+1
	succcount int
	//总数量,循环一次+1
	count int
	//线程
	wg sync.WaitGroup
	//开始时间
	start = time.Now()
	//接收cmd命令参数
	cmd1 = dbcp.CMD1()
	//cmd1       = `echo "系统版本:";echo;cat /etc/redhat-release;echo;echo;lsb_release -a;echo;echo "用户唯一性:";echo;cat /etc/passwd |awk -F":" '{print $1}'|uniq -c;echo;echo "特权用户:";echo;awk -F: '$3==0 {print $1}' /etc/passwd;echo;echo "空密码用户:";echo;cat /etc/shadow |awk -F: 'length($2)==0 {print $1}';echo;echo "查看用户登录信息:";last;echo;echo "查看所有用户上次登录信息";echo;lastlog;echo;echo;echo "系统修改密码天数:";echo;cat /etc/login.defs|grep -v "#" |grep -v "^$";echo;echo "用户复杂度:";echo;cat /etc/pam.d/system-auth|grep "password";echo;cat /etc/pam.d/sshd|grep "password";echo;cat /etc/pam.d/login|grep "password";echo;echo "审计服务：";echo;ps aux |grep -E "rsyslog|auditd" |grep -v "echo";echo;echo;echo;echo "audit日志位置:";echo;cat /etc/audit/auditd.conf |grep "log_file ="|grep -v "max";echo;echo "前三行审计日志";echo;head -n 3 /var/log/audit/audit.log;echo;echo "后三行审计日志:";echo;tail -n 3 /var/log/audit/audit.log;echo;echo;echo "前三行messages日志：";echo;head -n 3 /var/log/messages;echo;echo "后五行messages日志";echo;tail -n 5 /var/log/messages;echo;echo;echo "审计策略";echo;auditctl -l;echo;echo;echo "审计目录权限";echo;ls -l /var/log/audit/;echo;echo "超时时间:";echo;cat /etc/profile |grep "TMOUT" |grep -v "echo";set |grep "TMOUT"|grep -v "echo";echo $TMOUT;echo;echo "Selinux状态：";echo;sestatus;echo;echo "IP限制:";echo;cat /etc/hosts.deny |grep -v "#";echo;cat /etc/hosts.allow |grep -v "#";echo;echo "防火墙策略：";echo;iptables -L;echo;echo "防火墙状态";echo;firewall-cmd --state;echo;service ufw status;echo;echo "是否开启TELNET,FTT,MAIL高危服务：";echo;ps aux |grep -E "telnet|sshd|ftp|mail" |grep -v "echo";echo;echo "硬盘大小:";echo;df -h;echo;echo "所有开放端口";echo;netstat -anpt;echo;echo "查看syslog配置信息:";cat /etc/rsyslog.conf  |grep -v "#" |grep -v "^$"`
	cmd        = flag.String("cmd", cmd1, "此参数是自定义linux模式下执行的命令")
	run        = flag.String("run", "false", "此参数是定义采集模式,-run mysql采集mysql,-run linux采集linux，-run postgresql采集postgresql,")
	db         = flag.String("db", "false", "此参数是输出等保常用命令,现支持:oracle、mysql、linux、达梦、cisco、huawei、aix、postgresql")
	port       = flag.String("port", "false", "此参数是输出需要扫描端口的主机支持域名,如-port www.baidu.com")
	webserver  = flag.String("webserver", "false", "此参数是确认是否输出网页的的body信息,如需开启,-body true,只能是域名,不要加http、https")
	img        = flag.String("img", "false", "此参数的输出是先指定-webserver,通过-img true开启,webserver截图")
	ipinfo     = flag.String("ipinfo", "flase", "此参数的输出是获取IP地址得信息,所在地、经纬度,注意需要联网")
	fileshare  = flag.String("fileshare", "flase", "web文件共享当前目录,true开启,端口为11111")
	systeminfo = flag.String("systeminfo", "flase", "获取当前系统信息,true开启")
	ifconfig   = flag.String("ifconfig", "flase", "获取当前系统外网ip,true开启")
	checkpsswd = flag.String("checkpass", "flase", "验证密码复杂度是否合规")
	filesmd5   = flag.String("filesmd5", "flase", "通过对比文件MD5确认文件是否被更改")
	//goscan端口
	ar  = [...]int{1, 7, 9, 13, 19, 21, 22, 23, 25, 37, 42, 49, 53, 69, 79, 80, 1723, 81, 85, 105, 109, 111, 113, 123, 135, 137, 139, 143, 161, 179, 222, 264, 384, 389, 402, 407, 443, 446, 465, 500, 502, 512, 515, 523, 524, 540, 548, 554, 587, 617, 623, 689, 705, 771, 783, 873, 888, 902, 910, 912, 921, 993, 995, 998, 1000, 1024, 1030, 1035, 1090, 1098, 1103, 1128, 1129, 1158, 1199, 1211, 1220, 1234, 1241, 1300, 1311, 1352, 1433, 1435, 1440, 1494, 1521, 1530, 1533, 1581, 1582, 1604, 1720, 1723, 1755, 1811, 1900, 2000, 2001, 2049, 2082, 2083, 2100, 2103, 2121, 2199, 2207, 2222, 2323, 2362, 2375, 2380, 2381, 2525, 2533, 2598, 2601, 2604, 2638, 2809, 2947, 2967, 3000, 3037, 3050, 3057, 3128, 3200, 3217, 3273, 3299, 3306, 3311, 3312, 3389, 3460, 3500, 3628, 3632, 3690, 3780, 3790, 3817, 4000, 4322, 4433, 4444, 4445, 4659, 4679, 4848, 5000, 5038, 5040, 5051, 5060, 5061, 5093, 5168, 5247, 5250, 5351, 5353, 5355, 5400, 5405, 5432, 5433, 5498, 5520, 5521, 5554, 5555, 5560, 5580, 5601, 5631, 5632, 5666, 5800, 5814, 5900, 5910, 5920, 5984, 5986, 6000, 6050, 6060, 6070, 6080, 6082, 6101, 6106, 6112, 6262, 6379, 6405, 6502, 6504, 6542, 6660, 6661, 6667, 6905, 6988, 7001, 7021, 7071, 7080, 7144, 7181, 7210, 7443, 7510, 7579, 7580, 7700, 7770, 7777, 7778, 7787, 7800, 7801, 7879, 7902, 8000, 8001, 8008, 8014, 8020, 8023, 8028, 8030, 8080, 8082, 8087, 8090, 8095, 8161, 8180, 8205, 8222, 8300, 8303, 8333, 8400, 8443, 8444, 8503, 8800, 8812, 8834, 8880, 8888, 8890, 8899, 8901, 8903, 9000, 9002, 9060, 9080, 9081, 9084, 9090, 9099, 9100, 9111, 9152, 9200, 9390, 9391, 9443, 9495, 9809, 9815, 9855, 9999, 10001, 10008, 10050, 10051, 10080, 10098, 10162, 10202, 10203, 10443, 10616, 10628, 11000, 11099, 11211, 11234, 11333, 12174, 12203, 12221, 12345, 12397, 12401, 13364, 13500, 13838, 14330, 15200, 16102, 17185, 17200, 18881, 19300, 19810, 20010, 20031, 20034, 20101, 20111, 20171, 20222, 22222, 23472, 23791, 23943, 25000, 25025, 26000, 26122, 27000, 27017, 27888, 28222, 28784, 30000, 30718, 31001, 31099, 32764, 32913, 34205, 34443, 37718, 38080, 38292, 40007, 41025, 41080, 41523, 41524, 44334, 44818, 45230, 46823, 46824, 47001, 47002, 48899, 49152, 50000, 50004, 50013, 50500, 50504, 52302, 55553, 57772, 62078, 62514, 65535}
	arr = len(ar)
)

//主函数创建成功、失败目录，开启线程，调用ssh
func main() {

	//motd := `
	//
	//  ┏┓　　　┏┓
	//┏┛┻━━━┛┻┓
	//┃　　　　　　　┃
	//┃　　　━　　　┃
	//┃　＞　　　＜　┃
	//┃　　　　　　　┃
	//┃...　⌒　...　┃
	//┃　　　　　　　┃
	//┗━┓　　　┏━┛
	//    ┃　　　┃
	//    ┃　　　┃
	//    ┃　　　┃
	//    ┃　　　┃  神兽保佑
	//    ┃　　　┃  正常运行无bug
	//    ┃　　　┃  author：高业尚
	//    ┃　　　┗━━━┓
	//    ┃　　　　　　　┣┓
	//    ┃　　　　　　　┏┛
	//    ┗┓┓┏━┳┓┏┛
	//     ┃┫┫　┃┫┫
	//     ┗┻┛　┗┻┛`
	log.Println("------即将启动Golin")
	//log.Println(motd)
	runserver()
}

func runserver() {
	flag.Parse()
	//defer func() {
	//	if *run != "false" && *run != "windows" && *run != "oracle" {
	//		Breakgolin()
	//	}
	//}()
	switch *run {
	case "linux":
		golin()
	case "mysql":
		checkmysql()
	case "postgresql":
		checkpostsql()
		fmt.Println("看见像报错的信息别怕~只是读取的数据字段有空的,直接看文件即可")
	case "redis":
		redis.Run()
	case "windows":
		windows.Run()
	case "oracle":
		oracle.Run()
	case "xml":
		dbcp.Runxml()
	case "false":
	default:
		log.Println("-run 参数目前只支持linux、mysql、postgresql、redis、windows、oracle")
	}

	if *checkpsswd == "true" {
		checkpass.CheckPasswordLever(*run)
	}
	if *filesmd5 == "true" {
		files_md5.Run()
	}
	//文件共享
	if *fileshare == "true" {
		fileserver()
	}

	//外网ip
	if *ifconfig == "true" {
		ipip()
	}
	//系统信息
	if *systeminfo == "true" {
		osinfo.Osinfo()
	}
	//扫描IP信息
	if *ipinfo == "true" {
		ip_info(*ipinfo)
	}

	//扫描端口
	if *port == "true" {
		wg.Add(arr)
		aaa := *port
		for i := 0; i < arr; i++ {
			ip := fmt.Sprintf("%s:%d", aaa, ar[i])
			//fmt.Println(ip)
			go scn(ip)
		}
		wg.Wait()
	}

	//webserver确认
	if *webserver == "true" {
		wg.Add(arr)
		for i := 0; i < arr; i++ {
			if ar[i] != 443 {
				ip := fmt.Sprintf("http://%s:%d", *webserver, ar[i])
				go portbody(ip)
			}
			if ar[i] == 443 {
				ip := fmt.Sprintf("https://%s", *webserver)
				go portbody(ip)
			}
		}
		wg.Wait()
		winopen()
	}

	//等保命令输出参数
	switch {
	case *db == "oracle":
		dbcp.Oracle()
	case *db == "aix":
		dbcp.Aix()
	case *db == "huawei":
		dbcp.Huawei()
	case *db == "mysql":
		dbcp.Mysql()
	case *db == "linux":
		dbcp.Linux()
	case *db == "达梦":
		dbcp.Dm()
	case *db == "cisco":
		dbcp.Cisco()
	case *db == "postgresql":
		dbcp.Postsql()
	case *db == "nginx":
		dbcp.Nginx()
	case *db == "mongo":
		dbcp.Mongo()
	case *db == "false":
	default:
		log.Println("目前只支持输出oracle、aix、huawei、mysql、linux、达梦、cisco、postgresql、nginx、mongo测评命令")
	}

}

//采集linux服务器
func golin() {
	log.Printf("------即将启动采集功能：%s", *run)

	//flag.Parse()
	_, err := os.Stat("采集数量汇总.log")
	if os.IsNotExist(err) {
		os.Create("采集数量汇总.log")
	}
	fire, err := ioutil.ReadFile("ip.txt")
	if err != nil {
		fmt.Println("读取ip.txt文件错误--->")
		iptxt()
		//程序退出，状态码0表示成功,非0表示出错,程序会立刻终止，并且 defer 的函数不会被执行
		os.Exit(1)
	}
	lines := strings.Split(string(fire), "\n")
	wg.Add(len(lines))
	for i := 0; i < len(lines); i++ {
		firecount := strings.Count(lines[i], "~")
		if firecount != 4 {
			wg.Done()
			fmt.Println(lines[i], "格式错误")
			continue
		}

		a := lines[i]
		count++
		sshname := strings.Split(string(a), "~")[0]
		sshHost := strings.Split(string(a), "~")[1]
		sshUser := strings.Split(string(a), "~")[2]
		sshPasswrod := strings.Split(string(a), "~")[3]
		Port1 := strings.Split(string(a), "~")[4]
		//windos中换行符可能存在为/r/n,之前分割/n,还留存/r,清除它
		Port := strings.Replace(Port1, "\r", "", -1)
		sshPort, err := strconv.Atoi(Port)
		if err != nil {
			fmt.Println("转换失败--->", err)
		}
		fmt.Println("开启线程--->", count, nowtime(), "---->", "服务器名称:", sshname, "---->", sshHost)
		if *cmd == cmd1 {
			go cmd_ssh(sshname, sshHost, sshUser, sshPasswrod, sshPort, *cmd)
		} else {
			go cmd_ssh(sshname, sshHost, sshUser, sshPasswrod, sshPort, Read())
		}
		/* go func() {
		 	cmd_ssh(sshname, sshHost, sshUser, sshPasswrod, sshPort, Read())
		}()	*/

	}
	wg.Wait()
	defer et()

}

//通过调用ssh协议执行命令，写入到文件,并减一个线程数
func cmd_ssh(sshname string, sshHost string, sshUser string, sshPasswrod string, sshPort int, cmd string) {
	defer wg.Done()
	sshType := "password"
	// 创建ssh登录配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second, // ssh连接time out时间一秒钟,如果ssh验证错误会在一秒钟返回
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// config.KeyExchanges = append(config.KeyExchanges, "aes128-cbc")
	if sshType == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(sshPasswrod)}
	} else {
		fmt.Println("密码错误")
		return
	}
	// dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		errlog()
		errlog := nowtime() + "---->" + sshname + "---->" + sshHost + "------->连接失败\n"
		file, _ := os.OpenFile("采集数量汇总.log", os.O_WRONLY|os.O_APPEND, 0666)
		write := bufio.NewWriter(file)
		write.WriteString(errlog)
		write.Flush()
		defer file.Close()
		return
	}
	defer sshClient.Close()

	// 创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		fmt.Println("创建ssh session失败", err)
		return
	}

	defer session.Close()
	// 执行远程命令
	combo, err := session.CombinedOutput(cmd)
	if err != nil {
		errlog()
		errlog := nowtime() + "------>" + sshname + ":" + sshHost + "------->连接失败\n"
		file, _ := os.OpenFile("采集数量汇总.log", os.O_WRONLY|os.O_APPEND, 0666)
		write := bufio.NewWriter(file)
		write.WriteString(errlog)
		write.Flush()
		fmt.Println("远程执行cmd失败", err, "----->", sshname, "---->", sshHost)
		defer file.Close()
		return
	}
	succcount++
	//追加写出文件
	succlog()
	//timeUnix := time.Now().Unix() //已知的时间戳
	//formatTimeStr := time.Unix(timeUnix, 0).Format("15_04_05")
	fire := "采集完成目录//" + sshname + "---linux.log"
	datanew := []byte(string(combo))
	ioutil.WriteFile(fire, datanew, 0666)
	time.Sleep(time.Second * 1)

}

//检查postsql文件是否存在
func checkpostsql() {
	log.Printf("------即将启动采集功能：%s", *run)
	_, err := os.Stat("postgresql.txt")
	if os.IsNotExist(err) {
		file := "postgresql.txt"
		var port string = "5432"
		data := "数据库名称~127.0.0.1~user~password~" + port
		datanew := []byte(data)
		ioutil.WriteFile(file, datanew, 0600)
		pwd, _ := os.Getwd()
		echo := "创建读取服务器信息文件成功-->目录-->" + pwd
		fmt.Println(echo)
		os.Exit(1)
	}
	postsql()
}

//分析postsql文件
func postsql() {
	fire, _ := ioutil.ReadFile("postgresql.txt")
	lines := strings.Split(string(fire), "\n")
	for i := 0; i < len(lines); i++ {
		firecount := strings.Count(lines[i], "~")
		if firecount != 4 {
			fmt.Println(lines[i], "格式错误")
			continue
		}
		a := lines[i]
		count++
		myname := strings.Split(string(a), "~")[0]
		myhost := strings.Split(string(a), "~")[1]
		myuser := strings.Split(string(a), "~")[2]
		mypasswd := strings.Split(string(a), "~")[3]
		myport1 := strings.Split(string(a), "~")[4]
		//windos中换行符可能存在为/r/n,之前分割/n,还留存/r,清除它
		myport := strings.Replace(myport1, "\r", "", -1)
		runpostsql(myname, myuser, mypasswd, myhost, myport)

	}
	defer postsqlet()
}

//运行postsql采集
func runpostsql(myname string, myuser string, mypasswd string, myhost string, myport string) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable", myhost, myport, myuser, mypasswd)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return
		//fmt.Println("格式错误-->", myhost)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("连接失败--->", myhost, "检查用户密码端口，连接是默认数据库postgres,有吗?")
		return
	}
	succlog()
	succcount++
	fire := "采集完成目录//" + myname + "--postsql.log"
	succmylog(fire)
	file, _ := os.OpenFile(fire, os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString("------查看版本------\n")
	//查看版本
	rows, err := db.Query("select version()")
	if err != nil {
		return
	}
	var version string
	for rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(version)
	}
	write.WriteString(version)

	//查看审计模式：
	write.WriteString("\n\n------查看用户信息------\n")
	user_passwd_oid, _ := db.Query("select *from pg_shadow;	")
	var rolname, usesysid, usecreatedb, usesuper, userepl, usebypassrls, passwd, valuntil, useconfig string
	for user_passwd_oid.Next() {
		err = user_passwd_oid.Scan(&rolname, &usesysid, &usecreatedb, &usesuper, &userepl, &usebypassrls, &passwd, &valuntil, &useconfig)
		if err != nil {
			fmt.Println(err)
		}
		user_passwd_oid1 := "用户:" + rolname + "\nid:" + usesysid + "\n是否可以创建数据库" + usecreatedb + "\n是否是管理员: " + usesuper + "\n是否能开启流复制:" + userepl + "\n是否可以绕过所有安全性策略:" + usebypassrls + "\n密码字段:" + passwd + " \n口令过期时间:" + valuntil + "\n运行时配置变量的会话默认值:" + useconfig + "\n\n\n"
		write.WriteString(user_passwd_oid1)
	}

	//查看所有数据库：
	write.WriteString("\n\n------查看所有数据库------\n")
	databasesshow, _ := db.Query("SELECT datname FROM pg_database WHERE datistemplate = false;")
	var datname string
	for databasesshow.Next() {
		err = databasesshow.Scan(&datname)
		if err != nil {
			fmt.Println(err)
		}
		datname1 := "数据库名称:" + datname + "\n"
		write.WriteString(datname1)
	}

	//查看密码复杂度：
	write.WriteString("\n------查看密码复杂度------\n")
	shared_preload_libraries, _ := db.Query("show shared_preload_libraries")
	var shared string
	for shared_preload_libraries.Next() {
		err = shared_preload_libraries.Scan(&shared)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(shared)
	}
	write.WriteString(shared)

	//密码定期更换：查看valuntil字段
	write.WriteString("\n------密码过期时间-----\n")
	valuntils, _ := db.Query("select usename,valuntil from pg_shadow")
	var usename, valuntilname string
	for valuntils.Next() {
		err = valuntils.Scan(&usename, &valuntilname)
		if err != nil {
			fmt.Println(err)
		}
		val := "用户:" + usename + " 密码过期时间:" + valuntilname
		write.WriteString(val)
	}

	//查看存储加密算法
	write.WriteString("\n\n------查看存储加密算法段------\n")
	runssl, _ := db.Query("show password_encryption;")
	var password_encryption string
	for runssl.Next() {
		err = runssl.Scan(password_encryption)
		if err != nil {
			fmt.Println(err)
		}
		encryption := "加密算法:" + password_encryption
		write.WriteString(encryption)
	}

	//查看审计状态：
	write.WriteString("\n\n------查看审计状态------\n")
	postsqlaudit, _ := db.Query("show logging_collector")
	var loggin string
	for postsqlaudit.Next() {
		err = postsqlaudit.Scan(&loggin)
		if err != nil {
			fmt.Println(err)
		}
	}
	audlog := "审计状态:" + loggin
	write.WriteString(audlog)

	//查看审计模式：
	write.WriteString("\n\n------查看审计类型------\n")
	destination, _ := db.Query("show log_destination")
	var log_destination string
	for destination.Next() {
		err = destination.Scan(&log_destination)
		if err != nil {
			fmt.Println(err)
		}
	}
	logdestination := "审计模式:" + log_destination
	write.WriteString(logdestination)

	write.Flush()

}

//采集mysql入口。
func checkmysql() {
	log.Printf("------即将启动采集功能：%s", *run)
	_, err := os.Stat("mysql.txt")
	if os.IsNotExist(err) {
		file := "mysql.txt"
		var port string = "3306"
		data := "数据库名称~127.0.0.1~root~password~" + port
		datanew := []byte(data)
		ioutil.WriteFile(file, datanew, 0600)
		pwd, _ := os.Getwd()
		echo := "创建读取服务器信息文件成功-->目录-->" + pwd
		fmt.Println(echo)
		os.Exit(1)
	}
	my()
}

func my() {
	fire, _ := ioutil.ReadFile("mysql.txt")
	lines := strings.Split(string(fire), "\n")
	wg.Add(len(lines))
	for i := 0; i < len(lines); i++ {
		firecount := strings.Count(lines[i], "~")
		if firecount != 4 {
			fmt.Println(lines[i], "格式错误")
			wg.Done()
			continue
		}
		a := lines[i]
		count++
		fmt.Println(a)
		myname := strings.Split(string(a), "~")[0]
		myhost := strings.Split(string(a), "~")[1]
		myuser := strings.Split(string(a), "~")[2]
		mypasswd := strings.Split(string(a), "~")[3]
		myport1 := strings.Split(string(a), "~")[4]
		//windos中换行符可能存在为/r/n,之前分割/n,还留存/r,清除它
		myport := strings.Replace(myport1, "\r", "", -1)
		go sql_cmd(myname, myuser, mypasswd, myhost, myport)
	}
	wg.Wait()
	defer mysqlet()
}

//执行sql语句
func sql_cmd(myname string, myuser string, mypasswd string, myhost string, myport string) {
	defer wg.Done()
	dsn := myuser + ":" + mypasswd + "@tcp(" + myhost + ":" + myport + ")" + "/mysql?charset=utf8"
	//127.0.0.1:root@tcp(gys123..:3306)/mysql?charset=utf8
	//db, err := sql.Open("mysql", "root:gys123..@tcp(127.0.0.1:3306)/mysql?charset=utf8")
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		fmt.Println("格式错误", myname)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("连接失败", myname)
		errlog()
		errlog := nowtime() + "---->" + myname + "---->" + myhost + "------->连接失败\n"
		file, _ := os.OpenFile("采集数量汇总.log", os.O_WRONLY|os.O_APPEND, 0666)
		write := bufio.NewWriter(file)
		write.WriteString(errlog)
		write.Flush()
		defer file.Close()
		return
	}
	succlog()
	succcount++
	fire := "采集完成目录//" + myname + "--mysql.log"
	succmylog(fire)
	file, _ := os.OpenFile(fire, os.O_WRONLY|os.O_APPEND, 0666)
	write := bufio.NewWriter(file)
	defer file.Close()
	defer db.Close()

	//查看版本
	write.WriteString("------查看版本------\n")
	var version string
	rows := db.QueryRow("select version()")
	rows.Scan(&version)
	wriversion := "当前版本:" + version + "\n"
	write.WriteString(wriversion)

	//所有用户，登录权限
	write.WriteString("\n\n------查看所有用户------\n")
	var user, host, File_priv, Shutdown_priv, grant_priv string
	users, _ := db.Query("select user,host,File_priv,Shutdown_priv,grant_priv from user;")
	for users.Next() {
		users.Scan(&user, &host, &File_priv, &Shutdown_priv, &grant_priv)
		wriversion := "用户:" + user + "    远程权限:" + host + "   File_priv:  " + File_priv + "   Shutdown_priv  " + Shutdown_priv + "  grant_priv  " + grant_priv + "\n"
		write.WriteString(wriversion)
		fmt.Println("用户:", user, "远程权限:", host)
	}

	//超时时间
	write.WriteString("\n\n------查看超时时间------\n")
	var timeout_Variable_name, timeout_Value string
	time_conn, _ := db.Query("show variables like '%timeout%';")
	for time_conn.Next() {
		time_conn.Scan(&timeout_Variable_name, &timeout_Value)
		wriversion := timeout_Variable_name + "    字段值:" + timeout_Value + "\n"
		write.WriteString(wriversion)
		fmt.Println(" 查看超时时间:", timeout_Variable_name, "    字段值值:", timeout_Value)
	}

	//查看登录失败次数
	write.WriteString("\n\n------查看登录失败次数------\n")
	var connection_Variable_name, connection_Value string
	conne_conn, _ := db.Query(" show variables like '%connection_control%';")
	for conne_conn.Next() {
		conne_conn.Scan(&connection_Variable_name, &connection_Value)
		wriversion := connection_Variable_name + "    字段值:" + connection_Value + "\n"
		write.WriteString(wriversion)
		fmt.Println(" 查看auditd:", connection_Variable_name, "    字段值值:", connection_Value)
	}

	//密码策略
	write.WriteString("\n\n------查看密码策略------\n")
	var password_Variable_name, password_Value string
	password_conn, _ := db.Query("show variables like 'validate_password%'")
	for password_conn.Next() {
		password_conn.Scan(&password_Variable_name, &password_Value)
		wriversion := password_Variable_name + "    字段值:" + password_Value + "\n"
		write.WriteString(wriversion)
		fmt.Println(" 查看密码策略:", password_Variable_name, "    字段值值:", password_Value)
	}

	//查看开启ssl
	write.WriteString("\n\n------查看是否开启SSL------\n")
	var ssl_Variable_name, ssl_Value string
	ssl_conn, _ := db.Query(" show global variables like '%ssl%'")
	for ssl_conn.Next() {
		ssl_conn.Scan(&ssl_Variable_name, &ssl_Value)
		wriversion := ssl_Variable_name + "    字段值:" + ssl_Value + "\n"
		write.WriteString(wriversion)
		fmt.Println(" 查看SSL:", ssl_Variable_name, "    字段值值:", ssl_Value)
	}

	//查看开启审计功能
	write.WriteString("\n\n------查看是否审计------\n")
	var audit_Variable_name, audit_Value string
	auditd_conn, _ := db.Query("show global variables like '%general%'")
	for auditd_conn.Next() {
		auditd_conn.Scan(&audit_Variable_name, &audit_Value)
		wriversion := audit_Variable_name + "    字段值:" + audit_Value + "\n"
		write.WriteString(wriversion)
		fmt.Println(" 查看auditd:", audit_Variable_name, "    字段值值:", audit_Value)
	}

	write.Flush()

}

//创建采集失败存储的日志文件
func errlog() {
	path := "采集数量汇总.log"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Create(path)
	}
}

//创建采集完成存储的目录
func succlog() {
	path := "采集完成目录"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, 0666)
	}
}

//结束时执行函数，写入时间、总采集量、失败数量到文件
func et() {
	path := "采集数量汇总.log"
	// _, err := os.Stat(path)
	// if os.IsNotExist(err) {
	// 	os.Create(path)
	// }
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
	nn := "执行用户为:" + user + "\n执行模式为:" + *run + "\n完成时间:" + nowtime() + "\n总数量为:" + strconv.Itoa(count) + "\n" + "成功数量为:" + strconv.Itoa(succcount) + "\n" + "失败数量为:" + strconv.Itoa(errco) + "\n" + "-----------------------------\n"
	write.WriteString(nn)
	write.Flush()
	defer file.Close()
}

func mysqlet() {
	errlog()
	path := "采集数量汇总.log"
	// _, err := os.Stat(path)
	// if os.IsNotExist(err) {
	// 	os.Create(path)
	// }
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
	nn := "执行用户为:" + user + "\n执行模式为:" + *run + "\n完成时间:" + nowtime() + "\n总数量为:" + strconv.Itoa(count) + "\n" + "成功数量为:" + strconv.Itoa(succcount) + "\n" + "失败数量为:" + strconv.Itoa(errco) + "\n" + "-----------------------------\n"
	write.WriteString(nn)
	write.Flush()
	defer file.Close()
}

func postsqlet() {
	errlog()
	path := "采集数量汇总.log"
	// _, err := os.Stat(path)
	// if os.IsNotExist(err) {
	// 	os.Create(path)
	// }
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
	nn := "执行用户为:" + user + "\n执行模式为:" + *run + "\n完成时间:" + nowtime() + "\n总数量为:" + strconv.Itoa(count) + "\n" + "成功数量为:" + strconv.Itoa(succcount) + "\n" + "失败数量为:" + strconv.Itoa(errco) + "\n" + "-----------------------------\n"
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

func Read() string {
	f, err := ioutil.ReadFile(*cmd)
	//fmt.Println(string(f))
	if err != nil {
		fmt.Println("文件cmd读取失败")
		ipcmd()
		os.Exit(1)
	}
	return string(f)
}

func iptxt() {
	file := "ip.txt"
	data := "服务器名称~IP地址~用户名称~用户密码~端口"
	datanew := []byte(data)
	ioutil.WriteFile(file, datanew, 0600)
	pwd, _ := os.Getwd()
	echo := "创建读取服务器信息文件成功-->目录-->" + pwd
	fmt.Println(echo)

}

//创建linux，cmd命令执行文件
func ipcmd() {
	data := "pwd"
	datanew := []byte(data)
	ioutil.WriteFile(*cmd, datanew, 0600)
	pwd, _ := os.Getwd()
	echo := "创建执行服务器命令文件成功-->目录-->" + pwd
	fmt.Println(echo)
}

//判断存储mysql成功日志文件是否存在
func succmylog(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Create(path)
	}
}

func scn(ip_port string) {
	defer wg.Done()
	connTimeout := 5 * time.Second
	conn, err := net.DialTimeout("tcp", ip_port, connTimeout)
	if err != nil {
		return
	}
	fmt.Println(ip_port, ":----->open")
	conn.Close()
}

func portbody(body_port string) {
	defer wg.Done()
	//正常请求，但问题是连接超时会一直等待。
	//resp, err := http.Get(body_port)
	//设置http协议请求，超时等待10秒

	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(body_port)

	if err != nil {
		// if body_port == "https://192.168.0.121" {
		// 	fmt.Println(err, body_port)
		// }
		return

	} else {
		fmt.Println(body_port, "----webserver,runing")
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		file := "webserver/" + body_port + ".txt"
		filename1 := strings.Replace(file, "http://", "", -1)
		filename2 := strings.Replace(filename1, "https://", "", -1)
		filename := strings.Replace(filename2, ":", "--", -1)
		datanew := []byte(string(body))
		//pathname1 := strings.Replace(filename, "webserver/", "", -1)
		pathname := strings.Replace(filename, ".txt", ".png", -1)

		_, err = os.Stat("webserver")
		if err != nil {
			os.MkdirAll("webserver", 0666)
		}
		ioutil.WriteFile(filename, datanew, 0666)
		if *img == "true" {
			//fmt.Println(body_port, filename)
			webimg(body_port, pathname)
		}
		//fmt.Println(runtime.GOOS)
	}

	defer resp.Body.Close()

}

//执行-webserver后弹窗目录
func winopen() {
	if runtime.GOOS == "windows" {
		pwd, _ := os.Getwd()
		pwd += "\\webserver"
		cmd := exec.Command("explorer", pwd)
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}

// -webserver {域名}  -img true 开关、webserver网站截图
func webimg(url string, pathname string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run
	var b1 []byte
	if err := chromedp.Run(ctx,
		chromedp.Emulate(device.Reset),

		// set really large viewport
		chromedp.EmulateViewport(1920, 2000),
		chromedp.Navigate(url),
		chromedp.CaptureScreenshot(&b1),
	); err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	if err := ioutil.WriteFile(pathname, b1, 0777); err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}

}

//获取外网ip

func ipip() {
	resp, err := http.Get("https://ifconfig.me/")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

//基于腾讯的地址库信息。https://lbs.qq.com/
func ip_info(ip string) {
	log.Println("------正在调取腾讯API获取IP信息")
	url := fmt.Sprintf("https://apis.map.qq.com/ws/location/v1/ip?ip=%s&key=DDTBZ-SPYLJ-R6DF5-FIN4M-G7R5Z-RRBXZ", ip)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("亲亲~确认是否能联网呢？")

		return

	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("亲亲~确认是否能联网呢？")
		}

		fmt.Println(string(body))
		defer resp.Body.Close()
	}
}

//文件共享
func fileserver() {
	log.Println("------即将启动当前目录web共享功能")
	fmt.Println("浏览器打开当前地址,端口为:11111")
	http.Handle("/", http.FileServer(http.Dir("./"))) //把当前文件目录作为共享目录
	//如果是windos自动打开
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", "start "+"http://127.0.0.1:11111")
		cmd.Run()
	}
	http.ListenAndServe(":11111", nil)
}

func Breakgolin() {
	log.Println(`------任务完成即将退出:任务日志"采集数量汇总.log" 输出结果在"采集完成目录"中------`)
	_, err := os.Stat("采集完成目录")
	if os.IsNotExist(err) {
		return
	}
	if runtime.GOOS == "windows" {
		pwd, _ := os.Getwd()
		pwd += "\\采集完成目录"
		log.Println(pwd)
		cmd := exec.Command("explorer", pwd)
		err := cmd.Run()
		if err != nil {
			return
		}
	}
}
