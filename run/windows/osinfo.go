package windows

import (
	"fmt"
	"golin/global"
	"strings"
)

func osinfo() {
	name := cmdvalue(`wmic os get Caption /value`)
	Version := cmdvalue(`wmic os get Version /value`)
	arch := cmdvalue(`wmic os get OSArchitecture /value`)
	InstallDate := cmdvalue(`wmic os get InstallDate /value`)
	html = strings.ReplaceAll(html, "操作系统详细信息", fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", name, Version, arch, InstallDate))

}

func cmdvalue(cmd string) string {
	v := global.ExecCommands(cmd)
	v = strings.ReplaceAll(v, "\r\r\n", "")
	vlsit := strings.Split(v, "=")
	if len(vlsit) == 2 {
		return vlsit[1]
	}
	return ""

}
