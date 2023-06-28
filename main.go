package main

// Copyright © 2023 NAME HERE <selinuxg@163.com>

import (
	"fmt"
	"golin/cmd"
	"time"
)

func main() {
	start := time.Now()
	//global.Checkactivation() 验证许可
	cmd.Execute()
	end := time.Now().Sub(start)
	fmt.Printf("\n[*] 任务结束,耗时: %s\n", end)
}
