//go:build windows

package windows

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"golin/global"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"
)

var (
	ch       = make(chan struct{}, 3) //并发数量
	wg       = sync.WaitGroup{}
	Policy   = make(map[string]string) //安全策略map
	html     = Windowshtml()           //html字符串
	mu       sync.Mutex                //加锁
	auditmap = map[string]string{      //审计相关
		"AuditSystemEvents":           "是否审核系统事件",
		"AuditLogonEvents":            "是否审核登录事件",
		"AuditObjectAccess":           "是否审核对象访问事件",
		"AuditPrivilegeUse":           "是否审核权限使用事件",
		"AuditPolicyChange":           "是否审核策略更改事件",
		"AuditAccountManage":          "是否审核账户管理事件",
		"AuditDirectoryServiceAccess": "是否审核目录服务访问事件",
		"AuditAccountLogon":           "是否审核帐户登录事件",
	}
)

const (
	Yes = "是"
	No  = "否"
)

// replaceCommand 执行cmd或者powershell替换内容
type replaceCommand struct {
	Placeholder string
	Command     string
	Powershell  bool
}

type Policyone struct {
	Name   string //检查项
	Value  string //当前值
	Static string //状态
	Steer  string //建议
}

func Windows() {

	//使用 secedit 工具导出本地安全策略
	cmd := exec.Command("secedit", "/export", "/cfg", "policy.txt")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running secedit:", err)
		fmt.Println("是否用管理员执行了？")
		return
	}
	defer os.Remove("policy.txt")

	// 读取导出的安全策略文件
	data, err := os.ReadFile("policy.txt")
	if err != nil {
		fmt.Println("Error reading policy file:", err)
		return
	}

	// 查找 PasswordComplexity 设置项
	content := string(data)
	content = strings.Replace(content, "\x00", "", -1) //删除nul空字符
	content = strings.Replace(content, " ", "", -1)    //删除空格

	for _, i2 := range strings.Split(content, "\r\n") {
		if strings.Count(i2, "=") == 1 {
			Policy[strings.Split(i2, "=")[0]] = strings.Split(i2, "=")[1]
		}
	}

	runserver := []func(){mstsc, osinfo, iptables, usercheck, checkpasswd, lock, auditd, screen, patch, iptables, disk}
	wg.Add(len(runserver))
	commands := []replaceCommand{
		{"端口相关结果", "netstat -ano", false},
		{"进程列表结果", `tasklist | sort`, false},
		{"定时任务结果", `schtasks /query /fo LIST`, false},
		{"安装组件结果", `Dism /online /Get-Features /format:table`, true},
		{"安装程序结果", `Get-ItemProperty HKLM:\Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\*, HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\*, HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\* | Select-Object DisplayName, DisplayVersion, Publisher, InstallDate | Where-Object { $_.DisplayName -ne $null } | Sort-Object DisplayName -Unique | Format-Table –AutoSize`, true},
		{"Service结果", "Get-Service | Format-Table -AutoSize", true},
		{"共享资源结果", "net share", false},
		{"联网测试结果", "ping www.baidu.com", false},
		{"群组信息结果", `Get-CimInstance -ClassName Win32_Group | Select-Object Name, SID, Description | Format-List`, true},
		{"防病毒结果", `Get-MpComputerStatus`, true},
		{"安装驱动结果", `driverquery`, false},
		{"环境变量结果", `set`, false},
		{"家目录权限结果", `Get-ChildItem "C:\Users\" -Directory | ForEach-Object { $acl = Get-Acl $_.FullName; $properties = @{Path=$_.FullName; Owner=$acl.Owner; Access=$acl.Access | ForEach-Object { $_.FileSystemRights.ToString() + ', ' + $_.AccessControlType.ToString() + ', ' + $_.IdentityReference.ToString() }}; New-Object -TypeName PSCustomObject -Property $properties } | Format-List`, true},
		{"开机启动结果", `wmic startup get caption,command`, true},
		{"日志信息结果", `Get-ChildItem "C:\Windows\System32\winevt\Logs"`, true},
	}
	wg.Add(len(commands))
	for _, cmd := range commands {
		ch <- struct{}{}
		go replaceAsync(&html, cmd, &wg)
	}

	for _, v := range runserver {
		ch <- struct{}{}
		go func(v func()) {
			defer wg.Done()
			defer func() { <-ch }()
			v()
		}(v)
	}
	wg.Wait()

	//驱动

	//给结果增加颜色并写入文件
	html = strings.ReplaceAll(html, "<td>是</td>", `<td style="color: rgb(32, 199, 29)">是</td>`)
	html = strings.ReplaceAll(html, "<td>否</td>", `<td style="color: rgb(255, 0, 0)">否</td>`)
	html = strings.ReplaceAll(html, "生成日期", fmt.Sprintf("%s", time.Now().Format(time.DateTime)))
	os.Remove("windows.html")
	os.WriteFile("windows.html", []byte(html), os.FileMode(global.FilePer))
	if global.PathExists("windows.html") {
		_ = ExecCommands("start windows.html")
	}

	return
}

// replaceAsync 执行命令或替换模板结果
func replaceAsync(html *string, cmd replaceCommand, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() { <-ch }()

	// 因为等的有点久，所有增加一个输出
	fmt.Printf("[*] 正在执行命令:%s 用于获取:%s powershell：%v\n", cmd.Command, cmd.Placeholder, cmd.Powershell)

	var result string
	if cmd.Powershell {
		result = ExecCommandsPowershll(cmd.Command)
	} else {
		result = ExecCommands(cmd.Command)
	}
	mu.Lock()
	*html = strings.ReplaceAll(*html, cmd.Placeholder, result)
	mu.Unlock()
}

// ExecCommands 执行cmd命令
func ExecCommands(commands ...string) string {
	cmd := strings.Join(commands, " && ")
	execCmd := exec.Command("cmd", "/C", cmd)
	execCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} //关闭弹窗
	out, err := execCmd.CombinedOutput()
	if err != nil {
		return ""
	}
	var output []byte
	// 检查输出是否为有效的 UTF-8 编码
	if utf8.Valid(out) {
		output = out
	} else {
		output, _ = gbkToUtf8(out)
	}
	//output, _ := gbkToUtf8(out)
	return string(output)
}

// gbkToUtf8 转码utf8
func gbkToUtf8(input []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(input), simplifiedchinese.GBK.NewDecoder())
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(reader)
	if err != nil {
		return nil, fmt.Errorf("转换编码时发生错误: %v", err)
	}
	return buffer.Bytes(), nil
}

// ExecCommandsPowershll 执行Powershll命令
func ExecCommandsPowershll(commands ...string) string {
	cmdLine := strings.Join(commands, " ; ")
	execCmd := exec.Command("powershell", "-Command", cmdLine)
	execCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := exec.Command("powershell", "-Command", cmdLine).CombinedOutput()
	if err != nil {
		return ""
	}
	var output []byte
	// 检查输出是否为有效的 UTF-8 编码
	if utf8.Valid(out) {
		output = out
	} else {
		output, _ = gbkToUtf8(out)
	}
	return string(output)
}

/*
	mstsc()       //远程桌面
	osinfo()      //操作系统信息
	iptables()    //防火墙状态核查结果
	usercheck()   //用户详细信息
	checkpasswd() //密码策略
	lock()        //失败锁定策略
	auditd()      //日志策略
	screen()      //屏幕锁定策略
	patch()       //补丁信息
	iptables()    //防火墙状态核查结果
	disk()        //磁盘信息
	ExecCommands 执行cmd命令
*/
