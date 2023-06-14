//go:build linux

package crack

import (
	"context"
	"fmt"
	"os"
)

func rdpcon(ctx context.Context, cancel context.CancelFunc, ip, user, passwd string, port, timeout int) {
	fmt.Print("\033[2K") // 擦除整行
	fmt.Print("\r")      // 光标移动到行首
	fmt.Println("Linux操作系统不支持扫描RDP,使用Windows运行吧！")
	os.Exit(0)
}
