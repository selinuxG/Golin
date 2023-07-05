package Protocol

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"strings"
)

// IsMySqlProtocol 基于gorm的登录错误消息判断是否为MySQL
func IsMySqlProtocol(host, port string) bool {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=1.5s", "root", "1234567", host, port, "mysql")
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 设置日志级别为 silent
	})
	if err != nil {
		if strings.Contains(err.Error(), "Access denied for user") || strings.Contains(err.Error(), "MySQL server") {
			return true
		}
		return false
	} else {
		return true
	}

}

// IsPgsqlProtocol 基于gorm的登录错误消息判断是否为MySQL
func IsPgsqlProtocol(host, port string) bool {
	log.SetOutput(io.Discard)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable  TimeZone=Asia/Shanghai connect_timeout=%d", host, "postgres", "123456", port, 3)
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 设置日志级别为 silent
	})
	if err != nil {
		if strings.Contains(err.Error(), "failed") || strings.Contains(err.Error(), "user") || strings.Contains(err.Error(), "postgres") {
			return true
		}
		return false
	} else {
		return true
	}

}
