package crack

var (
	df_sshuser       = []string{"root", "ROOT", "test", "system", "super", "ceshi"} //ssh模式下默认用户
	df_mysqluser     = []string{"root", "mysql"}                                    //mysql模式下默认用户
	df_redisuser     = []string{"", "redis", "root"}                                //redis模式下默认用户
	df_pgsqluser     = []string{"postgres", "root"}                                 //pgsql模式下默认用户
	df_sqlserveruser = []string{"sa", "administrator"}                              //sqlserver模式下默认用户
	df_ftpuser       = []string{"ftp", "admin"}                                     //ftp模式下默认用户
	df_smbuser       = []string{"admin", "root"}                                    //smb模式下默认用户
	df_telnetuser    = []string{"admin", "root"}                                    //telnet模式下默认用户
	df_tomcatuser    = []string{"tomcat", "manager", "admin"}                       //tomcat模式下默认用户
	passwdlist       = []string{}
)
