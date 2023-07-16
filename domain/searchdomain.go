package domain

import (
	"fmt"
	"github.com/fatih/color"
	"net"
	"sync"
)

var (
	succeCount = 0
	lock       sync.Mutex
)

func searchDomain(subdomain string) {
	defer func() {
		wg.Done()
		<-ch
		doSomething() //计数
		percent()     //输出进度条
	}()
	_, err := net.LookupHost(subdomain)
	if err == nil {
		fmt.Printf("\r[√] %s         \n ", subdomain)
	}

}

// percent 输出进度条
func percent() {
	percent := (float64(succeCount) / float64(countall)) * 100.00
	spinChar := rotateSpinner()
	if percent == 100 {
		spinChar = "√"
	}
	percentStr := fmt.Sprintf("%.2f", percent) // 将百分比值格式化为字符串
	fmt.Printf("\r[%s] 当前进度: %s",
		spinChar,
		color.RedString("%s", fmt.Sprintf("%s%%", percentStr)),
	)
}

// 旋转进度条
func rotateSpinner() string {
	mu.Lock()
	defer mu.Unlock()

	spinChar := spinnerChars[counter%len(spinnerChars)]
	counter++
	return spinChar
}

// 加锁
func doSomething() {
	lock.Lock()
	defer lock.Unlock()
	succeCount++
}
