package windows

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Run() {
	log.Println("------即将调用华测工具：windows测评工具")
	pwd, _ := os.Getwd()
	//cmd := exec.Command("cmd.exe", "/c", "start "+pwd+"//windows//dbcp-windos.exe")
	cmd := exec.Command("cmd.exe", "/c", "start "+pwd+"//windows//等级保护测评辅助工具标准版.exe")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
