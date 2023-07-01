package global

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"sync"
)

var (
	spinnerChars = []string{"|", "/", "-", "\\"} //进度条更新动画
	counter      = 0                             //当前已扫描的数量，计算百分比
	greenColor   = "\033[32m"                    // 设置文本颜色为绿色
	resetColor   = "\033[0m"                     // 重置文本颜色
	redColor     = "\033[31m"                    // 设置文本颜色为红色
	colorOutput  = colorable.NewColorableStdout()
	//mu           sync.Mutex
)

// Percent 输出进度条
func Percent(mu *sync.Mutex, succeCount, countall int) {
	percent := (float64(succeCount) / float64(countall)) * 100.00

	// 根据百分比值选择相应颜色
	var colorCode string
	if percent < 100 {
		colorCode = redColor
	} else {
		colorCode = greenColor
	}
	spinChar := rotateSpinner(mu)
	if percent == 100 {
		spinChar = "√"
	}
	_, _ = fmt.Fprintf(colorOutput, "\r[%s] 当前进度: %s%.2f%%%s", spinChar, colorCode, percent, resetColor)
}

// 旋转进度条
func rotateSpinner(mu *sync.Mutex) string {
	mu.Lock()
	defer mu.Unlock()

	spinChar := spinnerChars[counter%len(spinnerChars)]
	counter++
	return spinChar
}
