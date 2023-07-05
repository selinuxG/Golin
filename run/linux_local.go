package run

import (
	"fmt"
	"go.uber.org/zap"
	"golin/global"
	"io/fs"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func LocalrunLinux(echo bool, cmd string) {
	sysType := runtime.GOOS
	if sysType == "linux" {
		cmd := exec.Command("/bin/bash", "-c", cmd)
		//创建获取命令输出管道
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
			zlog.Warn("错误", zap.String("本机执行状态", "失败"))
			return
		}
		//执行命令
		if err := cmd.Start(); err != nil {
			fmt.Println("Error:The command is err,", err)
			zlog.Warn("错误", zap.String("本机执行状态", "失败"))
			return
		}
		//读取所有输出
		bytes, err := io.ReadAll(stdout)
		if err != nil {
			fmt.Println("ReadAll Stdout:", err.Error())
			zlog.Warn("错误", zap.String("本机执行状态", "失败"))

			return
		}
		if err := cmd.Wait(); err != nil {
			fmt.Println("wait:", err.Error())
			zlog.Warn("错误", zap.String("本机执行状态", "失败"))
			return
		}

		pullfire := filepath.Join(succpath, "Linux")
		_, err = os.Stat(pullfire)
		if os.IsNotExist(err) {
			os.MkdirAll(pullfire, os.FileMode(global.FilePer))
		}

		fire := filepath.Join(pullfire, "local本机采集")
		err = os.WriteFile(fire, bytes, fs.FileMode(global.FilePer))
		if err != nil {
			zlog.Warn("错误", zap.String("本机执行状态", "失败"))
			return
		}
		//判断是否将结果进行输出
		if echo {
			fmt.Printf("%s\n", string(bytes))
		}
		zlog.Info("提示", zap.String("本机执行状态", "SUCCES"))
		return
	}
	zlog.Warn("操作系统错误,无法执行Linux模式", zap.String("当前操作系统", sysType))
	return
}
