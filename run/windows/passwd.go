package windows

import (
	"fmt"
	"strconv"
	"strings"
)

// checkpasswd 检查密码复杂度相关策略
func checkpasswd() {
	var passwd []Policyone
	v, ok := Policy["PasswordComplexity"] //是否开启启用密码复杂度要求
	if ok {
		one := Policyone{Name: "检查是否开启密码复杂度要求", Static: No, Value: "未开启", Steer: "开启密码复杂度要求"}
		if v == "1" {
			one.Static = Yes
			one.Value = "已开启"
		}
		passwd = append(passwd, one)
	}

	v, ok = Policy["MinimumPasswordLength"] //检查密码的最短长度要求
	intv, _ := strconv.Atoi(v)
	if ok {
		one := Policyone{Name: "检查密码的最短长度要求（位数）", Static: No, Value: v, Steer: "长度不低于8位"}
		if intv >= 8 {
			one.Static = Yes
		}
		passwd = append(passwd, one)
	}

	v, ok = Policy["PasswordHistorySize"] //在允许重新使用相同密码之前必须使用的唯一新密码的数量。
	intv, _ = strconv.Atoi(v)
	if ok {
		one := Policyone{Name: "检查在允许重新使用相同密码之前必须使用的唯一新密码的数量", Static: No, Value: v, Steer: "不允许"}
		if intv == 0 {
			one.Static = Yes
		}
		passwd = append(passwd, one)
	}
	//替换
	echo := ""
	for _, v := range passwd {
		echo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", v.Name, v.Value, v.Static, v.Steer)
	}
	html = strings.ReplaceAll(html, "密码复杂度结果", echo)
	//密码有效期
	passwd = []Policyone{}
	v, ok = Policy["MaximumPasswordAge"] //密码的最长有效期限（天数）
	intv, _ = strconv.Atoi(v)
	if ok {
		one := Policyone{Name: "检查密码的最长有效期限（天数）", Static: No, Value: v, Steer: "不高于90天定期更换口令"}
		if intv <= 90 {
			one.Static = Yes
		}
		passwd = append(passwd, one)
	}
	v, ok = Policy["MinimumPasswordAge"] //密码的最短有效期限（天数）
	intv, _ = strconv.Atoi(v)
	if ok {
		one := Policyone{Name: "检查密码的最短有效期限（天数）", Static: No, Value: v, Steer: "不少于1天定期更换口令"}
		if intv >= 1 {
			one.Static = Yes
		}
		passwd = append(passwd, one)
	}
	echo = ""
	for _, v := range passwd {
		echo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", v.Name, v.Value, v.Static, v.Steer)
	}
	html = strings.ReplaceAll(html, "密码有效期检查结果", echo)

}
