package web

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CreateZipFromFiles 将文件路径切片中的文件压缩成一个名为saveName的zip包
func CreateZipFromFiles(filePaths []string, saveName string) error {
	// 创建保存的zip文件
	zipFile, err := os.Create(saveName)
	if err != nil {
		return fmt.Errorf("创建zip文件失败: %v", err)
	}
	defer zipFile.Close()

	// 创建zip.Writer，它将向zip文件写入内容
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历文件路径切片
	for _, filePath := range filePaths {
		if err = addFileToZip(filePath, zipWriter); err != nil {
			return fmt.Errorf("添加文件到zip失败: %v", err)
		}
	}
	return nil
}

// addFileToZip 向zip包中添加文件
func addFileToZip(filePath string, zipWriter *zip.Writer) error {
	// 打开待压缩文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}
