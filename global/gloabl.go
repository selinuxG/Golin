package global

import (
	"os"
	"sync"
)

const (
	XlsxTemplateName = "golin上传文件模板文件.xlsx"
)

// 公共变量
var (
	SuccessLog       = "log.log"                                            //运行记录
	Split            = "~"                                                  //默认分割符号
	CmdLinuxPath     = "linux.txt"                                          //默认Linux模式多主机模式下读取的文件
	CmdMysqlPath     = "mysql.txt"                                          //默认Mysql模式多主机模式下读取的文件
	CmdRedisPath     = "redis.txt"                                          //默认Redis模式多主机模式下读取的文件
	CmdRoutepath     = "route.txt"                                          //默认route模式多主机模式下读取的文件
	CmdPgsqlPath     = "pgsql.txt"                                          //默认pgsql模式多主机模式下读取的文件
	CmdsqlServerPath = "sqlserver.txt"                                      //默认sqlserver模式多主机模式下读取的文件
	CmdOraclePath    = "oracle.txt"                                         //默认oracle模式多主机模式下读取的文件
	FilePer          = 0744                                                 //创建文件或目录时的默认权限，必须是0开头
	Succpath         = "采集完成目录"                                             //CLi模式成功主机写入日志的目录
	Succwebpath      = "webhistory.json"                                    //Web模式运行记录
	Denynametype     = []string{"\\", "\\/", "*", "?", "\"", "<", ">", "|"} //windos下不允许创建名称的特殊符号。
	PrintLock        sync.RWMutex                                           //并发输出写入
	WebURl           = ""                                                   //web扫描时临时后缀
	SaveIMG          = false                                                //web扫描时是否进行截图,本地需要有chrom浏览器
	SsaveIMGDIR      = "WebScreenshot"
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

// InSlice 判断字符串是否在 slice 中。
func InSlice(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates 切片去重
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var list []string

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
