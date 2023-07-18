package crack

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"golin/global"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
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

// 爆破成功的主机列表
type successip struct {
	ip     string
	port   string
	passwd string
}

var (
	userlist    = []string{}
	iplist      = []string{} //暴力破解主机列表
	port        = 0
	fireip      = map[string]int{} //读取文件中的主机列表
	successlist []successip
)

func parseFlags(cmd *cobra.Command) INFO {
	mode := cmd.Use
	checkMode(mode)

	ippath, _ := cmd.Flags().GetString("fire") //主机列表文件
	if global.PathExists(ippath) {
		Readfile(ippath)
	} else {
		ip, _ := cmd.Flags().GetString("ip")
		if strings.Count(ip, "/") == 0 && strings.Count(ip, ":") == 0 { //单个IP地址
			if net.ParseIP(ip) == nil {
				_, err := net.ResolveIPAddr("ip", ip) //ip失败时判断是不是域名
				if err != nil {
					fmt.Printf(" [x] IP地址不正确，退出！\n")
					os.Exit(1)
				}
				//ip = addr.String()
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

		if strings.Count(ip, ":") == 1 { //ip:port格式
			newip := strings.Split(ip, ":")[0]
			nPort := strings.Split(ip, ":")[1]
			atoi, _ := strconv.Atoi(nPort)
			port = atoi
			if net.ParseIP(newip) == nil {
				_, err := net.ResolveIPAddr("ip", newip)
				if err != nil {
					fmt.Printf(" [x] IP地址不正确，退出！\n")
					os.Exit(1)
				}
				//newip = addr.String()
			}
			iplist = append(iplist, newip)
		}
	}

	valueport, _ := cmd.Flags().GetInt("port") //端口
	if valueport != 0 {
		if valueport > 65535 {
			fmt.Printf(" [x] 端口范围不正确，按照默认端口进行！\n")
		} else {
			if port == 0 {
				port = valueport
			}
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

// defaule_port 各个服务默认端口
var defaule_port = map[string]int{
	"ssh":       22,
	"mysql":     3306,
	"redis":     6379,
	"pgsql":     5432,
	"sqlserver": 1433,
	"ftp":       21,
	"smb":       445,
	"telnet":    23,
	"tomcat":    8080,
}

// checkMode 基于模式设置默认端口以及用户列表
func checkMode(mode string) {
	switch mode {
	case "ssh":
		port = defaule_port[mode]
		userlist = append(userlist, df_sshuser...)
	case "mysql":
		port = defaule_port[mode]
		userlist = append(userlist, df_mysqluser...)
	case "redis":
		port = defaule_port[mode]
		userlist = append(userlist, df_redisuser...)
	case "pgsql":
		port = defaule_port[mode]
		userlist = append(userlist, df_pgsqluser...)
	case "sqlserver":
		port = defaule_port[mode]
		userlist = append(userlist, df_sqlserveruser...)
	case "ftp":
		port = defaule_port[mode]
		userlist = append(userlist, df_ftpuser...)
	case "smb":
		port = defaule_port[mode]
		userlist = append(userlist, df_smbuser...)
	case "telnet":
		port = defaule_port[mode]
		userlist = append(userlist, df_telnetuser...)
	case "tomcat":
		port = defaule_port[mode]
		userlist = append(userlist, df_tomcatuser...)
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
		cmd = exec.Command("ping", "-n", "1", "-w", "1", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-w", "1", ip)
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

// IsPortOpen 检查给定主机和端口是否开放，并设置超时时间。
func IsPortOpen(host string, port int, timeout int, mode string) bool {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			return false
		}
		return false
	}

	reader := bufio.NewReader(conn) //数据包返回值

	switch mode {
	case "ssh":
		line, err := reader.ReadString('\n')
		if err != nil {
			return false
		}
		check := strings.HasPrefix(line, "SSH-") // SSH协议要求服务器在建立连接时发送一个类似于SSH-2.0-OpenSSH_7.4的协议标识符，通常位于第一行。
		if check {
			return true
		}
		return false

	case "ftp":
		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}
		if !strings.HasPrefix(response, "220") { // 检查是否以 220 开头，FTP 服务器应该以状态码 220 开头
			return false
		}

	}
	_ = conn.Close()
	return true
}

// Readfile 读取主机列表保存到切片
func Readfile(name string) {
	data, err := os.ReadFile(name)
	if err != nil {
		return
	}
	str := string(data)
	str = strings.ReplaceAll(str, "\r\n", "\n")
	for _, s := range strings.Split(str, "\n") {
		if strings.Count(s, ":") == 1 {
			ip := strings.Split(s, ":")[0]
			intport, err := strconv.Atoi(strings.Split(s, ":")[1])
			if err != nil {
				continue
			}
			if net.ParseIP(ip) == nil {
				_, err := net.ResolveIPAddr("ip", ip)
				if err != nil {
					continue
				}
			}
			if port != 0 && port > 65535 {
				continue
			}
			iplist = append(iplist, ip)
			fireip[ip] = intport
		}
	}
	if len(fireip) == 0 {
		fmt.Printf("[-] 文件：%s 不存在格式正确的主机！")
		os.Exit(0)
	}
}
