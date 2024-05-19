package global

import (
	"fmt"
	"os"
	"time"
)

// LogLevel 定义了不同的日志级别。
type LogLevel int

const (
	// LevelError 代表错误级别的日志。
	LevelError LogLevel = iota
	// LevelWarning 代表警告级别的日志。
	LevelWarning
	// LevelInfo 代表信息级别的日志。
	LevelInfo
)

// LogToFile 将日志消息写入指定的日志文件。如果文件不存在，会自动创建。
func LogToFile(level LogLevel, message string) {

	logFileName := "crack.log"
	// 根据日志级别添加不同的前缀
	levelPrefix := map[LogLevel]string{
		LevelError:   "[ERROR]",
		LevelWarning: "[WARNING]",
		LevelInfo:    "[INFO]",
	}
	prefix, ok := levelPrefix[level]
	if !ok {
		return
	}

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	// 构建日志消息，格式为: time [level] msg
	logMessage := fmt.Sprintf("%s %s %s\n", time.Now().Format(time.DateTime), prefix, message)

	// 写入日志文件
	_, _ = file.WriteString(logMessage)

}
