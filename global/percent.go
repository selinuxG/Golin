package global

import (
	"fmt"
	"github.com/fatih/color"
)

var (
	spinnerChars = []string{"|", "/", "-", "\\"} //进度条更新动画
	counter      = 0                             //当前已扫描的数量，计算百分比
)

// Percent 输出进度条
func Percent(succeCount, countall uint32) {
	PrintLock.Lock()
	defer PrintLock.Unlock()
	percent := (float64(succeCount) / float64(countall)) * 100.00
	spinChar := rotateSpinner()
	if percent == 100 {
		spinChar = "√"
	}

	percentStr := fmt.Sprintf("%.2f", percent) // 将百分比值格式化为字符串
	fmt.Print("\033[2K")                       // 擦除整行
	fmt.Printf("\r[%s] 当前进度: %s",
		spinChar,
		color.RedString("%s", fmt.Sprintf("%s%%", percentStr)),
	)
}

// 旋转进度条
func rotateSpinner() string {
	spinChar := spinnerChars[counter%len(spinnerChars)]
	counter++
	return spinChar
}
