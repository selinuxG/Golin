package dirscan

import (
	"fmt"
	"github.com/mattn/go-colorable"
)

// percent 输出进度条
func percent() {
	percent := (float64(succeCount) / float64(countall)) * 100.00
	colorOutput := colorable.NewColorableStdout()
	// 设置文本颜色为红色
	redColor := "\033[31m"
	// 设置文本颜色为绿色
	greenColor := "\033[32m"
	// 重置文本颜色
	resetColor := "\033[0m"
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
