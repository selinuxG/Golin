package osinfo

import (
	"fmt"
	"os/user"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func Osinfo() {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Info()
	cc, _ := cpu.Percent(time.Second, false)
	d, _ := disk.Usage("/")
	n, _ := host.Info()

	fmt.Printf("        内存信息: 共%v MB  剩余: %v MB 使用:%v MB Usage:%f%%\n", v.Total/1024/1024, v.Available/1024/1024, v.Used/1024/1024, v.UsedPercent)
	if len(c) > 1 {
		for _, sub_cpu := range c {
			modelname := sub_cpu.ModelName
			cores := sub_cpu.Cores
			fmt.Printf("        CPU       : %v   %v cores \n", modelname, cores)
		}
	} else {
		sub_cpu := c[0]
		modelname := sub_cpu.ModelName
		cores := sub_cpu.Cores
		fmt.Printf("        CPU信息: %v   %v cores \n", modelname, cores)
	}
	fmt.Printf("        CPU使用率: used %f%% \n", cc[0])
	fmt.Printf("        硬盘信息: 共%v GB  剩余: %v GB Usage:%f%%\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)
	fmt.Printf("        系统版本: %v(%v)   %v  \n", n.Platform, n.PlatformFamily, n.PlatformVersion)
	fmt.Println("        系统架构:", runtime.GOARCH)
	fmt.Printf("        主机名称: %v  \n", n.Hostname)
	user, pwd := use(), pwd()
	fmt.Println("        当前用户:", user, " 家目录:", pwd)
}

func use() string {
	username, err := user.Current()
	if err != nil {
		fmt.Println(err.Error())
	}
	user := username.Username
	return user

}

func pwd() string {
	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}

	return u.HomeDir

}
