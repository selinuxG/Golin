package windows

import (
	"fmt"
	"golin/global"
	"os"
	"os/exec"
	"strings"
)

var (
	Policy   = make(map[string]string) //安全策略map
	html     = Windowshtml()           //html字符串
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

type Policyone struct {
	Name   string //检查项
	Value  string //但当前值
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
	osinfo()                                                                       //操作系统信息
	iptables()                                                                     //防火墙状态核查结果
	usercheck()                                                                    //用户详细信息
	checkpasswd()                                                                  //密码策略
	lock()                                                                         //失败锁定策略
	auditd()                                                                       //日志策略
	screen()                                                                       //屏幕锁定策略
	patch()                                                                        //补丁信息
	iptables()                                                                     //防火墙状态核查结果
	html = strings.ReplaceAll(html, "端口相关结果", global.ExecCommands("netstat -ano")) //开放端口

	//给结果增加颜色并写入文件
	html = strings.ReplaceAll(html, "<td>是</td>", `<td style="color: rgb(32, 199, 29)">是</td>`)
	html = strings.ReplaceAll(html, "<td>否</td>", `<td style="color: rgb(255, 0, 0)">否</td>`)
	os.Remove("windows.html")
	os.WriteFile("windows.html", []byte(html), os.FileMode(global.FilePer))
}
