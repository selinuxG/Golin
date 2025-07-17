//go:build windows

package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var WifiCmd = &cobra.Command{
	Use:   "wifi",
	Short: "获取本机WiFi信息",
	Long:  `查看曾经连接过并保存密码的WIFI网络信息,仅支持Windows`,
	Run:   WiFiEcho,
}

func init() {
	rootCmd.AddCommand(WifiCmd)
}

type WifiProfile struct {
	SSID     string
	Auth     string
	Password string
}

func WiFiEcho(cmd *cobra.Command, args []string) {
	profiles, err := GetWifiProfiles()
	if err != nil {
		fmt.Println("错误:", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"SSID", "认证方式", "密码"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
	})

	for _, p := range profiles {
		table.Append([]string{p.SSID, p.Auth, p.Password})
	}

	table.Render()
}

func GetWifiProfiles() ([]WifiProfile, error) {
	cmd := exec.Command("netsh", "wlan", "show", "profiles")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`:\s(.+)`)
	scanner := bufio.NewScanner(bytes.NewReader(output))

	var profiles []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "所有用户配置文件") {
			match := re.FindStringSubmatch(line)
			if len(match) > 1 {
				profiles = append(profiles, strings.TrimSpace(match[1]))
			}
		}
	}

	var results []WifiProfile
	for _, name := range profiles {
		p, err := GetProfileDetail(name)
		if err == nil {
			results = append(results, p)
		}
	}

	return results, nil
}

func GetProfileDetail(name string) (WifiProfile, error) {
	cmd := exec.Command("netsh", "wlan", "show", "profile", fmt.Sprintf("name=%q", name), "key=clear")
	output, err := cmd.Output()
	if err != nil {
		return WifiProfile{}, err
	}

	reKey := regexp.MustCompile(`关键内容\s+:\s(.+)`)
	reAuth := regexp.MustCompile(`身份验证\s+:\s(.+)`)

	scanner := bufio.NewScanner(bytes.NewReader(output))
	var auth, key string

	for scanner.Scan() {
		line := scanner.Text()
		if m := reAuth.FindStringSubmatch(line); len(m) > 1 {
			auth = strings.TrimSpace(m[1])
		}
		if m := reKey.FindStringSubmatch(line); len(m) > 1 {
			key = strings.TrimSpace(m[1])
		}
	}

	return WifiProfile{
		SSID:     name,
		Auth:     auth,
		Password: key,
	}, nil
}
