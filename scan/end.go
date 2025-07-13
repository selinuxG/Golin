package scan

import (
	"fmt"
	"golin/global"
	"golin/poc"
	"golin/scan/crack"
	"html/template"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// end 运行结束是输出,输出一些统计信息
func endEcho() {
	vulnerablehost, _ := calculateVulnerablePercentage(iplist, crack.MapCrackHost, poc.ListPocInfo)
	fmt.Printf("\r+------------------------------------------------------------+\n"+
		"[*] 漏洞主机:%v Linux:%v Windows:%v 存活端口:%v ssh:%v rdp:%v web:%v 指纹:%v 数据库:%v 弱口令:%v 漏洞:%v \n",
		printRed("%v%v", vulnerablehost, printGreen("/"+strconv.Itoa(len(iplist)))),
		printGreen("%v", linuxcount),
		printGreen("%v", windowscount),
		printGreen("%v", len(infolist)),
		printGreen("%v", protocolExistsAndCount("ssh")),
		printGreen("%v", protocolExistsAndCount("rdp")),
		printGreen("%v", protocolExistsAndCount("WEB应用")),
		printGreen("%v", len(global.AppMatchedRules)),
		printGreen("%v", protocolExistsAndCount("数据库")),
		printRed("%v", len(crack.MapCrackHost)),
		printRed("%v", len(poc.ListPocInfo)),
	)
}

func endHtml() {
	if len(iplist) == 0 {
		return
	}

	vulnerablehost, _ := calculateVulnerablePercentage(iplist, crack.MapCrackHost, poc.ListPocInfo)
	screenshotDir := global.SsaveIMGDIR
	couunt, _ := global.CountDirFiles(screenshotDir)

	// 获取截图图片路径
	var screenshots []string
	files, err := os.ReadDir(screenshotDir)
	if err == nil {
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".png") {
				screenshots = append(screenshots, filepath.ToSlash(filepath.Join(screenshotDir, f.Name())))
			}
		}
	}

	// 生成HTML报告
	html := generateHTMLReport(ReportData{
		Time:             time.Now().Format("2006-01-02 15:04:05"),
		TotalHosts:       len(iplist),
		VulnHosts:        vulnerablehost,
		LinuxCount:       linuxcount,
		WindowsCount:     windowscount,
		UnidentifiedOS:   len(iplist) - linuxcount - windowscount,
		PortsCount:       len(infolist),
		SSHCount:         protocolExistsAndCount("ssh"),
		RDPCount:         protocolExistsAndCount("rdp"),
		WebCount:         protocolExistsAndCount("WEB应用"),
		AppCount:         len(global.AppMatchedRules),
		DBCount:          protocolExistsAndCount("数据库"),
		ScreenshotCount:  couunt,
		ScreenshotDir:    screenshotDir,
		ScreenshotImages: screenshots,
		CrackList:        crack.MapCrackHost,
		PocList:          poc.ListPocInfo,
		PortServiceList:  infolist,
		IPList:           iplist,
		ChartJS:          template.JS(chartJS),
		ChartJSPlugin:    template.JS(chartJSPlugin),
	})

	filename := filepath.Join("ScanLog", time.Now().Format("200601021504")+"-report.html")
	if err := os.WriteFile(filename, []byte(html), 0644); err == nil {
		if runtime.GOOS == "windows" {
			_ = exec.Command("cmd", "/c", "start", filename).Run()
		}
	}
}

// calculateVulnerablePercentage 计算有漏洞的IP数量和百分比
func calculateVulnerablePercentage(iplist []string, crackHosts map[crack.HostPort]crack.SussCrack, pocInfos []poc.Flagcve) (int, float64) {
	uniqueIPs := make(map[string]struct{})

	// 正则表达式用于从URL中提取IP地址
	ipRegex := regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)

	// 从crack.MapCrackHost中收集所有Host
	for hostPort := range crackHosts {
		if ip := net.ParseIP(hostPort.Host); ip != nil {
			uniqueIPs[hostPort.Host] = struct{}{}
		}
	}

	// 从Flagcve的URL中提取IP地址
	for _, pocInfo := range pocInfos {
		if matches := ipRegex.FindStringSubmatch(pocInfo.Url); len(matches) > 0 {
			if ip := net.ParseIP(matches[0]); ip != nil {
				uniqueIPs[matches[0]] = struct{}{}
			}
		}
	}

	// 计算有漏洞的IP数量
	vulnerableCount := len(uniqueIPs)
	totalCount := len(iplist)

	// 计算百分比
	if totalCount == 0 {
		return vulnerableCount, 0.0
	}
	percentage := (float64(vulnerableCount) / float64(totalCount)) * 100.0
	return vulnerableCount, percentage
}
