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

var title = ""

func IsWeb(host, port string, timeout int) string {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	for _, v := range []string{"http", "https"} {
		url, htype := "", ""
		switch port {
		case "443":
			url = fmt.Sprintf("https://%s", host)
			htype = "https"

		case "80":
			url = fmt.Sprintf("http://%s", host)
			htype = "http"

		default:
			url = fmt.Sprintf("%s://%s:%s", v, host, port)
			htype = v
		}

		resp, err := client.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err == nil {
			title = doc.Find("title").Text()
			if title != "" {
				title = strings.ReplaceAll(title, "\n", "")
				title = strings.ReplaceAll(title, " ", "")
				title = fmt.Sprintf("Title:「%s」", title)
			}
		}

		serverType := resp.Header.Get("Server")
		if serverType != "" {
			serverType = fmt.Sprintf("server:「%s」", serverType)
		}

		// 匹配组件
		app, checkapp := "", CheckApp(string(body), resp.Header, resp.Cookies())
		if checkapp != "" {
			app = fmt.Sprintf("APP:「%s」", checkapp)
		}
		return fmt.Sprintf("%-3s | %-3d | %s %s %s",
			htype,
			resp.StatusCode,
			app,
			serverType,
			title,
		)
	}

	if port == "443" {
		return "https"
	}
	return ""
}

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

	return strings.Join(global.RemoveDuplicates(app), ",")

}
