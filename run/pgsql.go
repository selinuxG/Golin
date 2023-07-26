package run

import (
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	"golin/global"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"
	"strings"
)

func Pgsqlstart(cmd *cobra.Command, args []string) {
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
		Onlyonerun(value, spr, "pgsql")
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
	Rangefile(ippath, spr, "pgsql")
	wg.Wait()
	//完成前最后写入文件
	Deffile("pgsql", count, count-len(errhost), errhost)
}

// PgAuthID pg_authid表相关字段
type PgAuthID struct {
	Oid           int          `gorm:"column:oid"`
	Rolname       string       `gorm:"column:rolname"`
	Rolpassword   string       `gorm:"column:rolpassword"`
	Rolsuper      bool         `gorm:"column:rolsuper"`
	Rolcanlogin   bool         `gorm:"column:rolcanlogin"`
	Rolvaliduntil sql.NullTime `gorm:"column:rolvaliduntil"` // Use sql.NullTime for nullable time fields
	Rolcreaterole bool         `gorm:"column:rolcreaterole"`
	Rolcreatedb   bool         `gorm:"column:rolcreatedb"`
	Rolinherit    bool         `gorm:"column:rolinherit"`
}

// DatabasePrivileges 数据库权限
type DatabasePrivileges struct {
	DatabaseName string `gorm:"column:database_name"`
	RoleName     string `gorm:"column:role_name"`
	Privileges   string `gorm:"column:privileges"`
}

// RoleMember 角色关系
type RoleMember struct {
	RoleName string       `gorm:"column:role_name"`
	MemberOf sql.RawBytes `gorm:"column:memberof"`
}

// version 变量结构
type Version struct {
	PgVersion string `gorm:"column:version"`
}

// 连接pgsql运行内置命令
func Pgsql(name, host, user, passwd, port string) {
	defer wg.Done()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable  TimeZone=Asia/Shanghai connect_timeout=3", host, user, passwd, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		errhost = append(errhost, host)
		return
	}
	//确认采集完成目录是否存在
	fullPath := filepath.Join(succpath, "pgsql")
	_, err = os.Stat(fullPath)
	if os.IsNotExist(err) {
		os.MkdirAll(fullPath, os.FileMode(global.FilePer))
	}
	firenmame := filepath.Join(fullPath, fmt.Sprintf("%s_%s.html", name, host))
	//先删除之前的同名记录文件
	os.Remove(firenmame)

	html := pgsqlhtml()
	var version Version
	db.Raw("SELECT version()").Scan(&version)
	html = strings.ReplaceAll(html, "版本信息详细信息", fmt.Sprintf("<tr><td>%s</td></tr>", version.PgVersion))
	echoinfo := ""
	// 查询用户相关信息
	var result []PgAuthID
	db.Raw("SELECT oid,rolname,rolpassword,rolsuper,rolcanlogin,rolvaliduntil,rolcreaterole,rolcreatedb,rolinherit FROM pg_authid").Scan(&result)
	for _, user := range result {
		//当从数据库检索数据时，如果日期/时间字段非 null，则 Valid 字段将设置为 true，并且您可以通过访问 Time 字段来获取实际日期/时间值。如果日期/时间字段为 null，则 Valid 字段将设置为 false，并且 Time 字段的值未定义。
		rolvaliduntil := ""
		if user.Rolvaliduntil.Valid {
			rolvaliduntil = fmt.Sprintf("%v  	", user.Rolvaliduntil.Time)
		} else {
			rolvaliduntil = "永不过期"
		}

		echoinfo += fmt.Sprintf("<tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>",
			user.Rolname, user.Oid, user.Rolsuper, rolvaliduntil, user.Rolpassword, user.Rolcanlogin, user.Rolcreatedb, user.Rolcreaterole, user.Rolinherit)
	}

	html = strings.ReplaceAll(html, "user详细信息", echoinfo)

	echoinfo = ""
	var res []RoleMember
	db.Raw(`SELECT r.rolname AS role_name, ARRAY(SELECT b.rolname FROM pg_catalog.pg_auth_members m JOIN pg_catalog.pg_roles b ON (m.roleid = b.oid) WHERE m.member = r.oid) as memberof FROM pg_catalog.pg_roles r`).Scan(&res)
	for _, role := range res {
		str := strings.Trim(string(role.MemberOf), "{}") // 去除花括号
		echoinfo += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", role.RoleName, str)
	}
	html = strings.ReplaceAll(html, "role详细信息", echoinfo)

	echoinfo = ""
	rows, _ := db.Raw("SHOW shared_preload_libraries;").Rows()
	var libraries string
	for rows.Next() {
		rows.Scan(&libraries)
		if libraries != "" {
			echoinfo += fmt.Sprintf("<tr><td>%s</td></tr>", echoinfo)
		} else {
			echoinfo += fmt.Sprintf("<tr><td></td></tr>")
		}
	}
	html = strings.ReplaceAll(html, "插件详细信息", echoinfo)
	html = strings.ReplaceAll(html, "网卡地址详细信息", fmt.Sprintf("<tr><td>%s</td></tr>", sqlcon(db, `SHOW listen_addresses;`)))
	html = strings.ReplaceAll(html, "最低支持的TLS版本", fmt.Sprintf("<tr><td>%s</td></tr>", sqlcon(db, `SHOW ssl_min_protocol_version;`)))
	html = strings.ReplaceAll(html, "log信息", fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
		sqlcon(db, `SHOW logging_collector;`),
		sqlcon(db, `SHOW log_min_messages;`),
		sqlcon(db, `SHOW log_directory;`),
		sqlcon(db, `SHOW log_filename;`),
		sqlcon(db, `SHOW log_statement;`),
		sqlcon(db, `SHOW log_connections;`),
		sqlcon(db, `SHOW log_disconnections;`),
		sqlcon(db, `SHOW log_line_prefix;`),
		sqlcon(db, `SHOW log_destination;`),
	))
	html = strings.ReplaceAll(html, "替换名称", host)

	os.WriteFile(firenmame, []byte(html), os.FileMode(global.FilePer))
	if echorun {
		fmt.Println("当前模式已生成HTML报告,不在支持屏幕输出...")
	}

}

// sqlcon 执行sql命令
func sqlcon(db *gorm.DB, sql string) string {
	rows, err := db.Raw(sql).Rows()
	if err != nil {
		return ""
	}
	var sqlecho string
	for rows.Next() {
		rows.Scan(&sqlecho)
	}
	return sqlecho

}
