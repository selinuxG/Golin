package crack

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

func rediscon(cancel context.CancelFunc, ip, user, passwd string, port, timeout int) {
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", ip, port),
		Username:        user,
		Password:        passwd,
		DB:              0,
		DialTimeout:     time.Duration(timeout) * time.Second,
		MinRetryBackoff: time.Duration(timeout) * time.Second,
		ReadTimeout:     time.Duration(timeout) * time.Second,
	})
	_, err := client.Ping(ctx).Result()
	if err == nil {
		end(ip, user, passwd, port, "Redis")
		cancel()
	}
}
