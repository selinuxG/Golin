package run

import (
	"bufio"
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
	firenmame := filepath.Join(fullPath, fmt.Sprintf("%s_%s.log", name, host))
	//先删除之前的同名记录文件
	os.Remove(firenmame)
	file, _ := os.OpenFile(firenmame, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.FileMode(global.FilePer))
	defer file.Close()
	write := bufio.NewWriter(file)

	var version Version
	db.Raw("SELECT version()").Scan(&version)
	write.WriteString(fmt.Sprintf("-----------------版本信息：\n%s\n", version.PgVersion))
	write.WriteString(fmt.Sprintf("-----------------用户相关：\n"))

	// 查询用户相关信息
	var result []PgAuthID
	db.Raw("SELECT oid,rolname,rolpassword,rolsuper,rolcanlogin,rolvaliduntil,rolcreaterole,rolcreatedb,rolinherit FROM pg_authid").Scan(&result)
	for _, user := range result {
		write.WriteString(fmt.Sprintf("账户名：%s	", user.Rolname))
		write.WriteString(fmt.Sprintf("账户ID：%d	", user.Oid))
		write.WriteString(fmt.Sprintf("超级用户：%t	", user.Rolsuper))
		//当从数据库检索数据时，如果日期/时间字段非 null，则 Valid 字段将设置为 true，并且您可以通过访问 Time 字段来获取实际日期/时间值。如果日期/时间字段为 null，则 Valid 字段将设置为 false，并且 Time 字段的值未定义。
		if user.Rolvaliduntil.Valid {
			write.WriteString(fmt.Sprintf("口令过期时间：%v  	", user.Rolvaliduntil.Time))
		} else {
			write.WriteString("口令过期时间为：永不过期		")
		}
		pass := []rune(user.Rolpassword) // 如果密码啊长度大于 30
		if len(pass) > 30 {
			write.WriteString(fmt.Sprintf("密码：%s	", string(pass[:30])))
		} else {
			write.WriteString(fmt.Sprintf("密码：%s	", user.Rolpassword))
		}
		write.WriteString(fmt.Sprintf("\n权限相关：是否可登录：%t 是否可创建数据库：%t 是否可以创建角色：%t 是否可以继承其所属角色的权限%t \n", user.Rolcanlogin, user.Rolcreatedb, user.Rolcreaterole, user.Rolinherit))

		//根据用户查询数据库权限
		if !user.Rolsuper && user.Rolcanlogin {
			var datapro []DatabasePrivileges
			db.Raw(strings.ReplaceAll(datapriv, "gys", user.Rolname)).Scan(&datapro)
			for _, privileges := range datapro {
				write.WriteString(fmt.Sprintf("数据库：%s  详细权限：%s\n", privileges.DatabaseName, privileges.Privileges))
			}
			write.WriteString("\n")
		}
		write.WriteString("\n")
	}

	write.WriteString(fmt.Sprintf("\n-----------------角色关系：\n"))
	var res []RoleMember
	db.Raw(`SELECT r.rolname AS role_name, ARRAY(SELECT b.rolname FROM pg_catalog.pg_auth_members m JOIN pg_catalog.pg_roles b ON (m.roleid = b.oid) WHERE m.member = r.oid) as memberof FROM pg_catalog.pg_roles r`).Scan(&res)
	for _, role := range res {
		str := strings.Trim(string(role.MemberOf), "{}") // 去除花括号
		//strArr := strings.Split(str, ",")                // 按逗号分割字符串
		write.WriteString(fmt.Sprintf("角色：%s 隶属于：%v\n", role.RoleName, str))
	}

	write.WriteString(fmt.Sprintf("\n-----------------安全策略：\n"))
	write.WriteString(fmt.Sprintf("已安装插件：\n"))
	rows, _ := db.Raw("SHOW shared_preload_libraries;").Rows()
	var libraries string
	for rows.Next() {
		rows.Scan(&libraries)
		if libraries != "" {
			write.WriteString(libraries + "\n")
		}
	}
	write.WriteString(fmt.Sprintf("远程地址：%s\n", sqlcon(db, `SHOW listen_addresses;`)))
	write.WriteString(fmt.Sprintf("最低支持的TLS版本：%s\n", sqlcon(db, `SHOW ssl_min_protocol_version;`)))
	write.WriteString(fmt.Sprintf("\n-----------------日志相关：\n"))
	write.WriteString(fmt.Sprintf("错误日志(logging_collector)状态：%s\n", sqlcon(db, `SHOW logging_collector;`)))
	write.WriteString(fmt.Sprintf("错误日志记录级别：%s\n", sqlcon(db, `SHOW log_min_messages;`)))
	write.WriteString(fmt.Sprintf("日志文件存储目录：%s\n", sqlcon(db, `SHOW log_directory;`)))
	write.WriteString(fmt.Sprintf("文件命令格式：%s\n", sqlcon(db, `SHOW log_filename;`)))
	write.WriteString(fmt.Sprintf("查询语句开启状态：%s\n", sqlcon(db, `SHOW log_statement;`)))
	write.WriteString(fmt.Sprintf("用户登录记录开启状态：%s\n", sqlcon(db, `SHOW log_connections;`)))
	write.WriteString(fmt.Sprintf("用户登出记录开启状态：%s\n", sqlcon(db, `SHOW log_disconnections;`)))
	write.WriteString(fmt.Sprintf("日志内容记录字段：%s\n", sqlcon(db, `SHOW log_line_prefix;`)))
	write.WriteString(fmt.Sprintf("日志发送类型：%s\n", sqlcon(db, `SHOW log_destination;`)))

	write.Flush()
	if echorun {
		readFile, _ := os.ReadFile(firenmame)
		fmt.Println(string(readFile))
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

var datapriv = `
SELECT 
    d.datname AS database_name,
    r.rolname AS role_name,
    json_build_object(
        '是否允许连接数据库', has_database_privilege(r.rolname, d.datname, 'CONNECT'),
        '是否允许创建新表', has_database_privilege(r.rolname, d.datname, 'CREATE'),
        '是否允许创建临时表', has_database_privilege(r.rolname, d.datname, 'TEMPORARY'),
        '表权限', (
            SELECT 
                json_agg(
                    json_build_object(
                        'table', tp.table_name,
                        'privilege', tp.privilege_type
                    )
                )
            FROM 
                information_schema.table_privileges tp
            WHERE 
                tp.table_catalog = d.datname AND
                tp.grantee = r.rolname
        )
    ) as privileges
FROM 
    pg_roles r
CROSS JOIN 
    pg_database d
WHERE 
    r.rolname = 'gys';
`
