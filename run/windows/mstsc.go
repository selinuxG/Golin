//go:build windows

package windows

import (
	"golin/global"
	"strings"
)

func mstsc() {
	mst := global.ExecCommandsPowershll(`(Get-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Terminal Server" -Name "fDenyTSConnections").fDenyTSConnections`)
	mst = strings.ReplaceAll(mst, "\r\n", "")
	if mst == "1" {
		mst = Yes
	}
	if mst == "0" {
		mst = No
	}
	html = strings.ReplaceAll(html, "开启远程桌面结果", mst)
}
