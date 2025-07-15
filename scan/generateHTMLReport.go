package scan

import (
	"bytes"
	"embed"
	"golin/global"
	"golin/poc"
	"golin/scan/crack"
	"html/template"
	"regexp"
)

type ReportData struct {
	Time             string
	TotalHosts       int
	VulnHosts        int
	LinuxCount       int
	WindowsCount     int
	UnidentifiedOS   int
	PortsCount       int
	SSHCount         int
	RDPCount         int
	WebCount         int
	AppCount         int
	DBCount          int
	ScreenshotCount  int
	ScreenshotDir    string
	CrackList        map[crack.HostPort]crack.SussCrack
	PocList          []poc.Flagcve
	PortServiceList  []INFO
	IPList           []string
	ChartJS          template.JS
	ChartJSPlugin    template.JS
	ScreenshotImages []string
	Job              global.TaskJob
}

//go:embed template/*
var content embed.FS

//go:embed template/chat.js
var chartJS []byte

//go:embed template/chaplugin.js
var chartJSPlugin []byte

func generateHTMLReport(data ReportData) (string, error) {
	funcMap := template.FuncMap{
		"removeColor": func(input string) string {
			re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
			return re.ReplaceAllString(input, "")
		},
		"inc": func(i int) int {
			return i + 1
		},
	}

	tmpl, err := template.New("report.html").Funcs(funcMap).ParseFS(content, "template/report.html")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
