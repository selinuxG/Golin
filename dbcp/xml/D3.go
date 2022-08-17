package xml

import "strings"

func D3(filestring string) (str, steer, issue, analyze, success string) {
	str = "经核查，"
	//整改建议。
	steer = ""
	//实际问题
	issue = ""
	//危害分析
	analyze = ""
	//符合情况
	success = ""
	if sysbak := strings.Contains(string(filestring), "未做限制"); sysbak {
		str += `通过查看/etc/hosts.allow以及/etc/hosts.deny中未对网络地址进行限制。`
		issue += "未严格限制终端登录地址范围。"
		analyze += "恶意用户可从网内任意地址尝试对设备进行访问、攻击，存在非授权访问的风险。"
		steer += "建议对管理终端的网络地址范围进行限制。"
		success += "不符合"
		return str, steer, issue, analyze, success
	} else {
		str += `通过查看/etc/hosts.allow以及/etc/hosts.deny中以对网络地址进行限制。`
		success += "符合"
		return str, steer, issue, analyze, success
	}
}
