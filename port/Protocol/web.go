package Protocol

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"time"
)

func IsWeb(host, port string) string {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   3 * time.Second,
	}
	for _, v := range []string{"http", "https"} {
		url := ""
		htype := ""
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

		if (resp.StatusCode >= 200 && resp.StatusCode < 300) || resp.StatusCode == 401 || resp.StatusCode == 403 || resp.StatusCode == 404 {

			//查找title
			title := ""
			if resp.StatusCode == 200 {
				doc, err := goquery.NewDocumentFromReader(resp.Body)
				if err == nil {
					title = doc.Find("title").Text()
					if title != "" {
						title = fmt.Sprintf("Title:「%s」", title)
						title = strings.ReplaceAll(title, "\n", "")
					}
				}
			}
			serverType := resp.Header.Get("Server")
			if serverType != "" {
				serverType = fmt.Sprintf("server:「%s」", serverType)
			}
			return fmt.Sprintf("%s  %s  %s 状态码:「%d」", htype, serverType, title, resp.StatusCode)
		}

	}
	if port == "443" {
		return "https"
	}
	return ""
}
