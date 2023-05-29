package windows

import (
	"fmt"
	"golin/global"
	"strings"
)

type UserInformation struct {
	Name       string //用户名
	Nameall    string //全名
	Sid        string //SID
	Annotation string //注释
	Lock       string //是否启用
	AfTime     string //帐户到期
	Afpasswd   string //上次修改密码时间
	PasswdAf   string //密码到期
	Checkpass  string //需要密码
	Aflogin    string //上次登录
	Group      string //本地组
}

func usercheck() {
	echo := ""
	userlist := []string{}
	use := global.ExecCommands("wmic useraccount list full")
	use = strings.ReplaceAll(use, "\r\r\n", "\n")
	for _, i := range strings.Split(use, "\n") {
		if strings.Count(i, "Name=") == 1 && strings.Count(i, "FullName") == 0 {
			userlist = append(userlist, strings.Split(i, "=")[1])
		}
	}
	//var users []UserInformation //结构体用户信息
	for _, user := range userlist {
		uu := UserInformation{}
		u := global.ExecCommands("chcp 936", fmt.Sprintf("net user %s", user))
		u = strings.ReplaceAll(u, "\r\n", "\n")
		for _, s := range strings.Split(u, "\n") {
			if len(s) > 0 {
				words := strings.Fields(s)
				s = strings.Join(words, " ")
				uservaule := strings.Split(s, " ")
				if len(uservaule) >= 2 {
					switch uservaule[0] {
					case "用户名":
						uu.Name = uservaule[1]
					case "全名":
						uu.Nameall = uservaule[1]
					case "注释", "Comment":
						echo := ""
						if len(uservaule) >= 1 {
							for i := 1; i < len(uservaule); i++ {
								echo += uservaule[i]
							}
						}
						uu.Annotation = echo
					case "帐户启用":
						uu.Lock = uservaule[1]
					case "帐户到期":
						uu.AfTime = uservaule[1]
					case "上次设置密码":
						uu.Afpasswd = uservaule[1]
					case "密码到期":
						uu.PasswdAf = uservaule[1]
					case "需要密码":
						uu.Checkpass = uservaule[1]
					case "上次登录":
						uu.Aflogin = uservaule[1]
					case "本地组成员": //组可能有多个
						group := ""
						if len(uservaule) > 1 {
							for i := 1; i < len(uservaule); i++ {
								group += (uservaule[i])
							}
						}
						uu.Group = group
					}
				}
			}
		}
		//获取用户sid
		sid := global.ExecCommands(fmt.Sprintf("wmic useraccount where name='%s' get sid", user))
		sid = strings.ReplaceAll(sid, "\r\r\n", "")
		sid = strings.ReplaceAll(sid, "SID", "")
		sid = strings.ReplaceAll(sid, " ", "")
		uu.Sid = sid
		echo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", uu.Name, uu.Nameall, uu.Sid, uu.Annotation, uu.Lock, uu.AfTime, uu.Afpasswd, uu.Checkpass, uu.PasswdAf, uu.Aflogin, uu.Group)
	}
	html = strings.ReplaceAll(html, "用户详细信息", echo)

}
