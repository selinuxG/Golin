package oracle

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/cengsin/oracle"
	"golin/redis"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"os/user"
	"strconv"
	"time"
)

var (
	con       = flag.String("con", "flase", "system/oracle@1.1.2.135:1521/helowin")
	orname    = flag.String("name", "oracle", "名称")
	start     = time.Now()
	succcount int
	errco     int
)

// Version 版本
type Version struct {
	Banner string `json:"banner"`
}

// Users 所有用户状态
type Users struct {
	Username       string `json:"Username"`
	Account_status string `json:"Account_Status"`
}

// Psslockuser 查看是否具有复杂度要求并定期更换
type Psslockuser struct {
	Resource_name string `json:"Resource_Name"`
	Limit         string `json:"Limit"`
}

// Security 查看是否开启安全标记
type Security struct {
	Value string `json:"Value"`
}

// Auditlogs 审计日志
type Auditlogs struct {
	Username     string
	Timestamp    string
	Action_name  string
	Comment_text string
}

// Superuser  超级用户
type Superuser struct {
	Username string
	Sysdba   string
	Sysoper  string
	Sysasm   string
}

// Falldlogin   锁定次数
type Falldlogin struct {
	Profile string
	Limit   string
}

// Fesource   超时开关
type Fesource struct {
	Name  string
	Value string
}

// IDLETIME   超时时间
type IDLETIME struct {
	RESOURCE_NAME string
	LIMIT         string
}

