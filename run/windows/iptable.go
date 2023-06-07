//go:build windows

package windows

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"strings"
)

func iptables() {
	var checkiptable []Policyone
	onePolicyone := Policyone{}
	echo := ""
	//域防火墙状态
	domainState, err := getFirewallProfileState("DomainProfile")
	if err == nil {
		if domainState == "Enabled" {
			onePolicyone = Policyone{Name: "核查域网络防火墙状态", Value: "已开启", Static: Yes, Steer: "开启"}
		} else {
			onePolicyone = Policyone{Name: "核查域网络防火墙状态", Value: "未开启", Static: No, Steer: "未启"}
		}
	} else {
		onePolicyone = Policyone{Name: "核查域网络防火墙状态", Value: "未开启", Static: No, Steer: "开启"}
	}
	checkiptable = append(checkiptable, onePolicyone)
	//专用网络防火墙状态
	privateState, err := getFirewallProfileState("StandardProfile")
	if err == nil {
		if privateState == "Enabled" {
			onePolicyone = Policyone{Name: "核查专用网络防防火墙状态", Value: "已开启", Static: Yes, Steer: "开启"}
		}
	} else {
		onePolicyone = Policyone{Name: "核查专用网络防防火墙状态", Value: "未开启", Static: No, Steer: "开启"}
	}
	checkiptable = append(checkiptable, onePolicyone)
	//公共网络防火墙状态
	publicState, err := getFirewallProfileState("PublicProfile")
	if err == nil {
		if publicState == "Enabled" {
			onePolicyone = Policyone{Name: "核查公共网络防火墙状态", Value: "已开启", Static: Yes, Steer: "开启"}
		}
	} else {
		onePolicyone = Policyone{Name: "核查公共网络防火墙状态", Value: "未开启", Static: No, Steer: "开启"}
	}
	checkiptable = append(checkiptable, onePolicyone)
	for _, v := range checkiptable {
		echo += fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>", v.Name, v.Value, v.Static, v.Steer)
	}
	html = strings.ReplaceAll(html, "防火墙状态检查结果", echo)

	//域防火墙规则
	domainrlue := ExecCommands("netsh advfirewall firewall show rule name=all profile=domain")
	html = strings.ReplaceAll(html, "域防火墙规则结果", domainrlue)
	//专网防火墙规则
	privaterlue := ExecCommands("netsh advfirewall firewall show rule name=all profile=private")
	html = strings.ReplaceAll(html, "专网防火墙规则结果", privaterlue)
	//公共防火墙规则
	publicrlue := ExecCommands("netsh advfirewall firewall show rule name=all profile=public")
	html = strings.ReplaceAll(html, "公共防火墙规则结果", publicrlue)

}

// getFirewallProfileState "SYSTEM\\CurrentControlSet\\Services\\SharedAccess\\Parameters\\FirewallPolicy\\ 下DomainProfile、StandardProfile、PublicProfile
func getFirewallProfileState(profile string) (string, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Services\\SharedAccess\\Parameters\\FirewallPolicy\\"+profile, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer key.Close()

	enabled, _, err := key.GetIntegerValue("EnableFirewall")
	if err != nil {
		return "", err
	}

	if enabled == 0 {
		return "Disabled", nil
	}
	return "Enabled", nil
}
