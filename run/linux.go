package run

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golin/global"
	"os"
	"path"
	"strings"
)

var (
	echorun bool
	sudorun bool
)

func Linux(cmd *cobra.Command, args []string) {
	_ = global.MkdirAll(path.Join("采集完成目录", "Linux"))
	//是否sudo权限执行
	sudo, err := cmd.Flags().GetBool("sudo")
	if err != nil {
		fmt.Println(err)
		return
	}
	sudorun = sudo

	//确认结果是否输出
	echotype, err := cmd.Flags().GetBool("echo")
	if err != nil {
		fmt.Println(err)
		return
	}
	//是否输出记录内容
	echorun = echotype
	spr, err := cmd.Flags().GetString("spript")
	if err != nil {
		fmt.Println(err)
		return
	}
	cmdpath, err := cmd.Flags().GetString("cmd")
	if err != nil {
		fmt.Println(err)
		return
	}
	//如果cmdpath不为空，则判断是不是存在，存在则读取出来写入到runcmd变量中，为空则使用 Linux_cmd函数中的默认命令
	if len(cmdpath) > 0 {
		_, err := os.Stat(cmdpath)
		if os.IsNotExist(err) {
			zlog.Warn("自定义执行命令文件不存在！", zap.String("文件", cmdpath))
			os.Exit(3)
		}
		fire, _ := os.ReadFile(cmdpath)
		runcmd = string(fire)
		//新增： 去掉文件中的换行符，最后一个不是；自动增加然后保存成一条命令
		newcmd := ""
		checkcmd := strings.Split(runcmd, "\n")
		for i := 0; i < len(checkcmd); i++ {
			checkend := checkcmd[i]
			checkend = strings.Replace(checkend, "\r", "", -1)
			if len(checkend) == 0 {
				continue
			}
			if checkend[len(checkend)-1:] == ";" {
				newcmd += checkend
			} else {
				newcmd += checkend + ";"
			}
		}
		if newcmd != "" {
			runcmd = newcmd
		}
	}

	//判断是否有自定义执行的命令，如果有则处理他，不执行cmd文件中的命令。
	cmdvalue, err := cmd.Flags().GetString("cmdvalue")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(cmdvalue) > 0 {
		runcmd = string(cmdvalue)
	}
	//判断是不是本机执行的模式，
	localhosttype, err := cmd.Flags().GetBool("localhost")
	if err != nil {
		fmt.Println(err)
		return
	}
	if localhosttype {
		LocalrunLinux()
		return
	}

	//如果value值不为空则是运行一次的模式
	value, err := cmd.Flags().GetString("value")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(value) > 10 {
		Onlyonerun(value, spr, "Linux")
		wg.Wait()
		return
	}
	// 下面开始执行批量的
	ippath, err := cmd.Flags().GetString("ip")
	if err != nil {
		fmt.Println(err)
		return
	}
	//判断linux.txt文件是否存在
	Checkfile(ippath, fmt.Sprintf("名称%sip%s用户%s密码%s端口", Split, Split, Split, Split), global.FilePer, ippath)
	// 运行share文件中的函数
	Rangefile(ippath, spr, "Linux")
	wg.Wait()
	//完成前最后写入文件
	Deffile("Linux", count, count-len(errhost), errhost)
}
