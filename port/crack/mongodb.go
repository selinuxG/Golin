package crack

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func mongodbcon(ctx context.Context, cancel context.CancelFunc, ip, user, passwd string, port, timeout int) {
	defer func() {
		wg.Done()
		<-ch
	}()
	select {
	case <-ctx.Done():
		return
	default:
	}
	ctxx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	// 建立连接
	connOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", ip, port))
	if user != "" && passwd != "" {
		connOptions.SetAuth(options.Credential{
			Username: user,
			Password: passwd,
		})
	}
	conn, err := mongo.Connect(ctxx, connOptions)
	if err != nil {
		return
	}
	// 测试连接
	err = conn.Ping(ctx, readpref.Primary())
	if err == nil {
		end(ip, user, passwd, port, "MongoDB")
		cancel()
	}
	return
}
