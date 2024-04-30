package windows

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/net"
)

func netstat() string {

	td := ""
	connections, err := net.Connections("all")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for _, connection := range connections {
		td += "<tr>"
		switch connection.Type {
		case 1:
			td += fmt.Sprintf("<td>tcp</td>")
		case 2:
			td += fmt.Sprintf("<td>udp</td>")
		}
		td += fmt.Sprintf("<td>%s</td>", connection.Laddr.IP)
		td += fmt.Sprintf("<td>%s</td>", connection.Raddr.IP)
		td += fmt.Sprintf("<td>%s</td>", connection.Status)
		td += fmt.Sprintf("<td>%d</td>", connection.Laddr.Port)
		td += fmt.Sprintf("<td>%d</td>", connection.Pid)
		td += fmt.Sprintf("<td>%s</td>", ExecCommandsPowershll(fmt.Sprintf("Get-Process | Where-Object {$_.Id -eq %d} | Select-Object -ExpandProperty Path", connection.Pid)))
		td += "</tr>"
	}

	return td

}
