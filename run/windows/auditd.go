//go:build windows

package windows

import (
	"fmt"
	"strings"
)

func auditd() {
	auditlist := []Policyone{}
	one := Policyone{}
	for key, value := range auditmap {
		v, ok := Policy[key] //是否审核系统事件
		if ok {
			switch v {
			case "0":
				one = Policyone{Name: value, Static: No, Value: "未开启", Steer: "开启成功与失败事件"}
			case "1":
				one = Policyone{Name: value, Static: No, Value: "仅成功事件", Steer: "开启成功与失败事件"}
			case "2":
				one = Policyone{Name: value, Static: No, Value: "仅失败事件", Steer: "开启成功与失败事件"}
			case "3":
				one = Policyone{Name: value, Static: Yes, Value: "成功和失败事件", Steer: "开启成功与失败事件"}
			}
		}
		auditlist = append(auditlist, one)
	}
	echo := ""
	for _, v := range auditlist {
		echo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", v.Name, v.Value, v.Static, v.Steer)
	}
	html = strings.ReplaceAll(html, "审计相关结果", echo)
}
