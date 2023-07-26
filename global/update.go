package global

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type releaseInfo struct {
	TagName string        `json:"tag_name"`
	Assets  []BrowserDown `json:"assets"`
}
type BrowserDown struct {
	BrowserDownloadUrl string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// CheckForUpdate 检查更新
func CheckForUpdate() (releaseInfo, error) {
	var info releaseInfo
	response, err := http.Get(RepoUrl)
	if err != nil {
		return info, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return info, fmt.Errorf("failed to fetch latest release information with status code %d", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("无法读取响应正文:", err)
		return info, err
	}

	err = json.Unmarshal(body, &info)
	if err != nil {
		return info, err
	}

	return info, nil
}
