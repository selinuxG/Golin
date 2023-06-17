package dirscan

import (
	"encoding/json"
	"os"
)

var filename = "dirScan.json"

// AppendUrlStatusToFile 写入成功url到dirScan.json文件
func AppendUrlStatusToFile(status UrlStatus) error {
	var statuses []UrlStatus
	data, err := os.ReadFile(filename)
	if err == nil {
		err = json.Unmarshal(data, &statuses)
		if err != nil {
			return err
		}
	}

	statuses = append(statuses, status)

	data, err = json.MarshalIndent(statuses, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, data, 0644)
	return err
}
