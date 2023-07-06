//go:build windows

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"golin/global"
	"golin/run/windows"
	"os"
	"os/exec"
	"strings"
)

var batname = "GoLin_win.bat"

// windowsCmd represents the execl command
var windowsCmd = &cobra.Command{
	Use:   "windows",
	Short: "读取安全策略生成html",
	Long:  `读取安全策略生成html`,
	Run: func(cmd *cobra.Command, args []string) {
		switch global.PathExists(batname) {
		case true:
			windows.Windows() //这是主程序
			defer func(name string) {
				_ = os.Remove(name)
			}(batname)
		case false:
			restartWindows()
		default:
		}
	},
}

func init() {
	rootCmd.AddCommand(windowsCmd)
}

// restart_windows 此函数用于生成提权的bat文件，重新执行Windows模式
func restartWindows() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("[-] 获取Golin执行绝对路径失败！\n")
		os.Exit(0)
	}

	batstring := `
@echo off
net session >nul 2>&1
if %errorLevel% == 0 (
    echo Running with admin privileges...
) else (
    echo Requesting admin privileges...
    :: Run this script again with admin privileges
    PowerShell -Command "Start-Process cmd.exe -Verb runAs -ArgumentList '/c %~dpnx0'"
    exit /B
)
cd /D %~dp0
start "" "golinpath" windows
exit /B
`
	batstring = strings.ReplaceAll(batstring, "golinpath", exePath)
	err = global.AppendToFile(batname, batstring)
	if err != nil {
		fmt.Printf("[-] 写入提权bat文件失败！\n")
		os.Exit(0)
	}
	_ = exec.Command("cmd", "/C", batname).Run()
}
