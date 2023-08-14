package Protocol

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"golin/global"
	"golin/poc"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type WebInfo struct {
	url         string
	title       string
	statuscode  int
	ContentType string
	app         string
	server      string
}

func IsWeb(host, port string, timeout int, Poc bool) map[string]string {
	results := make(map[string]string)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	for _, v := range []string{"https", "http"} {
		info := WebInfo{}
		switch port {
		case "443":
			info.url = fmt.Sprintf("https://%s", host)
		case "80":
			info.url = fmt.Sprintf("http://%s", host)
		default:
			info.url = fmt.Sprintf("%s://%s:%s", v, host, port)
		}
		body, err := handleRequest(client, &info)
		if err != nil {
			continue
		}

		handlePocAndXss(&info, body, Poc)
		results[v] = chekwebinfo(info)
		return results
	}
	return results
}

// handleRequest 请求网页并补充WebInfo结构体
func handleRequest(client *http.Client, info *WebInfo) ([]byte, error) {
	resp, err := client.Get(info.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body)) //获取标题
	if err == nil {
		info.title = strings.TrimSpace(doc.Find("title").Text())
	}

	info.statuscode = resp.StatusCode
	info.ContentType = resp.Header.Get("Content-Type")
	info.app = CheckApp(string(body), resp.Header, resp.Cookies()) // 匹配组件
	info.server = resp.Header.Get("Server")

	if global.Debug {
		fmt.Println(string(body))
		fmt.Println(resp.Header)
	}

	return body, nil
}

// handlePocAndXss 漏洞扫描以及POC扫描
func handlePocAndXss(info *WebInfo, body []byte, Poc bool) {
	if !Poc {
		return
	}
	poc.CheckPoc(info.url, info.app) //POC扫描

	checkXSS, xssPayloads := CheckXss(info.url, body) //XSS扫描
	if checkXSS {
		poc.ListPocInfo = append(poc.ListPocInfo, poc.Flagcve{Url: info.url, Cve: "XSS", Flag: xssPayloads})
	}

	// 基于title确认是否url是目录浏览
	if strings.Contains(strings.ToLower(info.title), "index of") {
		poc.ListPocInfo = append(poc.ListPocInfo, poc.Flagcve{Url: info.url, Cve: "目录浏览漏洞"})
	}
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
					//fmt.Println(patterns, err)
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

func chekwebinfo(info WebInfo) string {
	output := fmt.Sprintf("%-23s ", info.url)

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
