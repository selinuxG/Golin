package xml

import "strings"

func C3(filestring string) (str, steer, issue, analyze, success string) {
	str = "经核查，"
	//整改建议。
	steer = ""
	//实际问题
	issue = ""
	//危害分析
	analyze = ""
	//符合情况
	success = ""
	if sysbak := strings.Contains(string(filestring), "通信正常"); sysbak {
		str += `通过执行ls -l /var/log中的关键日志权限配置合理仅root具备其权限；通过查看rsyslog.conf中以配置将审计日志推送到日志服务器中,且通信正常,可对审计记录进行保护。`
		success += "符合"
		return str, steer, issue, analyze, success
	} else {

		str += `通过执行ls -l /var/log中的关键日志权限配置合理仅root具备其权限；/通过查看rsyslog.conf中未配置将审计日志推送到日志服务器中,无法对审计记录进行保护。`
		issue += "未定期将审计记录进行备份，存在未预期的删除、修改、覆盖等风险。"
		analyze += "日志记录容易受到恶意篡改、删除，不便于对安全事件进行追溯和分析。"
		steer += "建议部署日志服务器对日志进行集中存放，并确保保存时间能够达到半年以上。"
		success += "部分符合"
		return str, steer, issue, analyze, success
	}
}
