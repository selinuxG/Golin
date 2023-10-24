package run

import "html/template"

type Data struct {
	CSS           template.CSS    //CSS样式
	Name          string          //资产名
	Info          ServerInfo      //服务器基本信息
	User          []LinUser       //现有用户信息
	Group         []LinGroup      //现有组信息
	CreateUser    []Logindefs     //新创建用户时信息
	Quality       Pwquality       ///etc/security/pwquality.conf安全配置
	Port          []PortList      //开放端口信息
	ConfigSSH     []SSH           //sshd_config安全配置信息
	HostAllow     string          ///etc/hosts.allow
	HostDeny      string          ///etc/hosts.Deny
	FilePer       []FileListPer   //重要文件权限
	FireWalld     []FireListWalld //防火墙/selinux
	IptablesInfo  string          //防护墙策略
	Address       string          //网卡信息
	Disk          string          //磁盘信息
	Dns           string          /////etc/resolv.conf文件信息
	PamSSH        string          ///etc/pam.d/sshd文件信息
	PamSystem     string          ///etc/pam.d/system-auth文件信息
	PamPasswd     string          ///etc/pam.d/passwd文件信息
	PwqualityConf string          ///etc/security/pwquality.conf文件信息
	PS            string          //ps aux命令结果
	Sudoers       string          //etc/sudoers文件结果
	Rsyslog       string          //etc/rsyslog.conf文件结果
	CronTab       string          //定时任务
	Share         string          //文件共享
	Env           string          //环境变量
	Version       string          //版本信息
	Docker        string          //docker ps -a
	ListUnit      string          //开机启动项
	HeadLog       string          //前十行日志
	TailLog       string          //后十行日志
	Logrotate     string          //日志切割配置
}

type LinUser struct {
	Name          string //用户
	Passwd        string //密码
	Uid           string //UID
	Gid           string //GID
	Description   string //注释
	Pwd           string //主目录
	Bash          string //命令解释器
	Login         bool   //是否可登录
	LastPasswd    string //上次修改密码时间
	PasswdExpired string //密码过期时间
	Lose          string //失效时间
	UserExpired   string //账户过期时间
	MaxPasswd     string //两次改变密码之间相距的最大天数
}

type LinGroup struct {
	Name     string
	Password string
	Gid      string
	UserList string
}

type Logindefs struct {
	PassMaxDays   string
	PassMinDays   string
	PassWarnAge   string
	UMASK         string
	EncryptMethod string
}

type Pwquality struct {
	Minlen      string
	Dcredit     string
	Ucredit     string
	Lcredit     string
	Ocredit     string
	Minclass    string
	Maxrepeat   string
	Maxsequence string
}

type PortList struct {
	Netid   string //协议
	State   string //状态
	Local   string //监听地址
	Process string //进程信息
}

type SSH struct {
	PermitRootLogin        bool   //是否可以root登录
	PasswordAuthentication bool   //是否允许密码进行验证
	PermitEmptyPasswords   bool   //是否允许空密码进行认证
	Protocol               string //协议
	MaxAuthTries           string //关闭连接之前允许的最大身份验证尝试次数
}

type FileListPer struct {
	Name          string //文件文件
	Permission    string //权限
	Size          string //字节大写
	Uid           string //文件所有者的用户名
	Gid           string //文件所属的组 ID
	LastReadTime  string //最后访问时间
	LastWriteTime string //最后修改时间
}

type FireListWalld struct {
	Name   string
	Status string
}

type ServerInfo struct {
	HostName    string
	Arch        string
	Cpu         string
	CpuPhysical string
	CpuCore     string
	Version     string
	ProductName string
	Free        string
}

// NewSSH 初始化SSH默认配置
func NewSSH() SSH {
	return SSH{
		PermitRootLogin:        true,
		PasswordAuthentication: true,
		PermitEmptyPasswords:   false,
		Protocol:               "SSHV2",
		MaxAuthTries:           "6",
	}
}
