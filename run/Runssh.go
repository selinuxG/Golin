package run

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golin/global"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Runssh 通过调用ssh协议执行命令，写入到文件,并减一个线程数
func Runssh(sshname string, sshHost string, sshUser string, sshPasswrod string, sshPort int, cmd string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("ssh panic: %v", r)
		}
		wg.Done()
	}()
	// 创建ssh登录配置
	configssh := &ssh.ClientConfig{
		Timeout:         time.Second * 3, // ssh连接timeout时间
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	configssh.Auth = []ssh.AuthMethod{ssh.Password(sshPasswrod)}
	//增加旧版本算法支持,部分机器会出现 ssh: handshake failed: ssh: packet too large 报错
	//configssh.Ciphers = []string{"aes128-cbc", "aes256-cbc", "3des-cbc", "aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "chacha20-poly1305@openssh.com"}

	// dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, err := ssh.Dial("tcp", addr, configssh)
	if err != nil {
		return fmt.Errorf("%s:连接失败,请检查密码或网络问题。", sshHost)
	}
	defer sshClient.Close()

	// 创建ssh-session
	session, err := sshClient.NewSession()
	defer session.Close()
	if err != nil {
		return fmt.Errorf("%s:创建session失败。", sshHost)
	}

	firepath := filepath.Join(succpath, "Linux")

	// 自定义命令存在则只执行自定义文件
	if runcmd != "" {
		combo, err := session.CombinedOutput(cmd)
		if err != nil {
			return fmt.Errorf("%s:执行自定义命令失败->%s。", sshHost, cmd)
		}

		//判断是否进行输出命令结果
		if echorun {
			fmt.Printf("%s\n%s\n", "<输出结果>", string(combo))
		}
		datanew := []byte(string(combo))
		err = os.WriteFile(filepath.Join(firepath, fmt.Sprintf("%s_%s.log", sshname, sshHost)), datanew, fs.FileMode(global.FilePer))
		if err != nil {
			return fmt.Errorf("%s:保存文件失败。", sshHost)
		}
		return nil
	}

	// 执行模板文件
	data := Data{
		CSS:  css,
		Name: fmt.Sprintf("%s_%s", sshname, sshHost),
		Info: ServerInfo{
			HostName:    runCmd("hostname", sshClient),
			Arch:        runCmd("arch", sshClient),
			Cpu:         runCmd(`cat /proc/cpuinfo | grep name | sort | uniq|awk -F ":" '{print $2}'| xargs`, sshClient),
			CpuPhysical: runCmd(`cat /proc/cpuinfo |grep "physical id"|sort|uniq| wc -l`, sshClient),
			CpuCore:     runCmd(`cat /proc/cpuinfo | grep "core id" | sort | uniq | wc -l`, sshClient),
			Version:     runCmd("rpm -q centos-release", sshClient),
			ProductName: runCmd(`dmidecode -t system | grep 'Product Name'|awk -F ":" '{print $2}'|xargs `, sshClient),
			Free:        runCmd(`free -g | grep Mem | awk '{print $2}'`, sshClient),
			Ping:        strings.Contains(runCmd(`ping www.baidu.com -c1 -w1 >/dev/null;echo $?`, sshClient), "0"),
		},
		SystemState: SystemInfo{
			Cpu:    runCmd(`top -b -n 1 | grep "Cpu(s)"|awk '{print $8}'`, sshClient),
			Memory: runCmd(`echo "$(free -g | awk 'NR==2 {printf "%.2f", ($3/($2)) * 100}')"`, sshClient),
			Load:   runCmd(`uptime`, sshClient),
			Time:   runCmd("date", sshClient),
		},
		Quality: Pwquality{
			Minlen:      runCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "minlen"|awk -F= '{print $2}'`, sshClient),
			Dcredit:     runCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "dcredit"|awk -F= '{print $2}'`, sshClient),
			Ucredit:     runCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "ucredit"|awk -F= '{print $2}'`, sshClient),
			Lcredit:     runCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "lcredit"|awk -F= '{print $2}'`, sshClient),
			Ocredit:     runCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "ocredit"|awk -F= '{print $2}'`, sshClient),
			Minclass:    runCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "minclass"|awk -F= '{print $2}'`, sshClient),
			Maxrepeat:   runCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "maxrepeat"|awk -F= '{print $2}'`, sshClient),
			Maxsequence: runCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "maxsequence"|awk -F= '{print $2}'`, sshClient),
		},
		Address:       runCmd("ifconfig", sshClient),
		Disk:          runCmd("df -h", sshClient),
		Dns:           runCmd(`cat /etc/resolv.conf|grep -v "^#"|grep -v "^$"`, sshClient),
		PamSSH:        runCmd(`cat /etc/pam.d/sshd|grep -v "^#"|grep -v "^$"`, sshClient),
		SSHAuthorized: runCmd(`cat /root/.ssh/authorized_keys`, sshClient),
		PamSystem:     runCmd(`cat /etc/pam.d/system-auth|grep -v "^#"|grep -v "^$"`, sshClient),
		PamPasswd:     runCmd(`cat /etc/pam.d/passwd|grep -v "^#"|grep -v "^$"`, sshClient),
		PwqualityConf: runCmd(`cat /etc/security/pwquality.conf|grep -v "^#"|grep -v "^$"`, sshClient),
		IptablesInfo:  runCmd("iptables -L", sshClient),
		PS:            runCmd("ps aux", sshClient),
		Sudoers:       runCmd(`cat /etc/sudoers|grep -v "^#"|grep -v "^$"`, sshClient),
		Rsyslog:       runCmd(`cat /etc/rsyslog.conf|grep -v "^#"|grep -v "^$"`, sshClient),
		CronTab:       runCmd(`cat /etc/crontab|grep -v "^#"|grep -v "^$";crontab -l`, sshClient),
		Share:         runCmd(`cat /etc/exports|grep -v "^#"|grep -v "^$"`, sshClient),
		Env:           runCmd(`cat /etc/profile |grep -v "^#" |grep -v "^$"`, sshClient),
		Version:       runCmd("cat /etc/*release;uname -r", sshClient),
		Docker:        runCmd("docker ps", sshClient),
		ListUnit:      runCmd("systemctl list-unit-files|grep enabled", sshClient),
		AuditCtl:      runCmd("auditctl -l", sshClient),
		HeadLog:       runCmd(`head -n 10 /var/log/messages /var/log/secure /var/log/audit/audit.log  /var/log/yum.log /var/log/cron`, sshClient),
		TailLog:       runCmd(`tail -n 10 /var/log/messages /var/log/secure /var/log/audit/audit.log  /var/log/yum.log /var/log/cron`, sshClient),
		Logrotate:     runCmd(`awk 'FNR==1{if(NR!=1)print "\nFile: " FILENAME; else print "File: " FILENAME}{if ($0 !~ /^#/ && $0 !~ /^$/) print $0}' /etc/logrotate.conf /etc/logrotate.d/*`, sshClient),
		RpmInstall:    runCmd("rpm -qa", sshClient),
		HomeLimits:    runCmd("ls -lh /home", sshClient),
		LastLog:       runCmd("lastlog", sshClient),
		User:          make([]LinUser, 0),
		CreateUser:    make([]Logindefs, 0),
		Port:          make([]PortList, 0),
		ConfigSSH:     make([]SSH, 0),
		FilePer:       make([]FileListPer, 0),
		FireWalld:     make([]FireListWalld, 0),
	}

	// 通过/etc/passwd 以及结合chage命令获取用户基本信息
	for _, v := range strings.Split(runCmd("cat /etc/passwd", sshClient), "\n") {
		userinfo := strings.Split(v, ":")
		if len(userinfo) != 7 {
			continue
		}
		var Login bool
		if userinfo[6] == "/bin/bash" || userinfo[6] == "/bin/zsh" {
			Login = true
		}
		shadow := strings.Split(strings.ReplaceAll(runCmd(fmt.Sprintf("chage -l %s", userinfo[0]), sshClient), "：", ":"), "\n")
		if len(shadow) != 8 {
			continue
		}
		passwdstatus := strings.Split(runCmd(fmt.Sprintf("passwd -S %s", userinfo[0]), sshClient), " ")
		passwdremark := strings.Join(passwdstatus[7:], "")
		passwdremark = strings.ReplaceAll(passwdremark, "\n", "")
		user := LinUser{
			Name:          userinfo[0],
			Passwd:        passwdstatus[1],
			Remark:        passwdremark,
			Uid:           userinfo[2],
			Gid:           userinfo[3],
			Pwd:           userinfo[5],
			Bash:          userinfo[6],
			Login:         Login,
			LastPasswd:    strings.Split(shadow[0], ":")[1],
			PasswdExpired: strings.Split(shadow[1], ":")[1],
			UserExpired:   strings.Split(shadow[3], ":")[1],
			MaxPasswd:     strings.Split(shadow[5], ":")[1],
		}
		data.User = append(data.User, user)
	}

	//组信息
	for _, v := range strings.Split(runCmd("cat /etc/group", sshClient), "\n") {
		if len(strings.Split(v, ":")) != 4 {
			continue
		}
		group := LinGroup{
			Name:     strings.Split(v, ":")[0],
			Password: strings.Split(v, ":")[1],
			Gid:      strings.Split(v, ":")[2],
			UserList: strings.Split(v, ":")[3],
		}
		data.Group = append(data.Group, group)
	}

	//读取/etc/login.defs获取新创建用户时的信息
	CreateUserLogindefs := Logindefs{
		PassMaxDays:   runCmd(`cat /etc/login.defs |grep -v "^#" |grep "PASS_MAX_DAYS"|awk -F " " '{print $2}'`, sshClient),
		PassMinDays:   runCmd(`cat /etc/login.defs |grep -v "^#" |grep "PASS_MIN_DAYS"|awk -F " " '{print $2}'`, sshClient),
		PassWarnAge:   runCmd(`cat /etc/login.defs |grep -v "^#" |grep "PASS_WARN_AGE"|awk -F " " '{print $2}'`, sshClient),
		UMASK:         runCmd(`cat /etc/login.defs |grep -v "^#" |grep "UMASK"|awk -F " " '{print $2}'`, sshClient),
		EncryptMethod: runCmd(`cat /etc/login.defs |grep -v "^#" |grep "ENCRYPT_METHOD"|awk -F " " '{print $2}'`, sshClient),
	}
	data.CreateUser = append(data.CreateUser, CreateUserLogindefs)

	//通过ss -tulnp获取端口信息
	for _, p := range strings.Split(runCmd(`ss -tulnp|grep -v "Netid"`, sshClient), "\n") {
		re := regexp.MustCompile(`\s+`)
		s := re.ReplaceAllString(p, " ")
		listen := strings.Split(s, " ")
		if len(listen) != 7 {
			continue
		}
		listenPort := PortList{
			Netid:   listen[0],
			State:   listen[1],
			Local:   listen[4],
			Process: listen[6],
		}
		data.Port = append(data.Port, listenPort)
	}

	//获取sshd_config配置
	sshdconfig := SSHConfig()
	if strings.Contains(runCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "PasswordAuthentication"|awk -F " " '{print $2}'`, sshClient), "no") {
		sshdconfig.PasswordAuthentication = false
	}
	if strings.Contains(runCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "PermitRootLogin"|awk -F " " '{print $2}'`, sshClient), "no") {
		sshdconfig.PermitRootLogin = false
	}
	if strings.Contains(runCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "PermitEmptyPasswords"|awk -F " " '{print $2}'`, sshClient), "yes") {
		sshdconfig.PermitEmptyPasswords = true
	}
	if runCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "Protocol"|awk -F " " '{print $2}'`, sshClient) != "" {
		sshdconfig.Protocol = runCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "Protocol"|awk -F " " '{print $2}'`, sshClient)
	}
	if runCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "MaxAuthTries"|awk -F " " '{print $2}'`, sshClient) != "" {
		sshdconfig.MaxAuthTries = runCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "MaxAuthTries"|awk -F " " '{print $2}'`, sshClient)
	}
	if strings.Contains(runCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "PubkeyAuthentication"|awk -F " " '{print $2}'`, sshClient), "no") {
		sshdconfig.PubkeyAuthentication = false
	}
	data.ConfigSSH = append(data.ConfigSSH, sshdconfig)

	//读取地址限制
	data.HostAllow = runCmd(` cat /etc/hosts.allow |grep -v "^#" |grep -v "^$"`, sshClient)
	data.HostDeny = runCmd(` cat /etc/hosts.Deny |grep -v "^#" |grep -v "^$"`, sshClient)

	//中文文件权限
	var FileList = []string{"/etc/passwd", "/etc/shadow", "/etc/group", "/etc/rsyslog.conf", "/etc/sudoers", "/etc/hosts.allow", "/etc/hosts.deny", "/etc/ssh/sshd_config", "/etc/pam.d/sshd", "/etc/pam.d/passwd", "/var/log/messages", "/var/log/audit/audit.log", "/etc/security/pwquality.conf", "/usr/lib64/security/pam_pwquality.so", "/etc/resolv.conf", "/etc/fstab", "/etc/sysctl.conf", "/etc/selinux/config", "/etc/sysctl.conf", "/etc/audit/auditd.conf"}
	for _, name := range FileList {
		FilePer := FileListPer{
			Name:          name,
			Permission:    runCmd(fmt.Sprintf("stat -c '%s' %s ", "%a", name), sshClient),
			Size:          runCmd(fmt.Sprintf("stat -c '%s' %s ", "%s", name), sshClient),
			Uid:           runCmd(fmt.Sprintf("stat -c '%s' %s ", "%U", name), sshClient),
			Gid:           runCmd(fmt.Sprintf("stat -c '%s' %s ", "%G", name), sshClient),
			LastReadTime:  runCmd(fmt.Sprintf("stat -c '%s' %s ", "%x", name), sshClient),
			LastWriteTime: runCmd(fmt.Sprintf("stat -c '%s' %s ", "%y", name), sshClient),
		}
		data.FilePer = append(data.FilePer, FilePer)
	}

	//防火墙 selinux状态
	data.FireWalld = append(data.FireWalld, FireListWalld{
		Name:   "firewalld",
		Status: runCmd(fmt.Sprintf(`systemctl status firewalld  |grep "Active"|awk -F " " '{print $2}'`), sshClient),
	})
	data.FireWalld = append(data.FireWalld, FireListWalld{
		Name:   "selinux",
		Status: runCmd(fmt.Sprintf(`cat /etc/selinux/config |grep -v "^#"|grep -v "^$"|awk -F 'SELINUX=' '{print $2}'`), sshClient),
	})

	// 读取模板文件
	tmpl, err := template.ParseFS(templateFile, "template/linux_html.html")
	if err != nil {
		return fmt.Errorf("%s:读取HTML模板文件失败。", sshHost)
	}

	//目录是否存在
	_, err = os.Stat(firepath)
	if err != nil {
		_ = os.MkdirAll(firepath, os.FileMode(global.FilePer))
	}

	// 创建一个新的文件
	newFile, err := os.Create(fmt.Sprintf("%s/%s_%s.html", firepath, sshname, sshHost))
	if err != nil {
		return fmt.Errorf("%s:创建文件失败。", sshHost)
	}

	defer newFile.Close()
	// 将模板执行的结果写入新的文件
	err = tmpl.Execute(newFile, data)
	if err != nil {
		return fmt.Errorf("%s:保存模板文件失败。", sshHost)
	}
	return nil
}

func runCmd(cmd string, Client *ssh.Client) string {
	newClient, err := Client.NewSession()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer newClient.Close()
	if sudorun {
		cmd = "sudo " + cmd
		fmt.Println(cmd)
	}
	combo, err := newClient.CombinedOutput(cmd)
	if err != nil {
		return ""
	}
	return string(combo)
}
