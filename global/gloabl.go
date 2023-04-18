package global

import (
	"os"
)

var (
	SuccessLog   = "log.log"                                            //运行记录
	Split        = "~"                                                  //默认分割符号
	CmdLinuxPath = "linux.txt"                                          //默认Linux模式多主机模式下读取的文件
	CmdMysqlPath = "mysql.txt"                                          //默认Mysql模式多主机模式下读取的文件
	CmdRedisPath = "redis.txt"                                          //默认Redis模式多主机模式下读取的文件
	CmdRoutepath = "route.txt"                                          //默认route模式多主机模式下读取的文件
	FilePer      = 0644                                                 //创建文件或目录时的默认权限，必须是0开头
	Succpath     = "采集完成目录"                                             //成功主机写入日志的目录
	Denynametype = []string{"\\", "\\/", "*", "?", "\"", "<", ">", "|"} //windos下不允许创建名称的特殊符号。
)

// 各类python程序路径
var (
	Python_path = "python"       //默认运行python路径
	Py_hw       = "python/hw.py" //python ssh hw
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
	if os.IsNotExist(err) {
		return false
	}
	return false
}
