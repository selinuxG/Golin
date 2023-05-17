package global

import (
	"net"
	"os"
)

const (
	XlsxTemplateName = "golin上传文件模板文件.xlsx"
)

// 文件相关公共变量
var (
	SuccessLog       = "log.log"       //运行记录
	Split            = "~"             //默认分割符号
	CmdLinuxPath     = "linux.txt"     //默认Linux模式多主机模式下读取的文件
	CmdMysqlPath     = "mysql.txt"     //默认Mysql模式多主机模式下读取的文件
	CmdRedisPath     = "redis.txt"     //默认Redis模式多主机模式下读取的文件
	CmdRoutepath     = "route.txt"     //默认route模式多主机模式下读取的文件
	CmdPgsqlPath     = "pgsql.txt"     //默认Linux模式多主机模式下读取的文件
	CmdsqlServerPath = "sqlserver.txt" //默认Linux模式多主机模式下读取的文件

	FilePer      = 0644                                                 //创建文件或目录时的默认权限，必须是0开头
	Succpath     = "采集完成目录"                                             //CLi模式成功主机写入日志的目录
	Succwebpath  = "webhistory.json"                                    //Web模式运行记录
	Denynametype = []string{"\\", "\\/", "*", "?", "\"", "<", ">", "|"} //windos下不允许创建名称的特殊符号。
)

// 各类python程序路径
var (
	PythonPath = "python" //默认运行python的环境变量
	PythonDir  = "python" //承载各类python脚本的目录位置
	PyHw       = "hw.py"  //python ssh hw
	PyGui      = "Gui.py" //python gui 程序
)

// AppendToFile 创建追加写入函数
func AppendToFile(filename string, content string) error {
	// 检测文件是否存在
	var file *os.File
	var err error
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		// 如果文件不存在，则创建文件
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
	} else {
		// 如果文件存在，则以追加模式打开文件
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.FileMode(FilePer))
		if err != nil {
			return err
		}
	}
	// 在文件末尾写入内容
	defer file.Close()
	if _, err = file.WriteString(content); err != nil {
		return err
	}
	return nil
}

// PathExists 文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	//fmt.Println(err)
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// GetLocalIPAddresses 获取本机网卡ip
func GetLocalIPAddresses() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}

		ip := ipNet.IP.To4()
		if ip == nil { // 排除IPv6地址
			continue
		}

		ips = append(ips, ip.String())
	}

	return ips, nil
}
