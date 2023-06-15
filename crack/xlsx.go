package crack

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
	"golin/global"
	"os"
	"path/filepath"
	"strings"
)

var (
	path          = "" //待破解文件路径
	df_xlsxpasswd = []string{"123456", "1234567", "12345678", "123456789", "1234567890", "666666", "888888", "000000", "111111"}
)

func XlsxCheck(cmd *cobra.Command, args []string) {

	path, _ = cmd.Flags().GetString("file")
	passpath, _ := cmd.Flags().GetString("passwdfile")
	chancount, _ := cmd.Flags().GetInt("chan")

	if !global.PathExists(path) {
		fmt.Printf("%s 文件不存在,检查后再来吧！拜拜\n", path)
		os.Exit(1)
	}
	if filepath.Ext(path) != ".xlsx" {
		fmt.Printf("%s 格式错误，只能是.xlsx文件! 检查后再来吧！拜拜\n", path)
		os.Exit(1)
	}

	if global.PathExists(passpath) {
		data, _ := os.ReadFile(passpath)
		strdata := string(data)
		strdata = strings.ReplaceAll(strdata, "\r\n", "\n")
		for _, v := range strings.Split(strdata, "\n") {
			df_xlsxpasswd = append(df_xlsxpasswd, v)
		}
	}
	ch = make(chan struct{}, chancount)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //确保所有的goroutine都已经退出
	fmt.Printf("[*] 开始检测文件：%s 尝试次数：%d 线程数：%d 祝君好运:)\n", path, len(df_xlsxpasswd), chancount)
	for _, pass := range removeDuplicates(df_xlsxpasswd) {
		ch <- struct{}{}
		wg.Add(1)
		go checkExeclPass(ctx, cancel, pass)
	}
	wg.Wait()
	fmt.Print("\033[2K") // 擦除整行
	fmt.Print("\r")      // 光标移动到行首
}

func checkExeclPass(ctx context.Context, cancel context.CancelFunc, pass string) {
	defer func() {
		wg.Done()
		<-ch
	}()
	select {
	case <-ctx.Done():
		return
	default:
	}
	fmt.Printf("\r[-] %s 尝试口令：%s", path, pass)
	f, err := excelize.OpenFile(path, excelize.Options{Password: pass})
	if err == nil {
		fmt.Printf("\r[*] 发现口令！%s----->%s\n", path, pass)
		defer f.Close()
		cancel()
		return
	}
	return
}
