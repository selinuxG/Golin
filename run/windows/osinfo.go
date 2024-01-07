//go:build windows

package windows

import (
	"fmt"
	"regexp"
	"strings"
)

func osinfo() {
	name := cmdvalue(`wmic os get Caption /value`)
	Version := cmdvalue(`wmic os get Version /value`)
	//获取架构后，正则匹配数字，设置代码页936后位字乱码
	arch := cmdvalue(`wmic os get OSArchitecture /value`)
	re := regexp.MustCompile(`\d+`)
	arch = re.FindAllString(arch, -1)[0]
	InstallDate := cmdvalue(`wmic os get InstallDate /value`)
	html = strings.ReplaceAll(html, "操作系统详细信息", fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td</tr>", ExecCommandsPowershll("hostname"), name, Version, arch, InstallDate, ExecCommandsPowershll(` typeperf "\Processor(_Total)\% Processor Time" -sc 1 | Select-String -Pattern '"(\d+\.\d+)"\s*' | ForEach-Object { [math]::Round([double]$_.Matches.Groups[1].Value, 2) }`), ExecCommandsPowershll(`typeperf "\Memory\% Committed Bytes In Use" -sc 1 | Select-String -Pattern '"(\d+\.\d+)"\s*' | ForEach-Object { [math]::Round([double]$_.Matches.Groups[1].Value, 2) }`), ExecCommandsPowershll("date")))
	systeminfo := ExecCommands("systeminfo")
	html = strings.ReplaceAll(html, "系统信息结果", systeminfo)

}

func cmdvalue(cmd string) string {
	v := ExecCommands(cmd)
	v = strings.ReplaceAll(v, "\r\r\n", "")
	vlsit := strings.Split(v, "=")
	if len(vlsit) == 2 {
		return vlsit[1]
	}
	return ""

}
