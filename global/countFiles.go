package global

import (
	"os"
	"path/filepath"
)

// CountDirFiles 接收一个目录返回目录中的文件数量
func CountDirFiles(dirPath string) (int, error) {
	var fileCount int
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileCount++
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return fileCount, nil
}
