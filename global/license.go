package global

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type LicenseName struct {
	Name       string    //名称
	UUID       string    //唯一标识
	ExpiryDate time.Time //过期时间
}

var LicenseList = []LicenseName{
	{Name: "高业尚", UUID: "68c78927-8aff-5d8b-9b7e-29052c888215", ExpiryDate: parseDate("2099-01-01")},
	{Name: "预览服务器", UUID: "fb2ac011-31c0-582b-8163-b01b1885c5cf", ExpiryDate: parseDate("2023-07-01")},
}

func Checkactivation() {
	mac, _ := MacSearch()
	uuid, _ := UuidFromMAC(mac)
	message := fmt.Sprintf("此设备未激活,激活后再来吧！\nUUID:%s", uuid)
	title := "许可提示"
	for _, listen := range LicenseList {
		if listen.UUID == uuid {
			if !isExpired(listen.ExpiryDate) {
				fmt.Printf("Licensed to %s ExpiryTime %v \n", listen.Name, listen.ExpiryDate)
				return
			} else {
				message = fmt.Sprintf("此设备已过期！\nUUID:%s", uuid)
			}
		}
	}
	if runtime.GOOS == "windows" {
		command := fmt.Sprintf("Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.MessageBox]::Show('%s', '%s')", message, title)
		exec.Command("powershell.exe", "-Command", command).Run()
	} else {
		fmt.Println(message)
	}
	os.Exit(1)
}

// parseDate 时间字符串解析为time.Time类型
func parseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		fmt.Printf("解析日期字符串时出错: %v\n", err)
	}
	return date
}

// isExpired 判断日期是否过期 true已过期,false未过期
func isExpired(expiryDate time.Time) bool {
	return time.Now().After(expiryDate)
}
