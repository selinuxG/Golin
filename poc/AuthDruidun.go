package poc

import (
	"io"
	"net/http"
	"strings"
)

func AuthDruidun(url string) {
	url += "/druid/index.html"
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := newRequest(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyBytes, err2 := io.ReadAll(resp.Body)
		if err2 != nil {
			return
		}
		bodyString := string(bodyBytes)
		if strings.Contains(bodyString, "Druid Stat Index") &&
			strings.Contains(bodyString, "DruidVersion") &&
			strings.Contains(bodyString, "DruidDrivers") {
			flags := Flagcve{
				url: url,
				cve: "Druid未授权访问",
			}
			echoFlag(flags)
		}
	}
}
