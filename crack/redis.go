package crack

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

func rediscon(ip, user, passwd string, port int) {
	defer func() {
		wg.Done()
		<-ch
	}()
	client := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", ip, port),
		Username:        user,
		Password:        passwd,
		DB:              0,
		DialTimeout:     1 * time.Second,
		MinRetryBackoff: 1 * time.Second,
		ReadTimeout:     1 * time.Second,
	})
	_, err := client.Ping(ctx).Result()
	if err == nil {
		end(ip, user, passwd, port)
	}
}
