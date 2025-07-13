package global

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

func HashMD5(url string) (string, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //跳过证书的验证
		},
		DisableKeepAlives: true, //禁用HTTP连接的keep-alive 特性
	}

	client := &http.Client{
		Transport: transport,
	}
	req, err := http.NewRequest("GET", url+"/favicon.ico", nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // 确保响应体在函数结束时关闭

	// 创建一个新的 MD5 哈希对象
	hash := md5.New()
	_, err = io.Copy(hash, resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	md5Hex := hex.EncodeToString(hash.Sum(nil))
	return md5Hex, nil
}
