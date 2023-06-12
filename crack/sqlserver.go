package crack

import (
	"context"
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func sqlservercon(ctx context.Context, cancel context.CancelFunc, ip, user, passwd string, port, timeout int) {
	defer func() {
		wg.Done()
		<-ch
	}()
	select {
	case <-ctx.Done():
		return
	default:
	}

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=master&timeout=%ds", user, passwd, ip, port, timeout)
	_, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err == nil {
		end(ip, user, passwd, port)
		cancel()
	}
}
