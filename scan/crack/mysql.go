package crack

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func mySql(cancel context.CancelFunc, ip, user, passwd string, port, timeout int) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%ds", user, passwd, ip, port, "mysql", timeout)
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Discard, // 设置为丢弃所有日志
	})
	if err == nil {
		end(ip, user, passwd, port, "MySql")
		cancel()
	}
}
