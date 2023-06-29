package dirscan

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
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
		//succeCount += 1 //每次结束数量加一
		percent()
	}()
	if request == nil {
		return
	}
	//resp, _, errs := request.Get(URL).End()
	req := request.Get(URL)
	if req == nil {
		return
	}
	resp, body, errs := req.End()
	if len(errs) > 0 || resp == nil {
		return
	}

	if statusCodeInRange(resp.StatusCode, code) {
		fmt.Print("\033[2K") // 擦除整行
		//查找title
		title := ""
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err == nil {
			title = doc.Find("title").Text()
		}

		line := strings.Count(body, "\n") //行数
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
		_ = AppendUrlStatusToFile(yesurl) //写入文件
		fmt.Printf("\r[√] Url:「%s」 State:「%d」 Title:「%s」 Length:「%s」 Type: 「%s」 Line:「%d」 \n", yesurl.Url, yesurl.Code, yesurl.Title, yesurl.ContentLength, yesurl.ContentType, yesurl.Line)
		return
	}
	return
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
