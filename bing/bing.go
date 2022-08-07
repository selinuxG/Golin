package bing

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

// 获取bing官网壁纸，调用官网API保存到本地。
func Getbing() {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	url := "https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1"
	resp, err := client.Get(url)

	if err != nil {
		fmt.Println("确保能联网？")
		return

	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("确保能联网？")
			return
		}
		//正则匹配地址
		jsonBlob := string(body)
		find := (`"url":"(.*?)","urlbase`)
		re := regexp.MustCompile(find)
		matches := re.FindStringSubmatch(jsonBlob)
		for i, m := range matches {
			// 过滤掉第一个元素
			if i > 0 {
				geturl := fmt.Sprintf("https://www.bing.com/%s", m)
				resp, err := http.Get(geturl)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				// 创建一个文件用于保存
				path := "bing壁纸"
				_, err = os.Stat(path)
				if os.IsNotExist(err) {
					os.Mkdir(path, 0666)
				}
				imgname := fmt.Sprintf("bing壁纸/%s.jpg", nowtime())
				out, err := os.Create(imgname)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer out.Close()

				// 然后将响应流和文件流对接起来
				_, err = io.Copy(out, resp.Body)
				if err != nil {
					panic(err)
				}
				fmt.Println("今日:", nowtime(), ",下载bing官网壁纸完成")

			}
		}
	}
}

func nowtime() string {
	timeObj := time.Now()
	year := timeObj.Year()
	month := timeObj.Month()
	day := timeObj.Day()
	// minute := timeObj.Minute()
	// second := timeObj.Second()
	timenow := fmt.Sprintf("%d-%d-%d", year, month, day)
	return timenow
}
