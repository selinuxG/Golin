package web

import "embed"

//go:embed tpl/*
var f embed.FS

// IndexHtml 单主机页面
func IndexHtml() string {
	data, _ := f.ReadFile("tpl/index.html")
	return string(data)
}

// IndexFileHtml 多主机模式页面

func IndexFileHtml() string {
	data, _ := f.ReadFile("tpl/indexFile.html")
	return string(data)
}

func ErrorHtml() string {
	data, _ := f.ReadFile("tpl/error.html")
	return string(data)
}

// GolinHomeHtml 首页
func GolinHomeHtml() string {
	data, _ := f.ReadFile("tpl/golinHome.html")
	return string(data)
}

func GolinHistoryIndexHtml() string {
	data, _ := f.ReadFile("tpl/golinHistoryIndex.html")
	return string(data)
}

// DjHtml 模拟定级
func DjHtml() string {
	data, _ := f.ReadFile("tpl/dj.html")
	return string(data)
}

// DjLevelHtml 输出定级结果
func DjLevelHtml() string {
	data, _ := f.ReadFile("tpl/djLevel.html")
	return string(data)

}
