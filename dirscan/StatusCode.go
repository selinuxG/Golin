package dirscan

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	Protocol2 "golin/Protocol"
	"strconv"
	"strings"
	"sync"
)

var (
	succeCount = 0
	lock       sync.Mutex
)

type UrlStatus struct {
	Url         string   //url地址
	Code        int      //状态码
	Title       string   //标题
	ContentType string   //媒体类型
	Line        int      //行数
	App         string   //组件
	info        []string //漏洞
}

var ContentType = map[string]string{
	"application/zip":               "文件下载",
	"application/octet-stream":      "文件下载",
	"text/xml":                      "疑似敏感接口",
	"application/xml":               "疑似敏感接口",
	"application/json":              "疑似敏感接口",
	"multipart/form-data":           "文件上传",
	"video/avi":                     "文件下载",
	"audio/x-wav":                   "文件下载",
	"audio/x-ms-wma":                "文件下载",
	"audio/mp3":                     "文件下载",
	"video/mpeg4":                   "文件下载",
	"application/javascript":        "JavaScript脚本文件",
	"image/jpeg":                    "图片文件",
	"image/png":                     "图片文件",
	"application/pdf":               "PDF文件下载",
	"application/msword":            "Word文档下载",
	"application/vnd.ms-excel":      "Excel文档下载",
	"image/gif":                     "图片文件",
	"image/tiff":                    "图片文件",
	"video/mp4":                     "视频文件下载",
	"video/webm":                    "视频文件下载",
	"audio/ogg":                     "音频文件下载",
	"audio/mpeg":                    "音频文件下载",
	"application/sql":               "数据库文件",
	"application/rtf":               "RTF文件下载",
	"application/x-shockwave-flash": "Flash文件下载",
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
			Url:         URL,
			Code:        resp.StatusCode,
			Title:       title,
			ContentType: contype,
			Line:        line,
		}

		if strings.Contains(strings.ToLower(title), "index of") {
			yesurl.info = append(yesurl.info, "目录浏览漏洞")
		}

		for k, v := range ContentType {
			if strings.Contains(contype, k) {
				yesurl.info = append(yesurl.info, v)
			}
		}

		check, xss := Protocol2.CheckXss(URL, body)
		if check {
			yesurl.info = append(yesurl.info, fmt.Sprintf("xss:%s", xss))
		}
		yesurl.App = Protocol2.CheckApp(string(body), resp.Header, req.Cookies, resp.Header.Get("Server"), "", "") // 匹配组件

		_ = AppendUrlStatusToFile(yesurl) // 写入文件

		fmt.Printf(echoinfo(yesurl, URL))
	}
}

// echoinfo 输出详细信息
func echoinfo(yesurl UrlStatus, URL string) string {
	echo := ""
	if yesurl.App != "" {
		echo += color.GreenString("%s", fmt.Sprintf(" APP:「%s」", yesurl.App))
	}

	if len(yesurl.info) > 0 {
		echo += fmt.Sprintf("%s", color.RedString("%s ", fmt.Sprintf("[%s]", strings.Join(yesurl.info, ","))))
	}

	if yesurl.Title != "" {
		echo += fmt.Sprintf("%s", color.BlueString(" %s", fmt.Sprintf("title:「%s」 ", yesurl.Title)))
	}
	echo += fmt.Sprintf("%s", color.MagentaString("%s", fmt.Sprintf("ContentType:%s  ", yesurl.ContentType)))
	echo += fmt.Sprintf("Code:%d  Line:%d", yesurl.Code, yesurl.Line)

	return fmt.Sprintf("\r%-45s| %s \n", URL, echo)

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
