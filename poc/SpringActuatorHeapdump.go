package poc

import "net/http"

func SpringActuatorHeapdump(url string) {
	url += "/actuator/heapdump"
	req, _ := http.NewRequest("HEAD", url, nil)

	resp, err := newRequest(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 && resp.Header.Get("Content-Type") == "application/octet-stream" {
		flags := Flagcve{
			url:  url,
			cve:  "Spring Actuator Heapdump",
			flag: "Spring内存文件下载",
		}
		echoFlag(flags)
	}

}
