package run

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golin/global"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

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

func Mysql(cmd *cobra.Command, args []string) {
	//分隔符
	spr, err := cmd.Flags().GetString("spript")
	if err != nil {
		fmt.Println(err)
		return
	}
	//是否是执行自定义sql命令
	sqlcmd, _ = cmd.Flags().GetString("cmd")

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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", myuser, mypasswd, myhost, myport, "mysql")
	db, err := gorm.Open("mysql", dsn)
	defer db.Close()
	if err != nil {
		errhost = append(errhost, myhost)
		return
	}
	//确认采集完成目录是否存在
	_, err = os.Stat(succpath)
	if os.IsNotExist(err) {
		os.Mkdir(succpath, os.FileMode(global.FilePer))
	}
	fire := global.Succpath + "//" + myname + "_" + myhost + "(mysql).log"
	//先删除之前的同名记录文件
	os.Remove(fire)
	file, err := os.OpenFile(fire, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.FileMode(global.FilePer))
	defer file.Close()
	if err != nil {
		errhost = append(errhost, myhost)
		return
	}
	write := bufio.NewWriter(file)

	//判读是否是自定义sql命令
	if sqlcmd != "" {
		rows, _ := db.Raw(sqlcmd).Rows()
		write.WriteString("----------------执行sql命令：" + sqlcmd + "\n\n")
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
	//执行预制sql命令
	var userlist []Userslist //用户列表信息
	db.Raw(`select user,host,authentication_string,plugin,ssl_type,account_locked,password_lifetime,password_expired,password_last_changed from mysql.user`).Scan(&userlist)
	for _, v := range userlist {
		write.WriteString("------------------------------用户: " + v.User + "\n")
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
			write.WriteString("是否为超级管理员(最好自己基于上述权限自己判断,仅供参考)：" + "	是\n")
		} else {
			write.WriteString("是否为超级管理员(最好自己基于上述权限自己判断,仅供参考)：" + "	否\n")
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
	write.WriteString("\n------------------------------新建用户密码过期策略（MySQL5.7测试不生效，必须手动配置，所以结果以单个用户的过期时间为准）\n")
	db.Raw(`show variables like 'default_password_lifetime'`).Scan(&variables)
	if len(variables) == 1 {
		write.WriteString(variables[0].Key + ": " + variables[0].Value + "	(过期天数)\n")
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

	//版本信息
	write.WriteString("\n查询版本:	")
	rows, _ := db.Raw("select version()").Rows()
	var version string
	for rows.Next() {
		rows.Scan(&version)
	}
	write.WriteString(version + "\n")

	//最后写入文件
	err = write.Flush()
	if err != nil {
		zlog.Warn("记录保存结果失败！", zap.String("IP：", myhost))
		errhost = append(errhost, myhost)
		return
	}
}

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
