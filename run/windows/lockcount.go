package windows

import (
	"fmt"
	"strconv"
	"strings"
)

// lock 核查锁定次数
func lock() []Policyone {
	locls := []Policyone{}
	v, ok := Policy["LockoutBadCount"] //在锁定账户之前允许的无效登录尝试次数
	intv, _ := strconv.Atoi(v)
	if ok {
		one := Policyone{Name: "检查在锁定账户之前允许的无效登录尝试次数", Static: No, Value: v, Steer: "不多于5次"}
		if intv > 0 && intv < 5 {
			one.Static = Yes
		}
		locls = append(locls, one)
	}

	v, ok = Policy["LockoutDuration"] // 账户被锁定的时间（分钟）
	intv, _ = strconv.Atoi(v)
	if ok {
		one := Policyone{Name: "检查账户被锁定的时间(分钟)", Static: No, Value: v, Steer: "锁定不小于1分钟"}
		if intv > 1 {
			one.Static = Yes
		}
		locls = append(locls, one)
	}
	v, ok = Policy["ResetLockoutCount"] // 在重置“无效登录尝试次数”计数器之前的时间（分钟）。
	intv, _ = strconv.Atoi(v)
	if ok {
		one := Policyone{Name: "检查在重置“无效登录尝试次数”计数器之前的时间（分钟）。", Static: No, Value: v, Steer: "不小于3分钟"}
		if intv > 1 {
			one.Static = Yes
		}
		locls = append(locls, one)
	}
	echo := ""
	for _, v := range locls {
		echo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", v.Name, v.Value, v.Static, v.Steer)
	}
	html = strings.ReplaceAll(html, "失败锁定结果", echo)

	return locls
}
