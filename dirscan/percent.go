package dirscan

import (
	"fmt"
	"github.com/fatih/color"
)

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
