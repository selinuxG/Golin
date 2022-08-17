package xml

import "strings"

func A3(filestring string) (str, steer, issue, analyze, success string) {
	str = "经核查，"
	//错误数量，0为符合。4为不符合。其他为部分符合。
	errPASScoung := 0
	//整改建议。
	steer = ""
	//实际问题
	issue = ""
	//危害分析
	analyze = ""
	//符合情况
	success = ""
	if findtelnet := strings.Contains(string(filestring), "xinetd"); findtelnet {
		errPASScoung += 1
		str += "使用命令ps –ef |grep xinetd发现存在telnet明文传输协议，进行远程管理时，鉴别信息以明文方式进行传输。"
		steer += "建议关闭telnet明文传输协议，仅保留ssh等加密协议进行远程管理。"
		issue += "未关闭telnet明文传输协议。"
		analyze = "账号、口令等通过明文进行传输，可能导致敏感信息被恶意人员嗅探并盗用，存在非授权访问的风险。"
		success = "不符合"
	} else {
		str += `使用ps -ef |grep sshd命令发现存在sshd进程，已开启sshd服务；使用ps -ef |grep xinetd命令未发现telnet进程，未开启telnet服务。可防止鉴别信息在网络传输过程中被窃听。`
		success += "符合"
	}

	return str, steer, issue, analyze, success
}
