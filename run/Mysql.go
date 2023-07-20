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
var defaultuser = []string{"root", "mysql.session", "mysql.sys", "mysql.infoschema", "debian-sys-maint"}

// VariablGlobal 系统变量
type VariablGlobal struct {
	Key   string `gorm:"column:Variable_name"`
	Value string `gorm:"column:Value"`
}

// Plugins 安装的插件
type Plugins struct {
	Name    string `gorm:"column:Name"`
	Status  string `gorm:"column:Status"`
	Type    string `gorm:"column:Type"`
	Library string `gorm:"column:Library"`
	License string `gorm:"column:License"`
}

// Userslist mysql.user表信息
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

// Database 数据库信息
type DatabaseSize struct {
	Database   string `gorm:"column:Database"`
	Size       string `gorm:"column:Size"`
	TableCount string `gorm:"column:TableCount"`
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
		fire = filepath.Join(fullPath, fmt.Sprintf("%s_%s.html", myname, myhost))
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
		for _, v := range scanRows2map(rows) {
			for k, vv := range v {
				write.WriteString(k + ":	" + vv + "\n")
			}
			write.WriteString("--------------------------------\n")
		}
		write.Flush()
		return
	}
	html := mysqlhtml()
	echoinfo := ""

	//版本信息
	rows, _ := db.Raw("select version()").Rows()
	var version string
	for rows.Next() {
		rows.Scan(&version)
	}
	echoinfo += fmt.Sprintf("<tr><td>%s</td>", version)

	//数据存储目录
	var variables []VariablGlobal
	db.Raw(`SHOW VARIABLES LIKE 'datadir'`).Scan(&variables)
	echoinfo += fmt.Sprintf("<td>%s</td>", variables[0].Value)

	//此次连接ID可与查询日志关联
	rows, _ = db.Raw("select connection_id()").Rows()
	var CONNECTION_ID string
	for rows.Next() {
		rows.Scan(&CONNECTION_ID)
	}
	echoinfo += fmt.Sprintf("<td>%s</td></tr>", CONNECTION_ID)
	html = strings.ReplaceAll(html, "版本详细信息", echoinfo)

	//数据库信息
	echoinfo = ""
	var datas []DatabaseSize
	db.Raw("SELECT TABLE_SCHEMA AS `Database`, COUNT(TABLE_NAME) AS `TableCount`, ROUND(SUM(DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024, 2) AS `Size` FROM INFORMATION_SCHEMA.TABLES GROUP BY TABLE_SCHEMA").Scan(&datas)
	for _, v := range datas {
		echoinfo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", v.Database, v.TableCount, v.Size)
	}
	html = strings.ReplaceAll(html, "数据库信息结果", echoinfo)

	//用户信息
	echoinfo = ""
	var userlist []Userslist //用户列表信息
	db.Raw(`select user,host,authentication_string,plugin,ssl_type,account_locked,password_lifetime,password_expired,password_last_changed from mysql.user`).Scan(&userlist)
	for _, v := range userlist {
		//查询权限
		usergrant := fmt.Sprintf("SHOW GRANTS FOR '%s'@'%s'", v.User, v.Host)
		rows, _ := db.Raw(usergrant).Rows()
		var grant, grantdata string
		for rows.Next() {
			rows.Scan(&grant)
			grantdata = grantdata + grant + "<br>"
		}
		if strings.Count(grantdata, "<br>") == 1 {
			grantdata = strings.ReplaceAll(grantdata, "<br>", "")
		}

		//是否为默认账户
		if checkdefultuser(v.User) {
			v.User = v.User + "(此账号为默认账户)"
		}
		super := "业务账户"
		if strings.Contains(grantdata, "GRANT ALL PRIVILEGES ON *.*") || strings.Contains(grantdata, "GRANT SUPER") {
			super = "超级管理员"
		}
		if v.User == "root(此账号为默认账户)" {
			super = "超级管理员"
		}
		echoinfo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			v.User, v.Host, v.AuthenticationString, v.Plugin, v.Ssl_type, v.Account_locked, v.Password_lifetime, v.Password_expired, v.Password_last_changed, super, grantdata)
	}
	html = strings.ReplaceAll(html, "用户详细信息", echoinfo)

	//全局密码复杂度
	db.Raw(`show global variables like "validate_password%"`).Scan(&variables)
	echoinfo = ""
	if len(variables) == 7 {
		echoinfo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", variables[0].Value, variables[2].Value, variables[3].Value, variables[4].Value, variables[6].Value, variables[5].Value)
	} else {
		echoinfo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", "", "", "", "", "", "")
	}
	html = strings.ReplaceAll(html, "密码复杂度详细信息", echoinfo)

	//全局密码过期时间
	echoinfo = ""
	db.Raw(`show variables like 'default_password_lifetime'`).Scan(&variables)
	if len(variables) == 1 {
		html = strings.ReplaceAll(html, "密码过期时间结果", variables[0].Value)
	}
	//查询失败锁定次数
	echoinfo = ""
	write.WriteString("\n------------------------------失败锁定策略:\n")
	db.Raw(`show variables like '%connection_control%'`).Scan(&variables)
	if len(variables) == 3 {
		echoinfo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", variables[0].Value, variables[1].Value, variables[2].Value)
	}
	html = strings.ReplaceAll(html, "密码过期时间详细信息", echoinfo)

	//查询超时功能，不生效，默认不符合
	echoinfo = ""
	db.Raw(`show global variables like 'connect_timeout'`).Scan(&variables)
	if len(variables) == 1 {
		echoinfo += fmt.Sprintf("<tr><td>%s</td>", variables[0].Value)
	}
	db.Raw(`show global variables like 'wait_timeout'`).Scan(&variables)
	if len(variables) == 1 {
		echoinfo += fmt.Sprintf("<td>%s</td></tr>", variables[0].Value)
	}
	html = strings.ReplaceAll(html, "超时策略详细信息", echoinfo)

	//日志相关
	echoinfo = ""
	//错误日志
	db.Raw(`show variables like 'log_error'`).Scan(&variables)
	if len(variables) == 1 {
		echoinfo += fmt.Sprintf("<tr><td>%s</td>", variables[0].Value)
	}
	//查询日志
	db.Raw(`show variables like 'general_log%'`).Scan(&variables)
	if len(variables) == 2 {
		echoinfo += fmt.Sprintf("<td>%s</td>", variables[0].Value)
	}
	//查询日志目录
	db.Raw(`show variables like 'general_log_file'`).Scan(&variables)
	if len(variables) == 1 {
		echoinfo += fmt.Sprintf("<td>%s</td>", variables[0].Value)
	} else {
		echoinfo += "<td>%s</td>"
	}
	//二进制文件状态
	db.Raw(`show variables like 'log_bin'`).Scan(&variables)
	if len(variables) == 1 {
		echoinfo += fmt.Sprintf("<td>%s</td>", variables[0].Value)
	} else {
		echoinfo += "<td>%s</td>"
	}
	//慢日志文件状态
	db.Raw(`show variables like 'slow_query_log'`).Scan(&variables)
	if len(variables) == 1 {
		echoinfo += fmt.Sprintf("<td>%s</td></tr>", variables[0].Value)
	} else {
		echoinfo += "<td>%s</td></tr>"
	}

	html = strings.ReplaceAll(html, "日志相关详细信息", echoinfo)

	//插件
	echoinfo = ""
	var plugs []Plugins
	db.Raw(`show plugins`).Scan(&plugs)
	for _, plug := range plugs {
		echoinfo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", plug.Name, plug.Status, plug.Type, plug.Library, plug.License)
	}
	html = strings.ReplaceAll(html, "插件信息详细信息", echoinfo)

	//全局配置变量
	echoinfo = ""
	db.Raw(`show global variables`).Scan(&variables)
	for i := 0; i < len(variables); i++ {
		echoinfo += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", variables[i].Key, variables[i].Value)
	}
	html = strings.ReplaceAll(html, "系统变量详细信息", echoinfo)

	//全局状态变量
	echoinfo = ""
	db.Raw(`show global status`).Scan(&variables)
	for i := 0; i < len(variables); i++ {
		echoinfo += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", variables[i].Key, variables[i].Value)
	}
	html = strings.ReplaceAll(html, "状态变量详细信息", echoinfo)

	html = strings.ReplaceAll(html, "替换名称", myhost)

	err = os.WriteFile(fire, []byte(html), os.FileMode(global.FilePer))
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
