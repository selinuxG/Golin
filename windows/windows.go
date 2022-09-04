package windows

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
)

var (
	//go:embed windows.exe
	windosdbcp string
)

func Run() {
	if runtime.GOOS != "windows" {
		log.Println("此模式只支持在windows模式下运行。")
		return
	}
	ioutil.WriteFile("等级保护测评辅助工具标准版.exe", []byte(windosdbcp), 0600)

	log.Println("------即将调用工具：windows测评工具")
	//pwd, _ := os.Getwd()
	cmd := exec.Command("等级保护测评辅助工具标准版.exe")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
