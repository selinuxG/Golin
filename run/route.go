package run

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"golin/global"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	routecmd          []string                                                                                                                                                                                                                                                                                                                                                             //执行的命令列表
	Defroutecmd       = []string{"display version", "display memory", "display  cpu-usage", "display local-user", "display role", "display current-configuration | begin local-user", "display password-control", "display info-center", "dir flash:/logfile/", "display acl all", "display ip interface brief description", "display ssh server status", "display current-configuration"} //默认h3c命令
	DefroutecmdHuawei = []string{"display local-user", "display local-aaa-user password policy access-user", "display aaa configuration"}                                                                                                                                                                                                                                                  //默认华为命令

)

func Route(cmd *cobra.Command, args []string) {
	//判断文件是否存在
	ippath, err := cmd.Flags().GetString("ip")
	if err != nil {
		fmt.Println(err)
		return
	}
	Checkfile(ippath, fmt.Sprintf("名称%sip%s用户%s密码%s端口", Split, Split, Split, Split), global.FilePer, ippath)

	//确认结果是否输出
	echotype, err := cmd.Flags().GetBool("echo")
	if err != nil {
		fmt.Println(err)
		return
	}
	//是否输出记录内容
	echorun = echotype
	//确认分隔符
	spr, err := cmd.Flags().GetString("spript")
	if err != nil {
		fmt.Println(err)
		return
	}
	//是不是自定义文件
	cmdpath, err := cmd.Flags().GetString("cmd")
	if err != nil {
		fmt.Println(err)
		return
	}
	//如果cmdpath不为空，则判断是不是存在，存在则读取出来写入到runcmd变量中，为空则使用函数中的默认命令
	if len(cmdpath) > 0 {
		_, err := os.Stat(cmdpath)
		if os.IsNotExist(err) {
			zlog.Warn("自定义执行命令文件不存在！", zap.String("文件", cmdpath))
			os.Exit(3)
		}
		fire, _ := os.ReadFile(cmdpath)
		runcmd = string(fire)
		for _, v := range strings.Split(string(fire), "\n") {
			routecmd = append(routecmd, v)
		}
	} else {
		routecmd = Defroutecmd
	}

	//判断是否有自定义执行的命令，如果有则处理他，不执行cmd文件中的命令。
	cmdvalue, err := cmd.Flags().GetString("cmdvalue")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(cmdvalue) > 0 {
		if len(strings.Split(cmdvalue, ";")) >= 1 {
			routecmd = []string{}
			for _, v := range strings.Split(cmdvalue, ";") {
				routecmd = append(routecmd, v)
			}
		} else {
			routecmd = []string{}
			routecmd = append(routecmd, cmdvalue)
		}
	}

	// 下面开始执行函数
	fmt.Println("-------------------------------------> run type:route")
	Rourange(ippath, spr, Defroutecmd)
	Deffile("Route", count, count-len(errhost), errhost)
}

func Rourange(path string, spr string, cmd []string) {
	fire, _ := os.ReadFile(path)
	lines := strings.Split(string(fire), "\n")
	for i := 0; i < len(lines); i++ {
		//如果是空行则跳过线程减1
		if len(lines[i]) == 0 {
			continue
		}
		firecount := strings.Count(lines[i], spr)
		if firecount != 4 {
			zlog.Warn("主机格式不正确，跳过！")
			continue
		}
		//总数量+1
		count += 1
		linedata := lines[i]
		Name := strings.Split(string(linedata), spr)[0]
		Host := strings.Split(string(linedata), spr)[1]
		User := strings.Split(string(linedata), spr)[2]
		Passwrod := strings.Split(string(linedata), spr)[3]
		Port1 := strings.Split(string(linedata), spr)[4]
		//windos中换行符可能存在为/r/n,之前分割/n,还留存/r,清除它
		Porttmp := strings.Replace(Port1, "\r", "", -1)
		Port, err := strconv.Atoi(Porttmp)
		if err != nil {
			zlog.Warn("端口转换失败: ", zap.String("IP", Host))
			errhost = append(errhost, Host)
			continue
		}
		//判断host是不是正确的IP地址格式
		address := net.ParseIP(Host)
		if address == nil {
			zlog.Warn("IP地址格式不正确，跳过！", zap.String("IP", Host))
			count = count - 1
			continue
		}
		//判断端口范围是否是1-65535
		if Port == 0 || Port > 65535 {
			zlog.Warn("端口范围不正确，跳过！", zap.String("IP", Host), zap.Int("Port:", Port))
			count = count - 1
			continue
		}
		//如果是Windows先判断保存文件是否存在特殊字符,是的话不执行直接记录为失败主机
		if runtime.GOOS == "windows" {
			if InSlice(denynametype, Name) {
				zlog.Warn("名称存在特殊符号，跳过！")
				errhost = append(errhost, Host)
				continue
			}
		}
		//保存路径是否存在
		firepath := filepath.Join(succpath, "Route")
		_, err = os.Stat(firepath)
		if err != nil {
			os.MkdirAll(firepath, os.FileMode(global.FilePer))
		}

		//拼接后的文件然后删除文件
		filename := fmt.Sprintf("%s_%s.log", Name, Host)
		filename = filepath.Join(firepath, filename)
		os.Remove(filename)
		for _, v := range cmd {
			//命令为空则跳过
			if len(v) == 0 {
				continue
			}
			Routessh(filename, Host, User, Passwrod, strconv.Itoa(Port), v)
		}
	}
}

// Routessh 连接一次执行一次命令。不确认是库本身的问题还是路由设备的问题，缓冲器有问题只能如此。
func Routessh(filename, Host, User, Passwrod, Port, Cmd string) {
	configssh := &ssh.ClientConfig{
		Timeout:         time.Second * 3, // ssh连接timeout时间
		User:            User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	//增加旧版本算法支持
	configssh.KeyExchanges = []string{"diffie-hellman-group-exchange-sha1"}
	configssh.Auth = []ssh.AuthMethod{ssh.Password(Passwrod)}
	//增加旧版本算法支持
	configssh.Ciphers = []string{"aes128-cbc", "aes256-cbc", "3des-cbc", "aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "chacha20-poly1305@openssh.com", "diffie-hellman-group-exchange-sha256", "curve25519-sha256"}
	// dial 获取ssh client
	addr := fmt.Sprintf("%s:%s", Host, Port)
	sshClient, err := ssh.Dial("tcp", addr, configssh)
	if err != nil {
		errhost = append(errhost, Host)
		fmt.Println(err)
		return
	}
	defer sshClient.Close()
	global.AppendToFile(filename, "------------------------------执行命令:"+Cmd+"\n")

	// 创建ssh-session
	var stdoutBuf bytes.Buffer
	session, err := sshClient.NewSession()
	defer session.Close()
	if err != nil {
		fmt.Println(err)
		errhost = append(errhost, Host)
		return
	}
	session.Stdout = &stdoutBuf
	err = session.Run(Cmd)
	if err != nil {
		fmt.Println(err)
		//return
	}
	lines := strings.Split(stdoutBuf.String(), "\n")
	for _, s := range lines {
		if len(s) > 2 && s[0] != '*' && len(strings.Split(s, "*")) < 5 {
			//是否输出
			if echorun {
				fmt.Println(s)
			}
			//替换换行符以及nul 空字符
			s = strings.Replace(s, "\r", "", -1)
			s = strings.Replace(s, "\x00", "", -1)
			s = strings.Replace(s, "\b", "", -1)
			global.AppendToFile(filename, s+"\n")
		}
	}
	global.AppendToFile(filename, "\n\n")
}
