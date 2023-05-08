package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"golin/global"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "检查更新程序",
	Long:  `通过api.github.com进行检查更新程序`,
	Run: func(cmd *cobra.Command, args []string) {
		newrelease, err := checkForUpdate()
		if err != nil {
			fmt.Println("更新失败:", err)
			return
		}
		if global.Version == newrelease.TagName {
			fmt.Println("当前版本已为最新版本,无需更新。")
			return
		}
		fmt.Println(fmt.Sprintf("当前版本为:%s,最新版本为:%s,可更新..", global.Version, newrelease.TagName))
		var input string
		fmt.Print("是否更新？y 或 n: ")
		for {
			_, err := fmt.Scanln(&input)
			if err != nil {
				fmt.Println("发生错误:", err)
				continue
			}
			input = strings.ToLower(input) //所有大写转换为小写
			switch input {
			case "y", "yes":
				proxy, _ := cmd.Flags().GetString("proxy")
				err := downloadFile(newrelease.Assets[0].BrowserDownloadUrl, "golin.exe", proxy)
				if err != nil {
					fmt.Printf("更新失败%s\n", err)
					return
				}
				fmt.Println("更新完成")
				os.Exit(0)
			case "n", "no":
				fmt.Println("已取消更新...")
				os.Exit(0)
			default:
				fmt.Printf("输入无效,请输入y/n:")
				continue
			}
		}
	}}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("proxy", "p", "", "此参数是指定代理ip(仅允许http/https代理哦)")
}

type releaseInfo struct {
	TagName string        `json:"tag_name"`
	Assets  []BrowserDown `json:"assets"`
}
type BrowserDown struct {
	BrowserDownloadUrl string `json:"browser_download_url"`
}

// checkForUpdate 检查更新
func checkForUpdate() (releaseInfo, error) {
	var info releaseInfo
	response, err := http.Get(global.RepoUrl)
	if err != nil {
		return info, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return info, fmt.Errorf("failed to fetch latest release information with status code %d", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("无法读取响应正文:", err)
		return info, err
	}

	err = json.Unmarshal(body, &info)
	if err != nil {
		return info, err
	}

	return info, nil
}

// downloadFile 下载更新
func downloadFile(downurl, localPath, proxy string) error {
	fmt.Println("开始更新....")
	client := http.Client{}
	if proxy != "" {
		urli := url.URL{}
		urlproxy, err := urli.Parse(proxy)
		if err != nil {
			fmt.Println("无法连接代理IP!", proxy)
			return nil
		}
		client = http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(urlproxy),
			},
		}
	}
	resp, err := client.Get(downurl)
	if err != nil {
		return fmt.Errorf("无法下载文件：%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码：%d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// 获取当前程序路径
	exe, err := os.Executable()
	if err != nil {
		return nil
	}
	//备份当前程序：重命名
	err = os.Rename(exe, exe+".bak")
	if err != nil {
		return err
	}
	//更新写入
	err = os.WriteFile(localPath, body, os.FileMode(global.FilePer))
	if err != nil {
		return err
	}
	return nil
}
