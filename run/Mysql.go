package run

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golin/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unsafe"
)

// MySQL默认用户，来源官网：https://dev.mysql.com/doc/refman/8.0/en/reserved-accounts.html
var defaultuser = []string{"root", "mysql.session", "mysql.sys", "mysql.infoschema"}

type VariablGlobal struct {
	Key   string `gorm:"column:Variable_name"`
	Value string `gorm:"column:Value"`
}

type Userslist struct {
	User                  string
	Host                  string
	AuthenticationString  string
	Plugin                string
	Ssl_type              string
	Account_locked        string
	Password_lifetime     string
	Password_expired      string
	Password_last_changed string
}

var sqlcmd string //自定义sqlcmd命令
var echo bool

func Mysql(cmd *cobra.Command, args []string) {
	//分隔符
	spr, err := cmd.Flags().GetString("spript")
	if err != nil {
		fmt.Println(err)
		return
	}
	//是否是执行自定义sql命令
	sqlcmd, _ = cmd.Flags().GetString("cmd")
	//是否输出结果
	echo, _ = cmd.Flags().GetBool("echo")

	//如果value值不为空则是运行一次的模式
	value, err := cmd.Flags().GetString("value")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(value) > 10 {
		Onlyonerun(value, spr, "Mysql")
		wg.Wait()
		zlog.Info("	单次运行Mysql模式完成！")
		return
	}
	//到这是运行批量采集的
	ippath, err := cmd.Flags().GetString("ip")
	if err != nil {
		fmt.Println(err)
		return
	}
	//判断Mysql.txt文件是否存在
	Checkfile(ippath, fmt.Sprintf("名称%sip%s用户%s密码%s端口", Split, Split, Split, Split), global.FilePer, ippath)
	// 运行share文件中的函数
	Rangefile(ippath, spr, "Mysql")
	wg.Wait()
	//完成前最后写入文件
	Deffile("Mysql", count, count-len(errhost), errhost)
}

