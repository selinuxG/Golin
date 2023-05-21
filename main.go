package main

// Copyright © 2023 NAME HERE <selinuxg@163.com>

import (
	"fmt"
	"golin/cmd"
	"golin/global"
	"time"
)

func main() {
	start := time.Now()
	global.Checkactivation()
	cmd.Execute()
	end := time.Now().Sub(start)
	fmt.Printf("\n[*] 任务结束,耗时: %s\n", end)
}
