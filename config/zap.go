package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golin/global"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	Log *zap.Logger
)

func init() {
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		CallerKey:      "caller",
		EncodeLevel:    zapcore.CapitalLevelEncoder, //CapitalColorLevelEncoder  彩色输出
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		// 添加字段分隔符和缩进
		ConsoleSeparator: "    ",
	})
	// 创建控制台输出的核心对象，设置级别为 info
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleDebugging, zap.InfoLevel)

	// 创建文件输出的核心对象，设置级别为 warn
	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		CallerKey:      "caller",
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	})
	fileDebugging := zapcore.Lock(zapcore.AddSync(&lumberjack.Logger{
		Filename:   global.SuccessLog,
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     30, // days
	}))
	fileCore := zapcore.NewCore(fileEncoder, fileDebugging, zap.WarnLevel)
	// 合并两个核心对象
	core := zapcore.NewTee(consoleCore, fileCore)
	// 创建logger
	Log = zap.New(core, zap.AddCaller())

}