// RunMysql 执行sql语句
func RunMysql(myname string, myuser string, mypasswd string, myhost string, myport string) {
	defer wg.Done()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=1.5s", myuser, mypasswd, myhost, myport, "mysql")
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		errhost = append(errhost, myhost)
		return
	}
	//确认采集完成目录是否存在
	fullPath := filepath.Join(succpath, "MySQL")
	_, err = os.Stat(fullPath)
	if os.IsNotExist(err) {
		os.MkdirAll(fullPath, os.FileMode(global.FilePer))
	}
	var fire string
	if sqlcmd != "" {
		fire = filepath.Join(fullPath, fmt.Sprintf("%s_%s(%s).log", myname, myhost, sqlcmd))
	} else {
		fire = filepath.Join(fullPath, fmt.Sprintf("%s_%s.log", myname, myhost))
	}
	defer echosqlfie(echo, fire)
	//先删除之前的同名记录文件
	os.Remove(fire)
	file, err := os.OpenFile(fire, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.FileMode(global.FilePer))
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		errhost = append(errhost, myhost)
		return
	}
	write := bufio.NewWriter(file)

	//判读是否是自定义sql命令
	if sqlcmd != "" {
		rows, _ := db.Raw(sqlcmd).Rows()
		write.WriteString("----------------" + myhost + "执行sql命令：" + sqlcmd + "\n\n")
		//fmt.Println()

		for _, v := range scanRows2map(rows) {
			for k, vv := range v {
				write.WriteString(k + ":	" + vv + "\n")
			}
			write.WriteString("--------------------------------\n")
		}
		write.Flush()
		return
	}

	//版本信息
	write.WriteString("版本信息:")
	rows, _ := db.Raw("select version()").Rows()
	var version string
	for rows.Next() {
		rows.Scan(&version)
	}
	write.WriteString(version + "  ")
	//此次连接ID可与查询日志关联
	write.WriteString("本次连接ID:")
	rows, _ = db.Raw("select connection_id()").Rows()
	var CONNECTION_ID string
	for rows.Next() {
		rows.Scan(&CONNECTION_ID)
	}
	write.WriteString(CONNECTION_ID + "  (MySQL服务器会为每个客户端连接分配一个唯一的连接ID，在执行SQL查询、更新、删除等操作时，可以使用该连接ID来标识当前的连接。如果开启general_log后可基于次ID查询记录。)\n\n")
	//用户信息
	write.WriteString("------------------------------查询用户相关：身份鉴别、访问控制基本都基于以下内容\n")
	var userlist []Userslist //用户列表信息
	db.Raw(`select user,host,authentication_string,plugin,ssl_type,account_locked,password_lifetime,password_expired,password_last_changed from mysql.user`).Scan(&userlist)
	for _, v := range userlist {
		//是否为默认账户
		if checkdefultuser(v.User) {
			write.WriteString("------------------------------用户: " + v.User + "(此账号为默认账户)\n")
		} else {
			write.WriteString("------------------------------用户: " + v.User + "\n")
		}
		write.WriteString("连接方式：" + v.Host + "\n")
		write.WriteString("密码信息：" + v.AuthenticationString + "\n")
		write.WriteString("密码加密插件：" + v.Plugin + "	(mysql_native_password采用SHA1+盐值机制；caching_sha2_password采用SHA256+加盐机制)\n")
		write.WriteString("加密连接类型：" + v.Ssl_type + "\n")
		write.WriteString("是否锁定：" + v.Account_locked + "\n")
		write.WriteString("过期时间：" + v.Password_lifetime + "\n")
		write.WriteString("是否过期：" + v.Password_expired + "\n")
		write.WriteString("上次修改密码时间：" + v.Password_last_changed + "\n")
		//查询权限
		usergrant := fmt.Sprintf("SHOW GRANTS FOR '%s'@'%s'", v.User, v.Host)
		rows, _ := db.Raw(usergrant).Rows()
		var grant, grantdata string
		for rows.Next() {
			rows.Scan(&grant)
			grantdata = grantdata + grant + "、"
		}
		write.WriteString("权限相关：" + grantdata + "\n")
		if strings.Contains(grantdata, "GRANT ALL PRIVILEGES ON *.*") || strings.Contains(grantdata, "GRANT SUPER") {
			write.WriteString("账户类别：超级管理员(人工基于上述权限自己判断,仅供参考)\n")
		} else {
			write.WriteString("账户类别：(人工基于上述权限自己判断,仅供参考)\n")
		}
		write.WriteString("\n\n")
	}

	//全局密码复杂度
	var variables []VariablGlobal
	db.Raw(`show global variables like "validate_password%"`).Scan(&variables)
	write.WriteString("------------------------------全局密码复策略：\n")
	if len(variables) == 7 {
		write.WriteString(variables[0].Key + ": " + variables[0].Value + "	(是否允许密码与账户同名)\n")
		write.WriteString(variables[2].Key + ": " + variables[2].Value + "	(密码长度要求)\n")
		write.WriteString(variables[3].Key + ": " + variables[3].Value + "	(大写字符长度要求)\n")
		write.WriteString(variables[4].Key + ": " + variables[4].Value + "	(数字字符长度要求)\n")
		write.WriteString(variables[6].Key + ": " + variables[6].Value + "	(特殊字符长度要求)\n")
		write.WriteString(variables[5].Key + ": " + variables[5].Value + "	(0/LOW：只检查长度。1/MEDIUM：检查长度、数字、大小写、特殊字符。2/STRONG：检查长度、数字、大小写、特殊字符、字典文件。)\n")
	}

	//全局密码过期时间
	write.WriteString("\n------------------------------密码过期时间：此变量定义全局自动密码过期策略，默认值为 0，即禁用自动密码过期。如果的值为正整数N，则表示允许的密码生存期，必须M天更改密码。密码期限是从其最近一次密码更改的日期和时间开始评估的。）\n")
	db.Raw(`show variables like 'default_password_lifetime'`).Scan(&variables)
	if len(variables) == 1 {
		write.WriteString(variables[0].Key + ": " + variables[0].Value + "	(过期天数)\n")
	}
	//查询失败锁定次数
	write.WriteString("\n------------------------------失败锁定策略:\n")
	db.Raw(`show variables like '%connection_control%'`).Scan(&variables)
	if len(variables) == 3 {
		write.WriteString(variables[0].Key + ": " + variables[0].Value + "	(在服务器为后续连接尝试添加延迟之前允许帐户连续失败的连接尝试次数)\n")
		write.WriteString(variables[1].Key + ": " + variables[1].Value + "	(超过阈值的连接失败的最大延迟,以毫秒为单位)\n")
		write.WriteString(variables[1].Key + ": " + variables[1].Value + "	(超过阈值的连接失败的最小延迟,以毫秒为单位)\n")
	}

	//查询超时功能，不生效，默认不符合
	write.WriteString("\n------------------------------超时策略(即使在存在超时情况下会提示重新连接然后正常执行命令返回结果，无法达到用户退出。):\n")
	db.Raw(`show global variables like 'connect_timeout'`).Scan(&variables)
	if len(variables) == 1 {
		write.WriteString(variables[0].Key + ": " + variables[0].Value + "	(登录时连接超时，单位为秒，配置后不生效，联盟也提出此问题。默认不符合)\n")
	}
	db.Raw(`show global variables like 'wait_timeout'`).Scan(&variables)
	if len(variables) == 1 {
		write.WriteString(variables[0].Key + ": " + variables[0].Value + "	(登录后连接超时，单位为秒，配置后不生效，联盟也提出此问题。默认不符合)\n")
	}

	//日志相关
	write.WriteString("\n------------------------------日志相关\n")
	write.WriteString("错误日志:\n")
	db.Raw(`show variables like 'log_error'`).Scan(&variables)
	if len(variables) == 1 {
		write.WriteString(variables[0].Key + ": " + variables[0].Value + "	(错误日志路径)\n")
	}
	write.WriteString("\n查询日志:\n")
	db.Raw(`show variables like 'general_log%'`).Scan(&variables)
	if len(variables) == 2 {
		write.WriteString(variables[0].Key + ": " + variables[0].Value + "	(查询日志开启状态)\n")
		write.WriteString(variables[1].Key + ": " + variables[1].Value + "	(查询日志存储路径)\n")
	}
	db.Raw(`show variables like 'log_output'`).Scan(&variables)
	if len(variables) == 1 {
		write.WriteString(variables[0].Key + ": " + variables[0].Value + "	(查询日志存放方式：FILE为文件（默认）TABLE为数据表)\n")
	}

	//最后写入文件
	err = write.Flush()
	if err != nil {
		zlog.Warn("记录保存结果失败！", zap.String("IP：", myhost))
		errhost = append(errhost, myhost)
		return
	}
}

