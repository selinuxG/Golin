package run

import (
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
	fire := filepath.Join(pullpath, fmt.Sprintf("%s_%s.html", myname, myhost))
	os.Remove(fire)
	file, err := os.OpenFile(fire, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.FileMode(global.FilePer))
	if err != nil {
		errhost = append(errhost, myhost)
		return
	}
	html := redishtml()

	defer file.Close()

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
	users, err := client.Do(ctx, "ACL", "LIST").StringSlice()
	if err == nil {
		if len(users) > 0 {
			userechho := ""
			userechho += fmt.Sprintf("aclfile存储路径:%s<br>", aclfile)
			for _, user := range users {
				userechho += user + "<br>"
			}
			html = strings.ReplaceAll(html, "acluser信息详细信息", userechho)
		}
	}

	aclcount := client.ConfigGet(ctx, "acllog-max-len").Val()
	//log := client.Do(ctx, "ACL", "LOG").Val()
	html = strings.ReplaceAll(html, "基本信息详细信息", fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", ipaddr[1], pass[1], redistimout[1]))
	// 端口信息
	html = strings.ReplaceAll(html, "端口信息详细信息", fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", redisport[1], redissslport, protocols))
	// 日志信息
	html = strings.ReplaceAll(html, "日志信息详细信息", fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td></tr>", lofile[1], loglevel[1], aclcount))
	// 整体info信息
	infoStr := client.Info(ctx).Val()
	infoStr = strings.ReplaceAll(infoStr, "\r\n", "\n")
	infoStr = strings.ReplaceAll(infoStr, "\n", "<br>")
	html = strings.ReplaceAll(html, "info信息详细信息", fmt.Sprintf("<tr>%s</tr>", infoStr))
	html = strings.ReplaceAll(html, "替换名称", adr)

	os.WriteFile(fire, []byte(html), os.FileMode(global.FilePer))

}
