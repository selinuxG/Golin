package web

import "embed"

//go:embed template/*
var f embed.FS

// IndexHtml 单主机页面
func IndexHtml() string {
	data, _ := f.ReadFile("template/index.html")
	return string(data)
}

// IndexFileHtml 多主机模式页面

func IndexFileHtml() string {
	data, _ := f.ReadFile("template/indexFile.html")
	return string(data)
}

func ErrorHtml() string {
	data, _ := f.ReadFile("template/error.html")
	return string(data)
}

// GolinHomeHtml 首页
func GolinHomeHtml() string {
	data, _ := f.ReadFile("template/golinHome.html")
	return string(data)
}

func GolinHistoryIndexHtml() string {
	data, _ := f.ReadFile("template/golinHistoryIndex.html")
	return string(data)
}

// DjHtml 模拟定级
func DjHtml() string {
	data, _ := f.ReadFile("template/dj.html")
	return string(data)
}

// DjLevelHtml 输出定级结果
func DjLevelHtml() string {
	data, _ := f.ReadFile("template/djLevel.html")
	return string(data)
}