// scanRows2map 执行自定义命令
func scanRows2map(rows *sql.Rows) []map[string]string {
	res := make([]map[string]string, 0)               //  定义结果 map
	colTypes, _ := rows.ColumnTypes()                 // 列信息
	var rowParam = make([]interface{}, len(colTypes)) // 传入到 rows.Scan 的参数 数组
	var rowValue = make([]interface{}, len(colTypes)) // 接收数据一行列的数组

	for i, colType := range colTypes {
		rowValue[i] = reflect.New(colType.ScanType())           // 跟据数据库参数类型，创建默认值 和类型
		rowParam[i] = reflect.ValueOf(&rowValue[i]).Interface() // 跟据接收的数据的类型反射出值的地址

	}
	// 遍历
	for rows.Next() {
		rows.Scan(rowParam...) // 赋值到 rowValue 中
		record := make(map[string]string)
		for i, colType := range colTypes {

			if rowValue[i] == nil {
				record[colType.Name()] = ""
			} else {
				record[colType.Name()] = Byte2Str(rowValue[i].([]byte))
			}
		}
		res = append(res, record)
	}
	return res
}

// []byte to string
func Byte2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func echosqlfie(stratic bool, path string) {
	if stratic {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(content))
	}

}

// checkdefultuser 严重是否为MySQL默认账户
func checkdefultuser(user string) bool {
	for _, s := range defaultuser {
		if user == s {
			return true
		}
	}
	return false

}
