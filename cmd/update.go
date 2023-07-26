package cmd

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"golin/global"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
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

		// 获取当前程序大小并对比远程仓库的大小
		exePath, _ := os.Executable()
		fileSize, _ := getFileSize(exePath)
		if global.Version == newrelease.TagName && fileSize == newrelease.Assets[0].Size {
			fmt.Println("当前版本已为最新版本,无需更新。")
			return
		}
		fmt.Println(fmt.Sprintf("当前版本为:%s 大小为:%v/KB  -> 最新版本为:%s 大小为:%v/KB", global.Version, fileSize, newrelease.TagName, newrelease.Assets[0].Size))
		var input string
		fmt.Print("是否更新？y/n: ")
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
					fmt.Println("Failed!")
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
