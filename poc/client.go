package poc

import (
	"crypto/tls"
	"net/http"
	"time"
)

func newRequest(req *http.Request, timeout time.Duration) (*http.Response, float64, error) {

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}

	if timeout == 0 {
		timeout = 3 * time.Second
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * timeout,
	}

	start := time.Now() // 记录开始时间
	resp, err := client.Do(req)
	elapsed := time.Since(start) // 计算耗时

	return resp, elapsed.Seconds(), err
}
