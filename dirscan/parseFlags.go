package dirscan

import (
	"crypto/tls"
	"embed"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
	"golin/global"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	spinnerChars = []string{"|", "/", "-", "\\"} //进度条更新动画
	code         []string                        //扫描的状态
	mu           sync.Mutex
	counter      = 0 //当前已扫描的数量，计算百分比
	ch           = make(chan struct{}, 30)
	wg           = sync.WaitGroup{}
	proxyurl     = ""     //代理地址
	file         = ""     //读取的字典文件
	checkurl     []string //扫描的url切片
	countall     = 0
	request      *gorequest.SuperAgent
)

//go:embed url.txt
var urlData embed.FS

func ParseFlags(cmd *cobra.Command, args []string) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	url, _ := cmd.Flags().GetString("url")
	url = strings.TrimSuffix(url, "/") // 删除最后的/

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		fmt.Printf("[-] URL需指定前缀！http or https！\n")
		return
	}
	chcount, _ := cmd.Flags().GetInt("chan")
	timeout, _ := cmd.Flags().GetInt("timeout")
	ch = make(chan struct{}, chcount)

	proxyurl, _ = cmd.Flags().GetString("proxy") //不为空则设置代理
	if proxyurl != "" {
		request = gorequest.New().
			Proxy(proxyurl).
			Timeout(time.Duration(timeout) * time.Second).
			TLSClientConfig(tlsConfig)
	} else {
		request = gorequest.New().
			Timeout(time.Duration(timeout) * time.Second).
			TLSClientConfig(tlsConfig)
	}
	Agent, _ := cmd.Flags().GetString("Agent") //如果Agent不为空则自定义User-Agent
	if Agent != "" {
		request.Set("User-Agent", Agent)
	}

	file, _ = cmd.Flags().GetString("file") //读取字典文件
	if !global.PathExists(file) {
		fmt.Printf("[-] 使用内置字典进行扫描,可通过-f指定扫描字典！\n")
		data, _ := urlData.ReadFile("url.txt")
		datastr := strings.ReplaceAll(string(data), "\r\n", "\n")
		for _, u := range strings.Split(datastr, "\n") {
			if len(u) == 0 {
				continue
			}
			checkurl = append(checkurl, u)
		}
	} else {
		data, _ := os.ReadFile(file)
		str := strings.ReplaceAll(string(data), "\r\n", "\n")
		for _, u := range strings.Split(str, "\n") {
			if len(u) == 0 {
				continue
			}
			checkurl = append(checkurl, u)
		}
	}

	countall = len(removeDuplicates(checkurl)) //去重

	waittime, _ := cmd.Flags().GetInt("wait")      //循环超时
	codese, _ := cmd.Flags().GetString("code")     //搜索的状态码
	for _, s := range strings.Split(codese, ",") { //根据分隔符写入到状态切片
		code = append(code, s)
	}
	codestr := strings.Join(code, ", ")

	fmt.Printf("[*] 开始运行dirsearch模式 共计尝试:%d次 超时等待:%d/s 循环等待:%d/s 并发数:%d 寻找状态码:%s 代理地址:%s\n ", countall, timeout, waittime, chcount, codestr, proxyurl)
	for _, checku := range removeDuplicates(checkurl) {
		if checku[0] != '/' { //判断第一个字符是不是/不是的话则增加
			checku = "/" + checku
		}
		ch <- struct{}{}
		wg.Add(1)
		go isStatusCodeOk(fmt.Sprintf("%s%s", url, checku)) //传递完整url\

		if waittime > 0 { //延迟等待
			time.Sleep(time.Duration(waittime) * time.Second)
		}
	}
	wg.Wait()
	time.Sleep(time.Second * 1) //等待1秒是因为并发问题，等待进度条。
	//percent()
	fmt.Print("\033[2K") // 擦除整行
	fmt.Printf("\r")

}

// removeDuplicates 切片去重
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
