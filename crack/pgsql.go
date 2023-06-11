package crack

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func pgsql(ctx context.Context, cancel context.CancelFunc, ip, user, passwd string, port int) {
	defer func() {
		wg.Done()
		<-ch
	}()
	select {
	case <-ctx.Done():
		return
	default:
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%d sslmode=disable  TimeZone=Asia/Shanghai connect_timeout=3", ip, user, passwd, port)
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err == nil {
		end(ip, user, passwd, port)
		cancel()
	}
}
