package crack

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func mySql(ctx context.Context, cancel context.CancelFunc, ip, user, passwd string, port int) {
	defer func() {
		wg.Done()
		<-ch
	}()
	select {
	case <-ctx.Done():
		return
	default:
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=1.5s", user, passwd, ip, port, "mysql")
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 设置日志级别为 silent
	})

	if err == nil {
		end(ip, user, passwd, port)
		cancel()
	}
}
