package run

import (
	"bytes"
	"fmt"
	"github.com/shirou/gopsutil/v3/docker"
	"go.uber.org/zap"
	"html/template"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func LocalrunLinux() {
	sysType := runtime.GOOS
	if sysType != "linux" {
		zlog.Warn("操作系统错误,无法执行Linux模式", zap.String("当前操作系统", sysType))
		return
	}

	// 执行模板文件
	data := Data{
		CSS:  css,
		Name: runLocalCmd("hostname"),
		Info: ServerInfo{
			HostName:    runLocalCmd("hostname"),
			Arch:        runLocalCmd("arch"),
			Cpu:         runLocalCmd(`cat /proc/cpuinfo | grep name | sort | uniq|awk -F ":" '{print $2}'| xargs`),
			CpuPhysical: runLocalCmd(`cat /proc/cpuinfo |grep "physical id"|sort|uniq| wc -l`),
			CpuCore:     runLocalCmd(`cat /proc/cpuinfo | grep "core id" | sort | uniq | wc -l`),
			Version:     runLocalCmd("rpm -q centos-release"),
			ProductName: runLocalCmd(`dmidecode -t system | grep 'Product Name'|awk -F ":" '{print $2}'|xargs `),
			Free:        runLocalCmd(`free -g | grep Mem | awk '{print $2}'`),
			Ping:        strings.Contains(runLocalCmd(`ping www.baidu.com -c1 -w1 >/dev/null;echo $?`), "0"),
		},
		SystemState: SystemInfo{
			Cpu:    runLocalCmd(`top -b -n 1 | grep "Cpu(s)"|awk '{print $8}'`),
			Memory: runLocalCmd(`echo "$(free -g | awk 'NR==2 {printf "%.2f", ($3/($2)) * 100}')"`),
			Load:   runLocalCmd(`uptime`),
			Time:   runLocalCmd("date"),
		},
		Quality: Pwquality{
			Minlen:      runLocalCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "minlen"|awk -F= '{print $2}'`),
			Dcredit:     runLocalCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "dcredit"|awk -F= '{print $2}'`),
			Ucredit:     runLocalCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "ucredit"|awk -F= '{print $2}'`),
			Lcredit:     runLocalCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "lcredit"|awk -F= '{print $2}'`),
			Ocredit:     runLocalCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "ocredit"|awk -F= '{print $2}'`),
			Minclass:    runLocalCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "minclass"|awk -F= '{print $2}'`),
			Maxrepeat:   runLocalCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "maxrepeat"|awk -F= '{print $2}'`),
			Maxsequence: runLocalCmd(`grep '^[^#]' /etc/security/pwquality.conf |grep "maxsequence"|awk -F= '{print $2}'`),
		},
		Address:       runLocalCmd("ifconfig"),
		Disk:          runLocalCmd("df -h"),
		Dns:           runLocalCmd(`cat /etc/resolv.conf|grep -v "^#"|grep -v "^$"`),
		PamSSH:        runLocalCmd(`cat /etc/pam.d/sshd|grep -v "^#"|grep -v "^$"`),
		SSHAuthorized: runLocalCmd(`cat /root/.ssh/authorized_keys`),
		PamSystem:     runLocalCmd(`cat /etc/pam.d/system-auth|grep -v "^#"|grep -v "^$"`),
		PamPasswd:     runLocalCmd(`cat /etc/pam.d/passwd|grep -v "^#"|grep -v "^$"`),
		PwqualityConf: runLocalCmd(`cat /etc/security/pwquality.conf|grep -v "^#"|grep -v "^$"`),
		IptablesInfo:  runLocalCmd("iptables -L"),
		PS:            runLocalCmd("ps aux"),
		Sudoers:       runLocalCmd(`cat /etc/sudoers|grep -v "^#"|grep -v "^$"`),
		Rsyslog:       runLocalCmd(`cat /etc/rsyslog.conf|grep -v "^#"|grep -v "^$"`),
		CronTab:       runLocalCmd(`cat /etc/crontab|grep -v "^#"|grep -v "^$";crontab -l`),
		Share:         runLocalCmd(`cat /etc/exports|grep -v "^#"|grep -v "^$"`),
		Env:           runLocalCmd(`cat /etc/profile |grep -v "^#" |grep -v "^$"`),
		Version:       runLocalCmd("cat /etc/*release;uname -r"),
		ListUnit:      runLocalCmd("systemctl list-unit-files|grep enabled"),
		AuditCtl:      runLocalCmd("auditctl -l"),
		HeadLog:       runLocalCmd(`head -n 10 /var/log/messages /var/log/secure /var/log/audit/audit.log  /var/log/yum.log /var/log/cron`),
		TailLog:       runLocalCmd(`tail -n 10 /var/log/messages /var/log/secure /var/log/audit/audit.log  /var/log/yum.log /var/log/cron`),
		Logrotate:     runLocalCmd(`awk 'FNR==1{if(NR!=1)print "\nFile: " FILENAME; else print "File: " FILENAME}{if ($0 !~ /^#/ && $0 !~ /^$/) print $0}' /etc/logrotate.conf /etc/logrotate.d/*`),
		RpmInstall:    runLocalCmd("rpm -qa"),
		HomeLimits:    runLocalCmd("ls -lh /home"),
		LastLog:       runLocalCmd("lastlog"),
		User:          make([]LinUser, 0),
		CreateUser:    make([]Logindefs, 0),
		Port:          make([]PortList, 0),
		ConfigSSH:     make([]SSH, 0),
		FilePer:       make([]FileListPer, 0),
		FireWalld:     make([]FireListWalld, 0),
		AuditD: AuditLog{
			Write:   runLocalCmd(`cat /etc/audit/auditd.conf |grep "write_logs" |awk -F "=" '{print $2}'`),
			Logfile: runLocalCmd(`cat /etc/audit/auditd.conf |grep "^log_file" |awk -F "=" '{print $2}'`),
			MaxSize: runLocalCmd(`cat /etc/audit/auditd.conf |grep "max_log_file " |awk -F "=" '{print $2}'`),
			NumLOG:  runLocalCmd(`cat /etc/audit/auditd.conf |grep "num_logs" |awk -F "=" '{print $2}'`),
		},
		DockerServerList: make([]DockerServer, 0),
	}

	data.XZ = func() bool {
		if strings.Contains(data.RpmInstall, "xz-5.6.0") {
			return true
		}
		if strings.Contains(data.RpmInstall, "xz-5.6.1") {
			return true
		}
		return false
	}()

	data.CVE20246387 = CVE_2024_6387Local()
	// 通过/etc/passwd 以及结合chage命令获取用户基本信息
	for _, v := range strings.Split(runLocalCmd("cat /etc/passwd"), "\n") {
		userinfo := strings.Split(v, ":")
		if len(userinfo) != 7 {
			continue
		}
		var Login bool
		if userinfo[6] == "/bin/bash" || userinfo[6] == "/bin/zsh" {
			Login = true
		}
		shadow := strings.Split(strings.ReplaceAll(runLocalCmd(fmt.Sprintf("chage -l %s", userinfo[0])), "：", ":"), "\n")
		if len(shadow) != 8 {
			continue
		}
		passwdstatus := strings.Split(runLocalCmd(fmt.Sprintf("passwd -S %s", userinfo[0])), " ")
		if len(passwdstatus) < 8 {
			continue
		}
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
	for _, v := range strings.Split(runLocalCmd("cat /etc/group"), "\n") {
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
		PassMaxDays:   runLocalCmd(`cat /etc/login.defs |grep -v "^#" |grep "PASS_MAX_DAYS"|awk -F " " '{print $2}'`),
		PassMinDays:   runLocalCmd(`cat /etc/login.defs |grep -v "^#" |grep "PASS_MIN_DAYS"|awk -F " " '{print $2}'`),
		PassWarnAge:   runLocalCmd(`cat /etc/login.defs |grep -v "^#" |grep "PASS_WARN_AGE"|awk -F " " '{print $2}'`),
		UMASK:         runLocalCmd(`cat /etc/login.defs |grep -v "^#" |grep "UMASK"|awk -F " " '{print $2}'`),
		EncryptMethod: runLocalCmd(`cat /etc/login.defs |grep -v "^#" |grep "ENCRYPT_METHOD"|awk -F " " '{print $2}'`),
	}
	data.CreateUser = append(data.CreateUser, CreateUserLogindefs)

	//通过ss -tulnp获取端口信息
	for _, p := range strings.Split(runLocalCmd(`ss -tulnp|grep -v "Netid"`), "\n") {
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
	if strings.Contains(runLocalCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "PasswordAuthentication"|awk -F " " '{print $2}'`), "no") {
		sshdconfig.PasswordAuthentication = false
	}
	if strings.Contains(runLocalCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "PermitRootLogin"|awk -F " " '{print $2}'`), "no") {
		sshdconfig.PermitRootLogin = false
	}
	if strings.Contains(runLocalCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "PermitEmptyPasswords"|awk -F " " '{print $2}'`), "yes") {
		sshdconfig.PermitEmptyPasswords = true
	}
	if runLocalCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "Protocol"|awk -F " " '{print $2}'`) != "" {
		sshdconfig.Protocol = runLocalCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "Protocol"|awk -F " " '{print $2}'`)
	}
	if runLocalCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "MaxAuthTries"|awk -F " " '{print $2}'`) != "" {
		sshdconfig.MaxAuthTries = runLocalCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "MaxAuthTries"|awk -F " " '{print $2}'`)
	}
	if strings.Contains(runLocalCmd(`cat /etc/ssh/sshd_config|grep -v "^#" |grep "PubkeyAuthentication"|awk -F " " '{print $2}'`), "no") {
		sshdconfig.PubkeyAuthentication = false
	}
	data.ConfigSSH = append(data.ConfigSSH, sshdconfig)

	//读取地址限制
	data.HostAllow = runLocalCmd(` cat /etc/hosts.allow |grep -v "^#" |grep -v "^$"`)
	data.HostDeny = runLocalCmd(` cat /etc/hosts.deny |grep -v "^#" |grep -v "^$"`)

	//重要文件权限
	var FileList = []string{"/etc/passwd", "/etc/shadow", "/etc/group", "/etc/rsyslog.conf", "/etc/sudoers", "/etc/hosts.allow", "/etc/hosts.deny", "/etc/ssh/sshd_config", "/etc/pam.d/sshd", "/etc/pam.d/passwd", "/var/log/messages", "/var/log/audit/audit.log", "/etc/security/pwquality.conf", "/usr/lib64/security/pam_pwquality.so", "/etc/resolv.conf", "/etc/fstab", "/etc/sysctl.conf", "/etc/selinux/config", "/etc/sysctl.conf", "/etc/audit/auditd.conf"}
	for _, name := range FileList {
		FilePer := FileListPer{
			Name:          name,
			Permission:    runLocalCmd(fmt.Sprintf("stat -c '%s' %s ", "%a", name)),
			Size:          runLocalCmd(fmt.Sprintf("stat -c '%s' %s ", "%s", name)),
			Uid:           runLocalCmd(fmt.Sprintf("stat -c '%s' %s ", "%U", name)),
			Gid:           runLocalCmd(fmt.Sprintf("stat -c '%s' %s ", "%G", name)),
			LastReadTime:  runLocalCmd(fmt.Sprintf("stat -c '%s' %s ", "%x", name)),
			LastWriteTime: runLocalCmd(fmt.Sprintf("stat -c '%s' %s ", "%y", name)),
		}
		data.FilePer = append(data.FilePer, FilePer)
	}

	//防火墙 selinux状态
	data.FireWalld = append(data.FireWalld, FireListWalld{
		Name:   "firewalld",
		Status: runLocalCmd(fmt.Sprintf(`systemctl status firewalld  |grep "Active"|awk -F " " '{print $2}'`)),
	})
	data.FireWalld = append(data.FireWalld, FireListWalld{
		Name:   "selinux",
		Status: runLocalCmd(fmt.Sprintf(`cat /etc/selinux/config |grep -v "^#"|grep -v "^$"|awk -F 'SELINUX=' '{print $2}'`)),
	})

	//docker信息读取
	dockerList, _ := docker.GetDockerStat()
	for _, v := range dockerList {
		info := DockerServer{
			Id:     v.ContainerID,
			Image:  v.Image,
			Name:   v.Name,
			Status: v.Status,
			Run:    v.Running,
		}
		data.DockerServerList = append(data.DockerServerList, info)
	}

	asString, _ := OutputTemplateAsString(data)
	if _, err := os.Stat("采集完成目录/Linux"); os.IsNotExist(err) {
		_ = os.MkdirAll("采集完成目录/Linux", 0644) //如果目录不存在，使用 MkdirAll 创建目录和任何必需的父目录
	}
	err := os.WriteFile("采集完成目录/Linux/local.html", []byte(asString), 0644)
	if err != nil {
		zlog.Warn("[err]", zap.String("文件", "写入失败!"))
		return
	}
	zlog.Info("[success]", zap.String("本地采集完成", "保存路径:采集完成目录/Linux/local.html"))
}

func runLocalCmd(cmd string) string {
	// 使用exec.Command执行本地命令
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func OutputTemplateAsString(data interface{}) (string, error) {
	// 读取模板文件
	tmpl, err := template.ParseFS(templateFile, "template/linux_html.html")
	if err != nil {
		return "", fmt.Errorf("读取HTML模板文件失败: %v", err)
	}

	// 使用bytes.Buffer来捕获模板的输出
	var tplBuffer bytes.Buffer
	err = tmpl.Execute(&tplBuffer, data)
	if err != nil {
		return "", fmt.Errorf("执行模板失败: %v", err)
	}

	// 将buffer的内容转换为字符串
	result := tplBuffer.String()
	return result, nil
}

func CVE_2024_6387Local() bool {
	out, err := exec.Command("bash", "-c", "ssh -V").CombinedOutput() //ssh -V 大概率是输出到标准错误，用CombinedOutput方法获取
	if err != nil {
		return false
	}
	re := regexp.MustCompile(`OpenSSH_([0-9\.]+[a-zA-Z]?[0-9]*)`)
	matches := re.FindStringSubmatch(string(out))
	if len(matches) < 2 {
		return false
	}

	version := matches[1]
	// Define version bounds
	lowVersion := "8.5p1"
	highVersion := "9.7p1"

	// Compare versions
	if compareVersions(version, lowVersion) >= 0 && compareVersions(version, highVersion) <= 0 {
		return true
	}
	return false
}
