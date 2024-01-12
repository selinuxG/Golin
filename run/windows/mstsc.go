//go:build windows

package windows

import (
	"strings"
)

func mstsc() {
	mst := ExecCommandsPowershll(`(Get-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Terminal Server" -Name "fDenyTSConnections").fDenyTSConnections`)
	mst = strings.ReplaceAll(mst, "\r\n", "")
	if mst == "1" {
		mst = Yes
	}
	if mst == "0" {
		mst = No
	}
	html = strings.ReplaceAll(html, "开启远程桌面结果", mst)
	html = strings.ReplaceAll(html, "开启远程桌面端口结果", ExecCommandsPowershll(`(Get-ItemProperty -Path 'HKLM:\System\CurrentControlSet\Control\Terminal Server\WinStations\RDP-Tcp\').PortNumber`))
}
