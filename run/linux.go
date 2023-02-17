package run

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"time"
)

var echorun bool

func Linux(cmd *cobra.Command, args []string) {

	//确认结果是否输出
	echotype, err := cmd.Flags().GetBool("echo")
	if err != nil {
		fmt.Println(err)
		return
	}
	//读取分隔符
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
			fmt.Printf("\x1b[%dm错误🤷‍ %s自定义执行命令文件不存在！ \x1b[0m\n", 31, cmdpath)
			os.Exit(3)
		}
		fire, _ := ioutil.ReadFile(cmdpath)
		runcmd = string(fire)
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

	//如果value值不为空则是运行一次的模式
	value, err := cmd.Flags().GetString("value")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(value) > 10 {
		Onlyonerun(value, spr, "Linux")
		wg.Wait()
		fmt.Printf("\x1b[%dm✔‍ 单次采集完成，请看「采集完成目录」！ \x1b[0m\n", 34)
		return
	}
	// 下面开始执行批量的
	ippath, err := cmd.Flags().GetString("ip")
	if err != nil {
		fmt.Println(err)
		return
	}
	//判断linux.txt文件是否存在
	Checkfile(ippath, fmt.Sprintf("名称%sip%s用户%s密码%s端口", Split, Split, Split, Split), pem, ippath)
	// 运行share文件中的函数
	Rangefile(ippath, spr, "Linux")
	wg.Wait()
	//完成前最后写入文件
	Deffile("Linux", count, count-len(errhost), errhost)
	fmt.Printf("\x1b[%dm✔‍ 完成! 共采集%d个主机,成功采集%d个主机,失败采集%d个主机。 \x1b[0m\n", 34, count, count-len(errhost), len(errhost))

}

// Runssh 通过调用ssh协议执行命令，写入到文件,并减一个线程数
func Runssh(sshname string, sshHost string, sshUser string, sshPasswrod string, sshPort int, cmd string) {
	defer wg.Done()
	sshType := "password"
	// 创建ssh登录配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second, // ssh连接time out时间一秒钟,如果ssh验证错误会在一秒钟返回
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if sshType == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(sshPasswrod)}
	} else {
		errhost = append(errhost, sshHost)
		return
	}
	// dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {

		errhost = append(errhost, sshHost)
		return
	}
	defer sshClient.Close()

	// 创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		errhost = append(errhost, sshHost)
		return
	}

	defer session.Close()
	// 执行远程命令
	combo, err := session.CombinedOutput(cmd)
	if err != nil {
		errhost = append(errhost, sshHost)
		return
	}

	//判断是否进行输出命令结果
	if echorun {
		fmt.Printf("%s\n%s\n", "<输出结果>", string(combo))
	}

	_, err = os.Stat(succpath)
	if os.IsNotExist(err) {
		os.Mkdir(succpath, pem)
	}
	fire := "采集完成目录//" + sshname + "_" + sshHost + "(linux).log"
	datanew := []byte(string(combo))
	err = ioutil.WriteFile(fire, datanew, pem)
	if err != nil {
		return
	}

}
