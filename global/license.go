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
	{Name: "k8s-node2", UUID: "1d0c209b-e56d-5506-b217-7b081f1a4d48", ExpiryDate: parseDate("2099-07-01")},
	{Name: "闫士华", UUID: "e71bdbab-27d2-50ec-9550-7f7ac95916e5", ExpiryDate: parseDate("2099-01-01")},
	{Name: "预览服务器", UUID: "fb2ac011-31c0-582b-8163-b01b1885c5cf", ExpiryDate: parseDate("2023-07-01")},
	{Name: "临时许可授权", UUID: "gggggg-yyyyyy-ssssss-000000-000000", ExpiryDate: parseDate("2023-07-01")},
}

// Checkactivation 先判断是否有激活文件，如没有判断在识别是否手动激活了
func Checkactivation() {
	var uuid string
	message := fmt.Sprintf("此设备未激活,激活后再来吧！\nUUID:%s", uuid)
	title := "许可提示"
	if PathExists(Licensename) {
		uuidbyte, _ := os.ReadFile(Licensename)
		uuid = string(uuidbyte)

	} else {
		mac, _ := MacSearch()
		uuid, _ = UuidFromMAC(mac)
	}
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
