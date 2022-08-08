package checkpass

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	lesscount  []string
	lessnum    []string
	lessaz     []string
	lessAZ     []string
	lesssymbol []string
	firelog    = "采集完成目录//" + "checkpass.log"
	count      int
)

//密码强度必须为字⺟⼤⼩写+数字+符号，8位以上
func CheckPasswordLever(path string) {
	log.Printf("------即将启动%s模式下的密码强度检测功能", path)
	var pathtxt string
	switch path {
	case "linux":
		pathtxt = "ip.txt"
	case "mysql":
		pathtxt = "mysql.txt"
	case "redis":
		pathtxt = "redis.txt"
	case "postgresql":
		pathtxt = "postgresql"
	default:
		log.Println("必须配合 -run （linux、mysql、redis、postgresql下使用）")
	}

	fire, err := ioutil.ReadFile(pathtxt)
	if err != nil {
		log.Printf("读取%s文件失败!", pathtxt)
		//程序退出，状态码0表示成功,非0表示出错,程序会立刻终止，并且 defer 的函数不会被执行
		os.Exit(1)
	}
	succmylog(firelog)
	lines := strings.Split(string(fire), "\n")
	for i := 0; i < len(lines); i++ {
		firecount := strings.Count(lines[i], "~")
		if firecount != 4 {
			log.Println(lines[i], "格式错误,跳过")
			continue
		}
		count++
		a := lines[i]
		name := strings.Split(string(a), "~")[0]
		ps := strings.Split(string(a), "~")[3]
		if len(ps) < 8 {
			lesscount = append(lesscount, name)
		}
		num := `[0-9]{1}`
		a_z := `[a-z]{1}`
		A_Z := `[A-Z]{1}`
		symbol := `[!@#~$%^&*()+|_]{1}`
		if b, err := regexp.MatchString(num, ps); !b || err != nil {
			lessnum = append(lessnum, name)
		}
		if b, err := regexp.MatchString(a_z, ps); !b || err != nil {
			lessaz = append(lessaz, name)
		}
		if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
			lessAZ = append(lessAZ, name)
		}
		if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
			lesssymbol = append(lesssymbol, name)
		}
	}

	aa := fmt.Sprintf("不包含数字的主机为:%d个\n%s", len(lessnum), lessnum)
	bb := fmt.Sprintf("小于8位的主机为:%d个\n%s", len(lesscount), lesscount)
	cc := fmt.Sprintf("不包含大小写的主机为:%d个\n%s", len(lessAZ), lessAZ)
	dd := fmt.Sprintf("不包含小写的主机为:%d个\n%s", len(lessaz), lessaz)
	ff := fmt.Sprintf("不包含特殊符号的主机为:%d个\n%s", len(lesssymbol), lesssymbol)
	file, _ := os.OpenFile(firelog, os.O_WRONLY|os.O_APPEND, 0666)
	write := bufio.NewWriter(file)
	write.WriteString("------------" + Nowtime() + "------------\n")
	line := fmt.Sprintf("检查格式为:字⺟⼤⼩写+数字+特殊符号，8位以上!此次共检查主机:%d个\n", count)
	write.WriteString(line)
	write.WriteString(aa + "\n\n")
	write.WriteString(bb + "\n\n")
	write.WriteString(cc + "\n\n")
	write.WriteString(dd + "\n\n")
	write.WriteString(ff + "\n\n")
	write.Flush()
	defer file.Close()
}

func succmylog(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Create(path)
	}
}

func Nowtime() string {
	timeObj := time.Now()
	year := timeObj.Year()
	month := timeObj.Month()
	day := timeObj.Day()
	hour := timeObj.Hour()
	minute := timeObj.Minute()
	second := timeObj.Second()
	timenow := fmt.Sprintf("%d-%d-%d %d:%d:%d", year, month, day, hour, minute, second)
	return timenow
}
