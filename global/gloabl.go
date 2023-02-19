package global

var (
	SuccessLog   = "Golin运行日志.log"                                      //运行记录
	Split        = "~"                                                  //默认分割符号
	CmdLinuxPath = "linux.txt"                                          //默认Linux模式多主机模式下读取的文件
	CmdMysqlPath = "mysql.txt"                                          //默认Mysql模式多主机模式下读取的文件
	CmdRedisPath = "redis.txt"                                          //默认Redis模式多主机模式下读取的文件
	FilePer      = 0744                                                 //创建文件或目录时的默认权限，必须是0开头
	Succpath     = "采集完成目录"                                             //成功主机写入日志的目录
	Denynametype = []string{"\\", "\\/", "*", "?", "\"", "<", ">", "|"} //windos下不允许创建名称的特殊符号。
)
