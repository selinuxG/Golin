package poc

import (
	"net/http"
	"time"
)

func newRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 3,
	}
	resp, err := client.Do(req)
	return resp, err
}
