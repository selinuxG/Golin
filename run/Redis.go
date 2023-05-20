package run

import (
	"bufio"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
	"golin/global"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Redis(cmd *cobra.Command, args []string) {
	//获取分隔符，默认是||
	spr, err := cmd.Flags().GetString("spript")
	if err != nil {
		fmt.Println(err)
		return
	}
	//如果value值不为空则是运行一次的模式
	value, err := cmd.Flags().GetString("value")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(value) > 10 {
		Onlyonerun(value, spr, "Redis")
		wg.Wait()
		zlog.Info("单次运行Redis模式结束！")
		return
	}
	//下面是多线程的模式
	ippath, err := cmd.Flags().GetString("ip")
	if err != nil {
		fmt.Println(err)
		return
	}
	//判断redis.txt文件是否存在
	Checkfile(ippath, fmt.Sprintf("名称%sip%s用户%s密码%s端口", Split, Split, Split, Split), global.FilePer, ippath)

	// 运行share文件中的函数
	Rangefile(ippath, spr, "Redis")
	wg.Wait()
	//完成前最后写入文件
	Deffile("Redis", count, count-len(errhost), errhost)

}

func Runredis(myname, myuser, myhost, mypasswd, myport1 string) {
	defer wg.Done()
	Port := strings.Replace(myport1, "\r", "", -1)
	ctx := context.Background()
	adr := myhost + ":" + Port
	//如果为user=null则改为空密码，web功能会用到此功能
	if myuser == "null" {
		myuser = ""
	}
	client := redis.NewClient(&redis.Options{
		Addr:            adr,
		Username:        myuser,
		Password:        mypasswd,
		DB:              0,
		DialTimeout:     1 * time.Second,
		MinRetryBackoff: 1 * time.Second,
		ReadTimeout:     1 * time.Second,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		errhost = append(errhost, myhost)
		return
	}

	pullpath := filepath.Join(succpath, "Redis")
	_, err = os.Stat(pullpath)
	if os.IsNotExist(err) {
		os.MkdirAll(pullpath, os.FileMode(global.FilePer))
	}
	fire := filepath.Join(pullpath, fmt.Sprintf("%s_%s.log", myname, myhost))
	os.Remove(fire)
	file, err := os.OpenFile(fire, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.FileMode(global.FilePer))
	if err != nil {
		errhost = append(errhost, myhost)
		return
	}
	defer file.Close()
	write := bufio.NewWriter(file)

	client.Get(ctx, "config").Val()
	ipaddr := client.ConfigGet(ctx, "bind").Val()
	lofile := client.ConfigGet(ctx, "logfile").Val()
	loglevel := client.ConfigGet(ctx, "loglevel").Val()
	pass := client.ConfigGet(ctx, "requirepass").Val()
	redistimout := client.ConfigGet(ctx, "timeout").Val()
	redisport := client.ConfigGet(ctx, "port").Val()
	redissslport := client.ConfigGet(ctx, "tls-port").Val()
	aclfile := client.ConfigGet(ctx, "aclfile").Val()
	protocols := client.ConfigGet(ctx, "tls-protocols").Val()

	//执行自定义命令,获取用户信息
	write.WriteString("----------------------------用户信息\n")

	users, err := client.Do(ctx, "ACL", "LIST").StringSlice()
	if err == nil {
		if len(users) > 0 {
			write.WriteString(fmt.Sprintf("密码信息存储路径：%s\n", aclfile))
			for _, user := range users {
				write.WriteString(fmt.Sprintf("%s\n", user))
			}
		}
	}

	write.WriteString("----------------------------基本信息\n")
	write.WriteString(fmt.Sprintf("监听网卡地址为:%s\n", ipaddr[1]))
	write.WriteString(fmt.Sprintf("日志存储为:%s  日志等级为:%s\n", lofile[1], loglevel[1]))
	write.WriteString(fmt.Sprintf("密码信息为:%s\n", pass[1]))
	write.WriteString(fmt.Sprintf("超时时间为:%s\n", redistimout[1]))
	write.WriteString(fmt.Sprintf("redis非加密运行端口为:%s\n", redisport[1]))
	write.WriteString(fmt.Sprintf("redis加密运行端口为:%s\n", redissslport[1]))
	write.WriteString(fmt.Sprintf("redis加密协议:%s (tls-protocols为设置服务端支持的TLS协议版本，如果为空则为默认支持TLSv1.2以及TLSv1.3)\n", protocols))
	write.WriteString("----------------------------acl日志信息\n")
	aclcount := client.ConfigGet(ctx, "acllog-max-len").Val()
	write.WriteString(fmt.Sprintf("存储acl日志条数最大为：%s\n", aclcount[1]))
	log := client.Do(ctx, "ACL", "LOG").Val()
	write.WriteString(fmt.Sprintf("%v", log))

	write.WriteString("\n----------------------------info信息\n")
	write.WriteString(client.Info(ctx).Val())
	write.Flush()
}
