//go:build windows

package cmd

import (
	"fmt"
	wapi "github.com/iamacarpet/go-win64api"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

var userCmd = &cobra.Command{
	Use:   "users",
	Short: "获取用户信息",
	Run:   checkusers,
}

func init() {
	rootCmd.AddCommand(userCmd)
}

func checkusers(cmd *cobra.Command, args []string) {
	users, err := wapi.ListLocalUsers()
	if err != nil {
		fmt.Printf("\033[31m错误：获取用户列表失败 (%v)\033[0m\n", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetHeader([]string{
		"用户名",
		"全名",
		"账户状态",
		"管理用户",
		"密码策略",
		"最后登录",
		"用户类型",
	})

	table.SetColWidth(100)       // 增大列宽
	table.SetAutoWrapText(false) // 禁止自动换行
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true) // 启用行间分隔线

	for _, u := range users {
		if u.Username == "" {
			continue
		}

		table.Append([]string{
			u.Username,
			u.FullName,
			formatAccountStatus(u.IsEnabled, u.IsLocked),
			formatAdminStatus(u.IsAdmin),
			formatPasswordPolicy(u.PasswordNeverExpires, u.NoChangePassword),
			formatLastLogon(u.LastLogon),
			detectShadowAccount(u.Username),
		})
	}

	table.Render()
}

func formatAccountStatus(enabled, locked bool) string {
	var status []string
	if !enabled {
		status = append(status, "\033[31m❌禁用")
	}
	if locked {
		status = append(status, "\033[33m⚠️锁定")
	}
	if len(status) == 0 {
		return "\033[32m✅正常"
	}
	return strings.Join(status, "\n")
}

// 管理员状态显示
func formatAdminStatus(isAdmin bool) string {
	if isAdmin {
		return "\033[32m● 是"
	}
	return "\033[37m○ 否"
}

func formatPasswordPolicy(neverExpires, noChange bool) string {
	var policies []string
	if neverExpires {
		policies = append(policies, "\033[31m永不过期")
	}
	if noChange {
		policies = append(policies, "\033[33m不可修改")
	}
	if len(policies) == 0 {
		return "\033[2m-"
	}
	return strings.Join(policies, " | ")
}

func formatLastLogon(t time.Time) string {
	if t.IsZero() || t.Unix() <= 0 {
		return "\033[2m从未登录"
	}
	return t.Local().Format("2006-01-02 15:04")
}

func detectShadowAccount(username string) string {
	if strings.HasSuffix(username, "$") {
		return "\033[31m🚨 影子账户"
	}
	defauleuser := []string{"defaultaccount", "administrator", "WDAGUtilityAccount", "guest"}
	for _, u := range defauleuser {
		if strings.EqualFold(username, u) {
			return "\033[33m⚠️ 默认账户"
		}
	}

	return "\033[2m-"
}
