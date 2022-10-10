package xml

import (
	"strings"
)

func B7(filestring string) (str, steer, issue, analyze, success string) {
	str = "经核查，"
	//整改建议。
	steer = ""
	//实际问题
	issue = ""
	//危害分析
	analyze = ""
	//符合情况
	success = ""
	if selinuxstatus := strings.Contains(string(filestring), "Selinux状态:enforcing"); selinuxstatus {
		str += `已开启安全增强型Linux子系统，SELinux设置为enforcing，对系统内重要主体和客体设置安全标记，并依据安全标记控制访问规则。`
		success += "符合"
		return str, steer, issue, analyze, success
	} else {
		str += `通过sestatus -v查看SELinux status为disabled状态，未开启安全增强型Linux子系统，未对重要主体和客体设置安全标记，未实现通过安全标记控制主体对客体的访问。`
		issue += "未对重要主体和客体设置安全标记，未实现通过安全标记控制主体对信息资源的访问。"
		analyze += "恶意用户可能通过修改用户权限等方法，非授权访问重要信息资源，存在潜在的安全隐患。"
		steer += "建议对资源进行严格划分，并对重要主体和客体进行分级标记，形成完整的资源分级和访问权限控制结构体系，依据安全标记控制主体对信息资源的访问。"
		success += "不符合"
		return str, steer, issue, analyze, success
	}
}
