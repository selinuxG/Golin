package crack

import (
	"embed"
	"golin/global"
	"os"
	"strings"
)

var (
	passwdlist []string
	userlist   []string
)

var userMap = map[string][]string{
	"ssh":        {"root"},
	"mysql":      {"root"},
	"redis":      {"", "root"},
	"postgresql": {"postgres", "root"},
	"sqlserver":  {"sa", "administrator"},
	"ftp":        {"ftp", "admin"},
	"smb":        {"administrator", "guest"},
	"telnet":     {"admin", "root"},
	"tomcat":     {"tomcat", "manager", "admin"},
	"rdp":        {"administrator"},
	"oracle":     {"orcl", "sys", "system"},
}

//go:embed password.txt
var passwd embed.FS

// Checkdistfile 用户自定义的字典文件
func Checkdistfile(userfile, passwdfile string) {
	if global.PathExists(userfile) {
		data, _ := os.ReadFile(userfile)
		datastr := strings.ReplaceAll(string(data), "\r\n", "\n")
		userlist = append(userlist, strings.Split(datastr, "\n")...)
	}
	if global.PathExists(passwdfile) {
		data, _ := os.ReadFile(passwdfile)
		datastr := strings.ReplaceAll(string(data), "\r\n", "\n")
		passwdlist = append(passwdlist, strings.Split(datastr, "\n")...)
	}

}

func Passwdlist() []string {

	if len(passwdlist) > 0 {
		return global.RemoveDuplicates(passwdlist)
	}

	data, _ := passwd.ReadFile("password.txt")
	datastr := strings.ReplaceAll(string(data), "\r\n", "\n")
	for _, u := range strings.Split(datastr, "\n") {
		passwdlist = append(passwdlist, u)
	}
	passwdlist = global.RemoveDuplicates(passwdlist)
	return passwdlist
}

func Userlist(mode string) []string {

	if len(userlist) > 0 {
		return global.RemoveDuplicates(userlist)
	}

	return userMap[mode]
}
