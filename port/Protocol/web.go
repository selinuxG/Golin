package Protocol

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golin/global"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type webinfo struct {
	url           string
	title         string
	servertype    string
	app           string
	statuscode    int
	ContentLength int64
	ContentType   string
}

func IsWeb(host, port string, timeout int) string {
	url := ""

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	for _, v := range []string{"http", "https"} {
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
		info.ContentLength = resp.ContentLength
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err == nil {
			info.title = doc.Find("title").Text()
			info.title = strings.ReplaceAll(info.title, "\n", "")
			info.title = strings.ReplaceAll(info.title, "  ", "")

		}
		info.servertype = resp.Header.Get("Server")
		info.statuscode = resp.StatusCode
		info.ContentLength = resp.ContentLength
		info.ContentType = resp.Header.Get("Content-Type")
		info.app = CheckApp(string(body), resp.Header, resp.Cookies()) // 匹配组件

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
	output := fmt.Sprintf("%s ", info.url)

	if info.app != "" {
		output += fmt.Sprintf(" APP:「%s」", info.app)
	}
	if info.title != "" {
		output += fmt.Sprintf(" title:「%s」", info.title)
	}
	if info.servertype != "" {
		output += fmt.Sprintf(" server:「%s」", info.servertype)
	}

	output += fmt.Sprintf(" Code:「%d」 Length:「%d」 ContentType:「%s」", info.statuscode, info.ContentLength, info.ContentType)

	return output
}
