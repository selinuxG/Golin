package run

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"golin/global"
	"os"
)

func Mysql(cmd *cobra.Command, args []string) {
	spr, err := cmd.Flags().GetString("spript")
	if err != nil {
		fmt.Println(err)
		return
	}
	//如果value值不为空则是运行一次的模式
	value, err := cmd.Flags().GetString("value")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(value) > 10 {
		Onlyonerun(value, spr, "Mysql")
		wg.Wait()
		zlog.Info("单次运行Mysql模式完成！")
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
	dsn := myuser + ":" + mypasswd + "@tcp(" + myhost + ":" + myport + ")" + "/mysql?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		errhost = append(errhost, myhost)
		return
	}
	err = db.Ping()
	if err != nil {
		errhost = append(errhost, myhost)
		return
	}
	_, err = os.Stat(succpath)
	if os.IsNotExist(err) {
		os.Mkdir(succpath, os.FileMode(global.FilePer))
	}
	fire := "采集完成目录//" + myname + "_" + myhost + "(mysql).log"
	//先删除之前的同名记录文件
	os.Remove(fire)
	file, err := os.OpenFile(fire, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.FileMode(global.FilePer))
	if err != nil {
		errhost = append(errhost, myhost)
		return
	}
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
	var user, host, File_priv, Shutdown_priv, grant_priv, ssl_type, password_expired, account_locked, authentication_string string
	users, _ := db.Query("SELECT  user,host,File_priv,Shutdown_priv,grant_priv,ssl_type,password_expired,account_locked,authentication_string from mysql.user")
	for users.Next() {
		users.Scan(&user, &host, &File_priv, &Shutdown_priv, &grant_priv, &ssl_type, &password_expired, &account_locked, &authentication_string)
		wriversion := "用户:" + user + "    远程登录权限:" + host + "\n" +
			"密码过期时间设置:" + password_expired + "  密码信息:" + authentication_string + "\n" +
			"是否锁定:" + account_locked +
			"	加密登录类型:" + ssl_type +
			"\n其他权限:   File_priv:  " + File_priv + "   Shutdown_priv  " + Shutdown_priv + "  grant_priv  " + grant_priv + "\n"
		write.WriteString(wriversion)

		//yuanc
		var usergrant string
		a := fmt.Sprintf(`SHOW GRANTS FOR "%s"@"%s"`, user, host)
		grantuser, _ := db.Query(a)
		write.WriteString("数据库相关权限: \n")
		for grantuser.Next() {
			grantuser.Scan(&usergrant)
			write.WriteString(usergrant + "\n")
		}
		write.WriteString("\n")
	}

	//超时时间
	write.WriteString("\n\n------查看超时时间------\n")
	var timeout_Variable_name, timeout_Value string
	time_conn, _ := db.Query("show variables like '%timeout%';")
	for time_conn.Next() {
		time_conn.Scan(&timeout_Variable_name, &timeout_Value)
		wriversion := timeout_Variable_name + "    字段值:" + timeout_Value + "\n"
		write.WriteString(wriversion)
	}

	//查看登录失败次数
	write.WriteString("\n\n------查看登录失败次数------\n")
	var connection_Variable_name, connection_Value string
	conne_conn, _ := db.Query(" show variables like '%connection_control%';")
	for conne_conn.Next() {
		conne_conn.Scan(&connection_Variable_name, &connection_Value)
		wriversion := connection_Variable_name + "    字段值:" + connection_Value + "\n"
		write.WriteString(wriversion)
	}

	//密码策略
	write.WriteString("\n\n------查看密码策略------\n")
	var password_Variable_name, password_Value string
	password_conn, _ := db.Query("show variables like 'validate_password%'")
	for password_conn.Next() {
		password_conn.Scan(&password_Variable_name, &password_Value)
		wriversion := password_Variable_name + "    字段值:" + password_Value + "\n"
		write.WriteString(wriversion)
	}

	//查看开启ssl
	write.WriteString("\n\n------查看是否开启SSL------\n")
	var ssl_Variable_name, ssl_Value string
	ssl_conn, _ := db.Query(" show global variables like '%ssl%'")
	for ssl_conn.Next() {
		ssl_conn.Scan(&ssl_Variable_name, &ssl_Value)
		wriversion := ssl_Variable_name + "    字段值:" + ssl_Value + "\n"
		write.WriteString(wriversion)
	}

	//查看开启审计功能
	write.WriteString("\n\n------查看是开启否审计------\n")
	var audit_Variable_name, audit_Value string
	auditd_conn, _ := db.Query("show global variables like '%general%'")
	for auditd_conn.Next() {
		auditd_conn.Scan(&audit_Variable_name, &audit_Value)
		wriversion := audit_Variable_name + "    字段值:" + audit_Value + "\n"
		write.WriteString(wriversion)
	}

	write.WriteString("\n\n\n-----以下是全局所有的配置参数\n")
	var Variable_name, Value string
	allValue, _ := db.Query("show global variables")
	for allValue.Next() {
		allValue.Scan(&Variable_name, &Value)
		wriversion := Variable_name + "    字段值:" + Value + "\n"
		write.WriteString(wriversion)
	}

	write.Flush()

}
