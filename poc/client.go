package poc

import (
	"crypto/tls"
	"net/http"
	"time"
)

func newRequest(req *http.Request, timeout time.Duration) (*http.Response, error) {

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

	resp, err := client.Do(req)
	return resp, err
}
