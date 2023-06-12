package crack

import (
	"fmt"
	"github.com/spf13/cobra"
	"golin/global"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	userlist = []string{}
	iplist   = []string{} //暴力破解主机列表
	port     = 0
)

type INFO struct {
	Mode    string   //运行模式
	IP      []string //暴力破解的地址
	User    []string //暴力破解的用户列表
	Passwd  []string //暴力破解的密码列表
	Prot    int      //暴力破解的端口
	NoPing  bool     //是否禁止ping监测
	Chan    int      //并发数量
	Timeout int      //超时等待时常
}

func parseFlags(cmd *cobra.Command) INFO {
	mode := cmd.Use
	checkMode(mode)

	ip, _ := cmd.Flags().GetString("ip")
	if strings.Count(ip, "/") == 0 { //单个IP地址
		if net.ParseIP(ip) == nil {
			fmt.Printf(" [x] IP地址不正确，退出！\n")
			os.Exit(1)
		}
		iplist = append(iplist, ip)
	}
	if strings.Count(ip, "/") == 1 { //范围IP地址
		_, ipnet, _ := net.ParseCIDR(ip)
		for ipl := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ipl); incrementIP(ipl) {
			if len(ipl.To4()) == net.IPv4len {
				lastByte := ipl.To4()[3]
				if lastByte == 0 || lastByte == 255 { //起初0开头以及255结尾
					continue
				}
			}
			iplist = append(iplist, ipl.String())
		}
	}

	valueport, _ := cmd.Flags().GetInt("port") //端口
	if valueport != 0 {
		if valueport > 65535 {
			fmt.Printf(" [x] 端口范围不正确，按照默认端口进行！\n")
		} else {
			port = valueport
		}
	}

	userpath, _ := cmd.Flags().GetString("user") //用户文件
	if global.PathExists(userpath) {
		userdata, _ := os.ReadFile(userpath)
		userstr := strings.ReplaceAll(string(userdata), "\r\n", "\n")
		userlist = userlist[0:0] //清空默认用户
		userlist = append(userlist, strings.Split(userstr, "\n")...)
	}

	passwdpath, _ := cmd.Flags().GetString("passwd") //密码文件
	if global.PathExists(passwdpath) {
		passwddata, _ := os.ReadFile(passwdpath)
		passwdstr := strings.ReplaceAll(string(passwddata), "\r\n", "\n")
		passwdlist = passwdlist[0:0] //清空默认密码
		passwdlist = append(passwdlist, strings.Split(string(passwdstr), "\n")...)
	}

	ping, _ := cmd.Flags().GetBool("noping") //是否禁止ping

	chancount, _ := cmd.Flags().GetInt("chan") //并发数量

	timeout, _ := cmd.Flags().GetInt("time") //超时等待时常
	if timeout <= 0 {
		timeout = 3
	}

	info := INFO{
		Mode:    mode,
		IP:      iplist,
		User:    removeDuplicates(userlist),
		Passwd:  removeDuplicates(passwdlist),
		Prot:    port,
		NoPing:  ping,
		Chan:    chancount,
		Timeout: timeout,
	}

	return info
}

// checkMode 基于模式设置默认端口以及用户列表
func checkMode(mode string) {
	switch mode {
	case "ssh":
		port = 22
		userlist = append(userlist, df_sshuser...)
	case "mysql":
		port = 3306
		userlist = append(userlist, df_mysqluser...)
	case "redis":
		port = 6379
		userlist = append(userlist, df_redisuser...)
	case "pgsql":
		port = 5432
		userlist = append(userlist, df_pgsqluser...)
	case "sqlserver":
		port = 1433
		userlist = append(userlist, df_sqlserveruser...)
	case "ftp":
		port = 21
		userlist = append(userlist, df_ftpuser...)
	}

}

// removeDuplicates 切片去重
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// NetWorkStatus 检查ping
func NetWorkStatus(ip string) bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", ip)
	}

	output, err := cmd.Output()
	if err != nil {
		return false
	}
	outttl := strings.ToLower(string(output)) //所有大写转换为小写
	if strings.Contains(outttl, "ttl") {
		return true
	}
	return false
}

// incrementIP 将ip地址增加1
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