func Run() {
	flag.Parse()
	log.Println("------即将启动Oracle测评功能")
	defer et()

	db, err := gorm.Open(oracle.Open(*con), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 1 * time.Millisecond,
			LogLevel:      logger.Error,
			Colorful:      true,
		}),
	})
	if err != nil {
		log.Println(err)
		errco++
		Errappend()
		return
	}
	redis.Succlog()
	fire := "采集完成目录//" + *orname + "--oracle.log"
	redis.Succmylog(fire)
	file, _ := os.OpenFile(fire, os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	write := bufio.NewWriter(file)

	var versionlist []Version
	db.Raw("select * from v$version").Scan(&versionlist)
	write.WriteString("------版本信息:------\n")
	for _, tmp := range versionlist {
		echo := fmt.Sprintf("%v\n", tmp.Banner)
		write.WriteString(echo)
	}

	//用户信息
	var userlist []Users
	write.WriteString("\n------用户信息:------\n")
	db.Raw("select username,account_status  from dba_users where profile = 'DEFAULT'").Scan(&userlist)
	for _, tmp := range userlist {
		echo := fmt.Sprintf("用户:%v	状态:%v\n", tmp.Username, tmp.Account_status)
		write.WriteString(echo)
	}

	//身份鉴别信息具有复杂度要求并定期更换
	var userpass []Psslockuser
	write.WriteString("\n密码以及等安全策略:\n")
	db.Raw("select resource_name,limit from dba_profiles where profile= 'DEFAULT'").Scan(&userpass)
	for _, tmp := range userpass {
		echo := fmt.Sprintf("字段:%v      值:%v\n", tmp.Resource_name, tmp.Limit)
		write.WriteString(echo)
	}
	//查看失败锁定次数
	var faildlogin []Falldlogin
	write.WriteString("\n------查看失败锁定次数(默认锁定一小时)------:\n")
	db.Raw(`select Profile,Limit from dba_profiles where resource_name = 'FAILED_LOGIN_ATTEMPTS'`).Scan(&faildlogin)
	for _, tmp := range faildlogin {
		echo := fmt.Sprintf("字段:%v      值:%v\n", tmp.Profile, tmp.Limit)
		write.WriteString(echo)
	}
	//查看连接超时状态
	var resource []Fesource
	write.WriteString("\n------查看连接超时是否开启------:\n")
	db.Raw(`SELECT name, value FROM gv$parameter WHERE name = 'resource_limit'`).Scan(&resource)
	for _, tmp := range resource {
		echo := fmt.Sprintf("字段:%v      状态:%v\n", tmp.Name, tmp.Value)
		write.WriteString(echo)
	}

	//查看连接超时时间
	var idletime []IDLETIME
	write.WriteString("\n------查看连接超时时间（分钟）------:\n")
	db.Raw(`select RESOURCE_NAME,LIMIT from dba_profiles where RESOURCE_NAME='IDLE_TIME'`).Scan(&idletime)
	for _, tmp := range idletime {
		echo := fmt.Sprintf("字段:%v      状态:%v\n", tmp.RESOURCE_NAME, tmp.LIMIT)
		write.WriteString(echo)
	}

	//查看是否开启安全标记
	write.WriteString("\n------安全标记状态:------\n")
	var sec []Security
	db.Raw("SELECT VALUE FROM V$OPTION WHERE PARAMETER = 'Oracle Label Security'").Scan(&sec)
	for _, tmp := range sec {
		echo := fmt.Sprintf("%v\n", tmp.Value)
		write.WriteString(echo)
	}

	//查看具备超级用户权限的用户
	var sup []Superuser
	write.WriteString("\n-----超级用户信息:------\n")
	db.Raw(`select * from V$PWFILE_USERS`).Scan(&sup)
	for _, tmp := range sup {
		echo := fmt.Sprintf("用户:%v	sysdba权限:%v	sysoper权限:%v	sysasm权限:%v\n", tmp.Username, tmp.Sysdba, tmp.Sysoper, tmp.Sysasm)
		write.WriteString(echo)
	}

	//select userid,userhost,ntimestamp#,COMMENT$TEXT FROM SYS.AUD$ WHERE ROWNUM<10;
	//查看前十条日志
	write.WriteString("\n-----前10条审计日志:------\n")
	var logs []Auditlogs
	db.Raw(`SELECT username,timestamp,action_name,comment_text FROM DBA_AUDIT_TRAIL WHERE ROWNUM<10`).Scan(&logs)
	for _, tmp := range logs {
		echo := fmt.Sprintf("用户:%v	时间:%v	事件类型:%v	事件:%v\n", tmp.Username, tmp.Timestamp, tmp.Action_name, tmp.Comment_text)
		write.WriteString(echo)
	}
	//最新的10条
	//select * from (SELECT username,timestamp,action_name,comment_text FROM DBA_AUDIT_TRAIL ORDER BY timestamp desc) where rownum<11
	write.WriteString("\n-----最新10条审计日志:------\n")
	db.Raw(`select * from (SELECT username,timestamp,action_name,comment_text FROM DBA_AUDIT_TRAIL ORDER BY timestamp desc) where rownum<11`).Scan(&logs)
	for _, tmp := range logs {
		echo := fmt.Sprintf("用户:%v	时间:%v	事件类型:%v	事件:%v\n", tmp.Username, tmp.Timestamp, tmp.Action_name, tmp.Comment_text)
		write.WriteString(echo)
	}
	write.WriteString("\n------远程加密oracle默认符合------\n")
	write.WriteString("\n------当你能用这个工具的时候很大可能是没做地址限制的，严格意义需查看listener.ora文件或访谈------\n")
	write.Flush()
	succcount++
	Errlog()
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
	//获取当前用户
	username, err := user.Current()
	if err != nil {
		fmt.Println(err.Error())
	}
	user := username.Username
	nn := "执行用户为:" + user + "\n执行模式为:oracle" + "\n完成时间:" + nowtime() + "\n总数量为:1" + "\n" + "成功数量为:" + strconv.Itoa(succcount) + "\n" + "失败数量为:" + strconv.Itoa(errco) + "\n" + "-----------------------------\n"
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

func Errappend() {
	Errlog()
	errlog := nowtime() + "---->" + *con + "---->" + "------->连接失败\n"
	file, _ := os.OpenFile("采集数量汇总.log", os.O_WRONLY|os.O_APPEND, 0666)
	write := bufio.NewWriter(file)
	write.WriteString(errlog)
	write.Flush()
	defer file.Close()
}
func Errlog() {
	path := "采集数量汇总.log"
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Create(path)
	}
}
