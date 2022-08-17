package xml

import (
	"regexp"
)

func A2(filestring string) (str, steer, issue, analyze, success string) {
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
	re := regexp.MustCompile(`超时退出时间为:(.*?)\n`)
	timout := re.FindAllStringSubmatch(string(filestring), -1)
	//inttimout, _ := strconv.Atoi(timout[0][1])
	if len(timout) == 1 {
		errPASScoung += 1
		str += "执行命令echo $TMOUT得到该变量值为空，不满足超时连接自动退出要求;"
		issue += "未配置登录连接超时策略;"
		steer += "建议配置登录失败处理策略;"
		analyze += "设备易被非授权人员恶意操作，存在非授权访问的风险。"
	} else {
		str += "执行命令echo $TMOUT得到该变量值为:" + timout[0][1] + "满足超时连接自动退出要求;"
	}
	re = regexp.MustCompile(`登陆错误失败次数为:(.*?)\n`)
	ueerlock := re.FindAllStringSubmatch(string(filestring), -1)
	if len(ueerlock) == 1 {
		errPASScoung += 1
		str += "/etc/pam.d/sshd中未配置deny锁定策略。"
		steer += "建议配置用户登陆失败锁定次数策略;"
		analyze += "建议配置登录失败处理策略，配置登录失败5次，锁定10分钟左右，防止恶意人员暴力破解账户口令。并配置登录连接超时策略，配置超时时间10分钟左右，降低设备被非授权访问的风险。"
	} else {
		re = regexp.MustCompile(`失败锁定时常为:(.*?)\n`)
		ueerlocktime := re.FindAllStringSubmatch(string(filestring), -1)
		str += "/etc/pam.d/sshd中以配置失败" + ueerlock[0][1] + "次后锁定" + ueerlocktime[0][1] + "分钟。"
	}

	switch {
	case errPASScoung == 0:
		success += "符合"
	case errPASScoung == 2:
		success += "不符合"
	default:
		success += "部分符合"
	}
	return str, steer, issue, analyze, success
}
