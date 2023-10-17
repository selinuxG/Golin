package run

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"golin/global"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SqlServer(cmd *cobra.Command, args []string) {
	//确认结果是否输出
	echotype, err := cmd.Flags().GetBool("echo")
	if err != nil {
		fmt.Println(err)
		return
	}
	echorun = echotype

	//获取分隔符
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
		Onlyonerun(value, spr, "sqlserver")
		wg.Wait()
		return
	}
	// 下面开始执行批量的
	ippath, err := cmd.Flags().GetString("ip")
	if err != nil {
		fmt.Println(err)
		return
	}
	//判断pgsql.txt文件是否存在
	Checkfile(ippath, fmt.Sprintf("名称%sip%s用户%s密码%s端口", Split, Split, Split, Split), global.FilePer, ippath)
	// 运行share文件中的函数
	Rangefile(ippath, spr, "sqlserver")
	wg.Wait()
	//完成前最后写入文件
	Deffile("pgsql", count, count-len(errhost), errhost)

}

// User 用户安全策略相关
type User struct {
	UserName            string    //登录名（用户名）
	IsDisabled          bool      //表示用户是否被禁用，false 表示未禁用，true 表示已禁用
	IsPolicyChecked     bool      //表示该用户是否遵循Windows密码策略，true 表示遵循，false 表示不遵循
	IsExpirationChecked bool      //表示是否允许用户密码过期，true 表示允许过期，false 表示不允许过期
	DaysUntilExpiration *int64    //剩余天数直至密码过期（如果启用了密码过期策略）
	PasswordLastSetTime time.Time //表示上次密码更改的时间
}

// ConfigTimeout 超时退出相关
type ConfigTimeout struct {
	Name        string //remote query timeout
	Minimum     int    //此配置选项允许的最小值。对于 "remote query timeout"，最小值为0，表示禁用远程查询超时
	Maximum     int    //此配置选项允许的最大值。对于 "remote query timeout"，最大值为 2147483647（约68年）。
	ConfigValue int    //当前为该配置选项设置的值。这里显示的值可能会在下次服务器启动时生效。
	RunValue    int    //当前正在运行的值。服务器实际上使用的值可能与配置值不同，特别是当修改了配置并未重新启动服务时。通常情况下，在执行RECONFIGURE命令之后，config_value和run_value应该相同。
}

// SqlServerrun 连接SQLServer运行内置命令
func SqlServerrun(name, host, user, passwd, port string) {
	defer wg.Done()
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=master&timeout=1.5s&encrypt=disable", user, passwd, host, port)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		errhost = append(errhost, host)
		return
	}
	//确认采集完成目录是否存在
	fullPath := filepath.Join(succpath, "sqlserver")
	_, err = os.Stat(fullPath)
	if os.IsNotExist(err) {
		os.MkdirAll(fullPath, os.FileMode(global.FilePer))
	}
	firenmame := filepath.Join(fullPath, fmt.Sprintf("%s_%s.log", name, host))
	//先删除之前的同名记录文件
	os.Remove(firenmame)
	file, _ := os.OpenFile(firenmame, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.FileMode(global.FilePer))
	defer file.Close()
	write := bufio.NewWriter(file)
	//查询版本
	write.WriteString("\n-----版本信息：\n")
	rows, _ := db.Raw("select @@version").Rows()
	var version string
	for rows.Next() {
		rows.Scan(&version)
		write.WriteString(fmt.Sprintf("%s\n", strings.Replace(version, "\n", "", -1)))
	}

	//查询当前连接是否加密
	write.WriteString("\n-----查询当前连接是否加密：\n")
	rows, _ = db.Raw(`SELECT encrypt_option FROM sys.dm_exec_connections WHERE session_id = @@SPID;`).Rows()
	var nowssl string
	for rows.Next() {
		rows.Scan(&nowssl)
		write.WriteString(fmt.Sprintf("加密传输：%s\n", strings.Replace(nowssl, "\n", "", -1)))
	}

	//查询用户信息
	write.WriteString("\n-----用户信息查询：\n")

	var result []User
	sqlQuery := "SELECT  P.name AS UserName, P.is_disabled AS IsDisabled, L.is_policy_checked AS IsPolicyChecked, L.is_expiration_checked AS IsExpirationChecked, LOGINPROPERTY(P.name, 'DaysUntilExpiration') AS DaysUntilExpiration, LOGINPROPERTY(P.name, 'PasswordLastSetTime') AS PasswordLastSetTime FROM sys.server_principals P JOIN sys.sql_logins L ON P.principal_id = L.principal_id WHERE P.type_desc = 'SQL_LOGIN' ORDER BY P.name;"
	res := db.Raw(sqlQuery).Scan(&result)

	if res.Error != nil {
		fmt.Println("Error:", res.Error)
		return
	}
	for _, m := range result {
		username := m.UserName
		write.WriteString("用户名称：" + username + "\n")
		isdisable := m.IsDisabled
		write.WriteString(fmt.Sprintf("是否锁定：%t\n", isdisable))
		ispolicychecked := m.IsPolicyChecked
		write.WriteString(fmt.Sprintf("是否遵循windows密码规则：%t\n", ispolicychecked))
		IsExpirationChecked := m.IsExpirationChecked
		write.WriteString(fmt.Sprintf("是否允许密码过期：%t\n", IsExpirationChecked))
		//DaysUntilExpiration如果是空指针则为永不过期
		DaysUntilExpiration := m.DaysUntilExpiration
		daysUntilExpirationStr := "永不过期"
		if DaysUntilExpiration != nil {
			daysUntilExpirationStr = fmt.Sprintf("%d 天", *m.DaysUntilExpiration)
		}
		write.WriteString(fmt.Sprintf("密码过期时间：%v\n", daysUntilExpirationStr))
		PasswordLastSetTime := m.PasswordLastSetTime
		write.WriteString(fmt.Sprintf("上次修改密码时间：%v\n\n", PasswordLastSetTime))
	}

	write.WriteString("\n------空密码账户查询：\n")
	rows, _ = db.Raw(`SELECT name FROM sys.sql_logins WHERE PWDCOMPARE('', password_hash) = 1 ;`).Rows()
	var PWDCOMPARE string
	for rows.Next() {
		rows.Scan(&PWDCOMPARE)
		write.WriteString(fmt.Sprintf("空密码账户：%s\n", strings.Replace(PWDCOMPARE, "\n", "", -1)))
	}

	write.WriteString("\n------远程管理查询：\n")
	rows, _ = db.Raw(`SELECT value_in_use FROM sys.configurations WHERE name = 'remote access';`).Rows()
	var value_in_use string
	for rows.Next() {
		rows.Scan(&value_in_use)
		write.WriteString(fmt.Sprintf("远程登录管理状态：%s(0为关闭，1为开启)\n", strings.Replace(value_in_use, "\n", "", -1)))
	}

	write.WriteString("\n------超时相关：\n")

	var timeout ConfigTimeout
	db.Raw("EXEC sp_configure 'remote query timeout'").Scan(&timeout)
	write.WriteString(fmt.Sprintf("此项允许的最小值：%d（最小值如为0表示禁用远程查询超时）\n", timeout.Minimum))
	write.WriteString(fmt.Sprintf("当前运行的值：%d(秒)\n", timeout.RunValue))
	write.Flush()
	if echorun {
		readFile, _ := os.ReadFile(firenmame)
		fmt.Println(string(readFile))
	}
}
