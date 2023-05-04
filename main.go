/*
Copyright © 2022 高业尚
*/
package main

import (
	"fmt"
	"golin/cmd"
	"time"
)

func main() {
	start := time.Now()
	cmd.Execute()
	end := time.Now().Sub(start)
	fmt.Printf("\n[*] 任务结束,耗时: %s\n", end)
}
