//go:build windows

package windows

import (
	"fmt"
	"strconv"
	"strings"
)

// screen 屏保相关
func screen() {
	echo := fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", "检查屏幕保护开启状态", "未开启", No, "开启并配置等待时间不高于10分钟")
	//查询是否开启屏保
	v, err := regedt(`Control Panel\Desktop`, "SCRNSAVE.EXE")
	if err == nil {
		v, err = regedt(`Control Panel\Desktop`, "ScreenSaveTimeOut")
		if err == nil {
			intv, _ := strconv.Atoi(v)
			if intv <= 600 {
				echo = fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", "检查屏幕保护开启状态", "已开启", Yes, "开启并配置等待时间不高于10分钟")
				echo += fmt.Sprintf("<tr><td>%s</td><td>%s(秒)</td><td>%s</td><td>%s</td></tr>", "检查屏幕保护等待时间", v, Yes, "等待时间不高于10分钟")
			}
		}
		v, err = regedt(`Control Panel\Desktop`, "ScreenSaverIsSecure")
		if err == nil {
			intv, _ := strconv.Atoi(v)
			if intv > 0 {
				echo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", "检查在恢复时显示登录界面", "已开启", Yes, "开启")
			}
			if intv == 0 {
				echo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", "检查在恢复时显示登录界面", "未启用", No, "开启")
			}
		}
	}
	html = strings.ReplaceAll(html, "屏幕保护相关结果", echo)
}
