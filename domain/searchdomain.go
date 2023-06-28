package domain

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"net"
	"sync"
)

var (
	succeCount  = 0
	lock        sync.Mutex
	colorOutput = colorable.NewColorableStdout()
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
		//ip := strings.Join(host, " ")
		fmt.Printf("\r[√] %s         \n ", subdomain)
	}

}

// percent 输出进度条
func percent() {
	percent := (float64(succeCount) / float64(countall)) * 100.00

	// 根据百分比值选择相应颜色
	var colorCode string
	if percent < 100 {
		colorCode = redColor
	} else {
		colorCode = greenColor
	}
	spinChar := rotateSpinner()
	if percent == 100 {
		spinChar = "√"
	}
	fmt.Fprintf(colorOutput, "\r[%s] 当前进度: %s%.2f%%%s", spinChar, colorCode, percent, resetColor)
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
