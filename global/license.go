package global

import (
	"encoding/base64"
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

// 临时授权
var LicenseList_temporary = LicenseName{Name: "河北华测信息技术有限公司", UUID: "Z2dnZ2dnLXl5eXl5eS1zc3Nzc3MtMDAwMDAwLTAwMDAwMA==", ExpiryDate: parseDate("2023-07-01")}

// Checkactivation 先判断是否有激活文件
func Checkactivation() {
	//先读取lic文件
	if PathExists(Licensename) {
		uuidbyte, _ := os.ReadFile(Licensename)
		sDec := base64.StdEncoding.EncodeToString(uuidbyte)
		if sDec != LicenseList_temporary.UUID || isExpired(LicenseList_temporary.ExpiryDate) {
			message := fmt.Sprintf("此设备未激活,激活后再来吧！")
			waring(message)
		}
		fmt.Printf("激活账户：%s,有效期至: %s\n", LicenseList_temporary.Name, LicenseList_temporary.ExpiryDate)
		return
	}
	waring("此设备未激活,激活后再来吧！")

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

func waring(message string) {
	title := "许可提示"
	if runtime.GOOS == "windows" {
		command := fmt.Sprintf("Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.MessageBox]::Show('%s', '%s')", message, title)
		exec.Command("powershell.exe", "-Command", command).Run()
	} else {
		fmt.Println(message)
	}
	os.Exit(1)
}
