package cmd

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golin/global"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "检查更新程序",
	Long:  `通过api.github.com进行检查更新程序`,
	Run: func(cmd *cobra.Command, args []string) {
		newrelease, err := global.CheckForUpdate()
		if err != nil {
			fmt.Println("更新失败:", err)
			return
		}
		var input string
		fmt.Print("是否下载GitHub仓库最新版本程序？y/n: ")
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
				BrowserDownloadUrl, savaname := "", ""
				switch runtime.GOOS {
				case "windows":
					savaname = "golin.exe"
					BrowserDownloadUrl = fmt.Sprintf("https://github.com/selinuxG/Golin/releases/download/%s/%s", newrelease.TagName, savaname)
				case "linux":
					savaname = "golin_linux_amd64"
					BrowserDownloadUrl = fmt.Sprintf("https://github.com/selinuxG/Golin/releases/download/%s/%s", newrelease.TagName, savaname)
				case "drawin":
					savaname = "golin_drawin_amd64"
					BrowserDownloadUrl = fmt.Sprintf("https://github.com/selinuxG/Golin/releases/download/%s/%s", newrelease.TagName, savaname)
				}
				err := downloadFile(BrowserDownloadUrl, savaname, proxy)
				if err != nil {
					fmt.Println("更新失败->", err)
					return
				}
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

// getFileSize 返回文件大小
func getFileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	fileInfo.ModTime()

	return fileInfo.Size(), nil
}

// downloadFile 下载更新
func downloadFile(downurl, localPath, proxy string) error {
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

	// 获取文件大小
	fileSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	var progressed int

	// 获取当前程序完整位置
	exe, err := os.Executable()
	if err != nil {
		return nil
	}
	// 当前旧版程序重命名实现备份效果
	os.Remove(exe + ".bak") //先移除之前的备份文件
	err = os.Rename(exe, exe+".bak")
	if err != nil {
		return err
	}
	// 写入最新版
	out, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer out.Close()

	reader := bufio.NewReader(resp.Body)
	buff := make([]byte, 1024)
	for {
		n, err := reader.Read(buff)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		_, err = out.Write(buff[:n])
		if err != nil {
			return err
		}
		progressed += n
		percentage := float64(progressed) / float64(fileSize) * 100

		percentStr := fmt.Sprintf("%.2f", percentage) // 将百分比值格式化为字符串
		fmt.Printf("\r更新进度: %s",
			color.RedString("%s", fmt.Sprintf("%s%%", percentStr)),
		)

	}
	fmt.Printf("\r更新进度: %s\n",
		color.GreenString("%s", fmt.Sprintf("%-7s", "100%")),
	)
	return nil
}
