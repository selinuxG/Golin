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

// 连接pgsql运行内置命令
func SqlServerrun(name, host, user, passwd, port string) {
	defer wg.Done()
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=master&timeout=1.5s", user, passwd, host, port)
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
	rows, _ := db.Raw("select @@version").Rows()
	var version string
	for rows.Next() {
		rows.Scan(&version)
		write.WriteString(fmt.Sprintf("当前版本为：%s\n", strings.Replace(version, "\n", "", -1)))
	}

	write.Flush()
	if echorun {
		readFile, _ := os.ReadFile(firenmame)
		fmt.Println(string(readFile))
	}
}
