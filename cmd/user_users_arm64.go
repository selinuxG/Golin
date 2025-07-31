//go:build windows && arm64

package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

var userCmd = &cobra.Command{
	Use:   "users",
	Short: "获取用户信息 (Windows ARM64 简化版)",
	Run:   checkusers,
}

func init() {
	rootCmd.AddCommand(userCmd)
}

func checkusers(cmd *cobra.Command, args []string) {
	// 使用 PowerShell 命令获取用户信息，因为 go-win64api 不支持 ARM64
	psCmd := `Get-LocalUser | Select-Object Name, FullName, Enabled, LastLogon, UserMayChangePassword, PasswordRequired | Format-Table -AutoSize`

	execCmd := exec.Command("powershell", "-Command", psCmd)
	output, err := execCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("\033[31m错误：获取用户列表失败 (%v)\033[0m\n", err)
		return
	}

	fmt.Println("Windows ARM64 用户信息 (通过 PowerShell 获取):")
	fmt.Println(strings.TrimSpace(string(output)))

	// 简化的表格显示
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetHeader([]string{
		"说明",
		"状态",
	})

	table.Append([]string{"用户信息获取方式", "PowerShell (ARM64 兼容)"})
	table.Append([]string{"详细信息", "请查看上方 PowerShell 输出"})
	table.Append([]string{"注意", "ARM64 版本功能有限"})

	table.Render()
}
