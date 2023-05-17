package run

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	"golin/global"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"path/filepath"
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
}

// version 变量结构
type Version struct {
	PgVersion string `gorm:"column:version"`
}

// 连接pgsql运行内置命令
func Pgsql(name, host, user, passwd, port string) {
	defer wg.Done()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable  TimeZone=Asia/Shanghai connect_timeout=3", host, user, passwd, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		errhost = append(errhost, host)
		//fmt.Println("Error connecting to the database:", err)
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
	write.WriteString(fmt.Sprintf("版本信息：%s\n", version.PgVersion))
	write.WriteString(fmt.Sprintf("-----------------用户相关：\n"))

	// 查询用户相关信息
	var result []PgAuthID
	db.Raw("SELECT oid,rolname,rolpassword,rolsuper,rolcanlogin,rolvaliduntil FROM pg_authid").Scan(&result)
	for _, user := range result {
		write.WriteString(fmt.Sprintf("账户名：%s\n", user.Rolname))
		write.WriteString(fmt.Sprintf("账户ID：%d\n", user.Oid))
		write.WriteString(fmt.Sprintf("密码：%s\n", user.Rolpassword))
		write.WriteString(fmt.Sprintf("是否为超级用户：%t\n", user.Rolsuper))
		write.WriteString(fmt.Sprintf("是否可登录：%t\n", user.Rolcanlogin))
		//当从数据库检索数据时，如果日期/时间字段非 null，则 Valid 字段将设置为 true，并且您可以通过访问 Time 字段来获取实际日期/时间值。如果日期/时间字段为 null，则 Valid 字段将设置为 false，并且 Time 字段的值未定义。
		if user.Rolvaliduntil.Valid {
			write.WriteString(fmt.Sprintf("口令过期时间：%v\n", user.Rolvaliduntil.Time))
		} else {
			write.WriteString("口令过期时间为：空！")
		}
		write.WriteString("\n\n")
	}
	write.Flush()
	if echorun {
		readFile, _ := os.ReadFile(firenmame)
		fmt.Println(string(readFile))
	}
}
