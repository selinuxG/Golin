package crack

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

func pgsql(ctx context.Context, cancel context.CancelFunc, ip, user, passwd string, port, timeout int, ch <-chan struct{}, wg *sync.WaitGroup) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%d sslmode=disable  TimeZone=Asia/Shanghai connect_timeout=%d", ip, user, passwd, port, timeout)
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err == nil {
		end(ip, user, passwd, port, "PgSQL")
		cancel()
	}
}
