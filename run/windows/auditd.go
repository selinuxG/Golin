//go:build windows

package windows

import (
	"fmt"
	"regexp"
	"strings"
)

func auditd() {
	var auditlist []Policyone
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

	//高级审核策略
	aud := ExecCommands(`auditpol /get /category:*`)
	aud = strings.ReplaceAll(aud, "\r\r\n", "\n")
	for _, v := range strings.Split(aud, "\n") {
		// 创建一个正则表达式匹配多个空格 将多个连续空格替换为一个
		re := regexp.MustCompile(`\s{2,}`)
		v = re.ReplaceAllString(v, " ")
		if len(strings.Split(v, " ")) == 3 {
			one := Policyone{Name: strings.Split(v, " ")[1], Value: strings.Split(v, " ")[2], Steer: "配置为成功和失败"}
			switch strings.Split(v, " ")[2] {
			case "成功":
				one.Static = No
			case "无审核":
				one.Static = No
			case "成功和失败":
				one.Static = Yes
			}
			auditlist = append(auditlist, one)
		}
	}
	echo = ""
	for _, v := range auditlist {
		echo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", v.Name, v.Value, v.Static, v.Steer)
	}
	html = strings.ReplaceAll(html, "高级审计策略结果", echo)

	echo = ""
	eventLogs := []string{"Application", "Security", "Setup", "System"}
	for _, v := range eventLogs {
		cmd := ExecCommands(fmt.Sprintf("wevtutil get-log %s", v))
		echo += cmd + "\n"
	}
	html = strings.ReplaceAll(html, "日志属性结果", echo)
}
