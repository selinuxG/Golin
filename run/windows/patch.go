//go:build windows

package windows

import (
	"fmt"
	"golin/global"
	"regexp"
	"strings"
)

func patch() {
	//补丁信息
	patchecho := ""
	patch := global.ExecCommands("wmic qfe list")
	patch = strings.ReplaceAll(patch, "\r\n", "\n")
	for i := 1; i < len(strings.Split(patch, "\n")); i++ {
		line := strings.Split(patch, "\n")[i]
		var url, csname, desc, fixid, installedBy, date string
		fmt.Sscanf(line, "%s %s %s %s %s %s", &url, &csname, &desc, &fixid, &installedBy, &date)
		patchecho += fmt.Sprintf("%s %s %s %s %s %s\n", url, csname, desc, fixid, installedBy, date)
		// 创建一个正则表达式匹配多个空格
		re := regexp.MustCompile(`\s{2,}`)
		// 将多个连续空格替换为一个
		patchecho = re.ReplaceAllString(patchecho, " ")
	}
	html = strings.ReplaceAll(html, "补丁相关结果", patchecho) //补丁

}
