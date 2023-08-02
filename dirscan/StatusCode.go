package dirscan

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"strconv"
	"strings"
	"sync"
)

var (
	succeCount = 0
	lock       sync.Mutex
)

type UrlStatus struct {
	Url           string //url地址
	Code          int    //状态码
	Title         string //标题
	ContentLength string //大小
	ContentType   string //媒体类型
	Line          int    //行数
}

func isStatusCodeOk(URL string) {
	defer func() {
		wg.Done()
		<-ch
		doSomething()
		// succeCount += 1 // 每次结束数量加一
		percent()
	}()
	if request == nil {
		return
	}

	req := request.Clone().Get(URL) // 使用 Clone() 方法创建请求对象
	if req == nil {
		return
	}

	resp, body, errs := req.EndBytes() // 使用 EndBytes() 方法获取响应和响应体
	if len(errs) > 0 || resp == nil {
		return
	}

	if statusCodeInRange(resp.StatusCode, code) {
		fmt.Print("\033[2K") // 擦除整行
		// 查找 title
		title := ""
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body)) // 使用响应体创建文档对象
		if err == nil {
			title = doc.Find("title").Text()
		}

		line := bytes.Count(body, []byte("\n")) // 行数
		contype := resp.Header.Get("Content-Type")
		contype = strings.Split(contype, ";")[0]

		yesurl := UrlStatus{
			Url:           URL,
			Code:          resp.StatusCode,
			Title:         title,
			ContentLength: FormatBytes(resp.ContentLength),
			ContentType:   contype,
			Line:          line,
		}
		info := ""
		if strings.Contains(strings.ToLower(title), "index of") {
			info += " 目录浏览漏洞"
		}
		if contype == "application/zip" {
			info += " 文件下载"
		}
		if contype == "text/xml" {
			info += " 疑似敏感信息泄露"
		}
		_ = AppendUrlStatusToFile(yesurl) // 写入文件
		fmt.Printf("\r%s %-40s | code:%s | Title:%-30s | Length:%s | Type:%-5s | Line:%-5s |%s\n",
			color.GreenString("%s", "[√]"),
			yesurl.Url,
			color.GreenString("%d", yesurl.Code),
			color.GreenString("%s", yesurl.Title),
			color.GreenString("%s", yesurl.ContentLength),
			color.GreenString("%s", yesurl.ContentType),
			color.GreenString("%d", yesurl.Line),
			color.RedString("%s", info),
		)

		return
	}
}

// statusCodeInRange 确认切片状态是否在搜索队列中
func statusCodeInRange(status int, rangeSlice []string) bool {
	strcode := strconv.Itoa(status)
	for _, v := range rangeSlice {
		if v == strcode {
			return true
		}
	}
	return false
}

// 加锁
func doSomething() {
	lock.Lock()
	defer lock.Unlock()
	succeCount++
}

// FormatBytes 大小转换
func FormatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
	)

	switch {
	case bytes < KB:
		return fmt.Sprintf("%d B", bytes)
	case bytes < MB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	}
}
