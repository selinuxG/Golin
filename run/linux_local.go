package run

import (
	"fmt"
	"go.uber.org/zap"
	"golin/config"
	"golin/global"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
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
			config.Log.Warn("错误", zap.String("本机执行状态", "失败"))
			return
		}
		//执行命令
		if err := cmd.Start(); err != nil {
			fmt.Println("Error:The command is err,", err)
			config.Log.Warn("错误", zap.String("本机执行状态", "失败"))
			return
		}
		//读取所有输出
		bytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			fmt.Println("ReadAll Stdout:", err.Error())
			config.Log.Warn("错误", zap.String("本机执行状态", "失败"))

			return
		}
		if err := cmd.Wait(); err != nil {
			fmt.Println("wait:", err.Error())
			config.Log.Warn("错误", zap.String("本机执行状态", "失败"))
			return
		}

		_, err = os.Stat(succpath)
		if os.IsNotExist(err) {
			os.Mkdir(succpath, os.FileMode(global.FilePer))
		}
		fire := "采集完成目录//" + "localhost" + "(linux).log"
		err = ioutil.WriteFile(fire, bytes, fs.FileMode(global.FilePer))
		if err != nil {
			config.Log.Warn("错误", zap.String("本机执行状态", "失败"))
			return
		}
		//判断是否将结果进行输出
		if echo {
			fmt.Printf("%s\n", string(bytes))
		}
		config.Log.Info("提示", zap.String("本机执行状态", "SUCCES"))
		return
	}
	config.Log.Warn("操作系统错误,无法执行Linux模式", zap.String("当前操作系统", sysType))
	return
}
