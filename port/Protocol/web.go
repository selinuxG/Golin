package Protocol

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"golin/global"
	"golin/poc"
	"golin/port/crack"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type webinfo struct {
	url         string //web地址
	title       string //网站标题
	app         string //识别到的组件
	statuscode  int    //状态码
	ContentType string //ContentType
	xss         string //是否存在xss漏洞
	server      string //ContentType中的server
}

func IsWeb(host, port string, timeout int, xss, Poc bool) string {
	url := ""

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	for _, v := range []string{"https", "http"} {
		info := webinfo{}
		switch port {
		case "443":
			info.url = fmt.Sprintf("https://%s", host)
		case "80":
			info.url = fmt.Sprintf("http://%s", host)
		default:
			info.url = fmt.Sprintf("%s://%s:%s", v, host, port)
		}
		url = info.url

		resp, err := client.Get(info.url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if global.Debug {
			fmt.Println(string(body))
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err == nil {
			info.title = doc.Find("title").Text()
			info.title = strings.ReplaceAll(info.title, "\n", "")
			info.title = strings.ReplaceAll(info.title, "  ", "")

		}
		info.statuscode = resp.StatusCode
		info.ContentType = resp.Header.Get("Content-Type")
		info.app = CheckApp(string(body), resp.Header, resp.Cookies()) // 匹配组件

		if strings.Contains(strings.ToLower(info.app), "elasticsearch") {
			intport, _ := strconv.Atoi(port)
			crack.ListCrackHost = append(crack.ListCrackHost, crack.SussCrack{Host: host, Port: intport, Mode: "ElasticSearch"})
		}

		//xss扫描
		if xss {
			checkXSS, xssPayloads := CheckXss(url, body)
			if checkXSS {
				info.xss = xssPayloads
				poc.ListPocInfo = append(poc.ListPocInfo, poc.Flagcve{url, "XSS", xssPayloads})
			}
		}

		//poc扫描
		if Poc {
			go poc.CheckPoc(info.url, info.app)
		}

		// 基于title确认是否url是目录浏览
		if strings.Contains(strings.ToLower(info.title), "index of") {
			poc.ListPocInfo = append(poc.ListPocInfo, poc.Flagcve{Url: url, Cve: "目录浏览漏洞"})
		}

		info.server = resp.Header.Get("Server")

		return chekwebinfo(info)
	}

	if port == "443" {
		return url
	}

	return ""
}

// CheckApp 基于返回的body、headers、cookies判定组件信息
func CheckApp(body string, head map[string][]string, cookies []*http.Cookie) string {
	var app []string
	for _, rule := range RuleDatas {
		switch rule.Type {
		case "body":
			patterns, err := regexp.Compile(rule.Rule)
			if err == nil && patterns.MatchString(body) {
				app = append(app, rule.Name)
			}

		case "headers":
			for _, values := range head {
				for _, value := range values {
					patterns, err := regexp.Compile(`(?i)` + rule.Rule) //不区分大小写
					if err == nil && patterns.MatchString(value) {
						app = append(app, rule.Name)
					}
				}
			}

		case "cookie":
			for _, cookie := range cookies {
				patterns, err := regexp.Compile(`(?i)` + rule.Rule) //不区分大小写
				if err == nil && patterns.MatchString(cookie.Name) {
					app = append(app, rule.Name)
				}
			}
		}
	}

	return strings.Join(global.RemoveDuplicates(app), "、")

}

func chekwebinfo(info webinfo) string {
	output := fmt.Sprintf("%-23s ", info.url)

	if info.xss != "" {
		output += color.RedString("%s", fmt.Sprintf(" [XSS漏洞:%s]", info.xss))
	}

	if info.app != "" {
		output += color.GreenString("%s", fmt.Sprintf(" APP:「%s」", info.app))
	}
	if info.title != "" {
		info.title = strings.ReplaceAll(info.title, "  ", "")
		output += color.BlueString("%s", fmt.Sprintf(" title:「%s」", info.title))
	}
	if info.server != "" {
		output += fmt.Sprintf("%s", color.MagentaString("%s", fmt.Sprintf(" Server:「%s」", info.server)))
	}

	output += fmt.Sprintf(" Code:%d", info.statuscode)

	if info.ContentType != "" {
		output += fmt.Sprintf(" ContentType:%s", info.ContentType)
	}

	return output
}
