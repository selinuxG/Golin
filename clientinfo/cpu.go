// link:d-eyes

package clientinfo

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"time"
)

type Process struct {
	Process []*process.Process
}

func CPU(cmd *cobra.Command, args []string) {

	count, _ := cmd.Flags().GetInt("count")
	Top(count)

}

func Top(count int) {
	ps, err := process.Processes()
	if err != nil {
		fmt.Println(err)
		return
	}

	sort.Slice(
		ps, func(i, j int) bool {
			pic, _ := ps[i].CPUPercent()
			pjc, _ := ps[j].CPUPercent()
			return pic > pjc
		},
	)
	pss := Process{Process: ps}
	CPUSum := 0
	color.Greenp("==============================================================================================\n")
	for _, ps := range pss.Process {
		pid := os.Getpid()
		if pid == int(ps.Pid) {
			continue
		}
		CPUSum++
		color.Greenp("* CPU Top ", CPUSum)
		fmt.Println()
		_pct, _ := ps.CreateTime()
		_pPath, _ := ps.Exe()
		_pCpuP, _ := ps.CPUPercent()
		startDate := time.Unix(_pct/1000, 0).Format("2006-01-02 15:04:05")
		username, _ := ps.Username()
		MemPer, _ := ps.MemoryPercent()
		fmt.Printf(
			"[User]:%s | [Pid]:%d  | [Path]:%s | [CPU]:%.2f | [Memory]:%.2f | [Createdtime]:%v \n",
			username, ps.Pid, _pPath, _pCpuP, MemPer, startDate,
		)
		//network
		_ps, _ := ps.Connections()
		if len(_ps) == 0 {
			fmt.Println("[netstat]: null")
		} else {
			netSum := 0

			for _, conn := range _ps {
				if conn.Family == 1 {
					continue
				}
				netSum++
				fmt.Printf(
					"[netstat %d]: %v:%v<->%v:%v(%v)\n",
					netSum, conn.Laddr.IP, conn.Laddr.Port, conn.Raddr.IP, conn.Raddr.Port, conn.Status,
				)
			}
		}
		color.Greenp("==============================================================================================\n")
		if count != 0 && CPUSum == count {
			break
		}
	}
}
