package poc

import (
	"crypto/tls"
	"net/http"
	"time"
)

func newRequest(req *http.Request) (*http.Response, error) {

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 3,
	}

	resp, err := client.Do(req)
	return resp, err
}
