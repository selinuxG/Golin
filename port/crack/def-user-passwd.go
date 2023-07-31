package crack

import (
	"embed"
	"golin/global"
	"strings"
)

var (
	passwdlist = []string{}
	userlist   = []string{}
)

//go:embed password.txt
var passwd embed.FS

func Passwdlist() []string {
	data, _ := passwd.ReadFile("password.txt")
	datastr := strings.ReplaceAll(string(data), "\r\n", "\n")
	for _, u := range strings.Split(datastr, "\n") {
		passwdlist = append(passwdlist, u)
	}
	passwdlist = global.RemoveDuplicates(passwdlist)
	return passwdlist
}

func Userlist(mode string) []string {
	switch mode {
	case "ssh":
		userlist = []string{"root", "admin"}
	case "mysql":
		userlist = []string{"root", "mysql"}
	case "redis":
		userlist = []string{"", "root"}
	case "pgsql":
		userlist = []string{"postgres", "root", "admin"}
	case "sqlserver":
		userlist = []string{"sa", "administrator"}
	case "ftp":
		userlist = []string{"ftp", "admin", "www", "web", "root", "db", "wwwroot", "data"}
	case "smb":
		userlist = []string{"administrator", "admin", "guest"}
	case "telnet":
		userlist = []string{"admin", "root"}
	case "tomcat":
		userlist = []string{"tomcat", "manager", "admin"}
	}
	return userlist
}
