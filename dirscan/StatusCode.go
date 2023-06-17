package dirscan

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strconv"
	"sync"
	"time"
)

var (
	succeCount = 0
	lock       sync.Mutex
)

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
	t1 := time.Now() // 记录发送请求前的时间
	resp, _, errs := request.Get(URL).End()
	t2 := time.Now() // 记录收到响应后的时间

	if len(errs) > 0 || resp == nil {
		//fmt.Println(errs)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	//fmt.Println(URL, resp.StatusCode, code)
	if statusCodeInRange(resp.StatusCode, code) {
		fmt.Print("\033[2K") // 擦除整行
		//查找title
		title := ""
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err == nil {
			title = doc.Find("title").Text()
		}
		fmt.Printf("\r[√] %s 发现url：%s 状态码:%d Title:%s 响应时间:%v\n", time.Now().Format(time.DateTime), URL, resp.StatusCode, title, t2.Sub(t1))
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
