package global

import (
	"context"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CaptureScreenshot 要截图的网页的 URL，截图的质量，以及保存截图的目录
func CaptureScreenshot(url string, quality int64, dir string) error {
	// 创建一个上下文
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 设置浏览器选项
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("ignore-certificate-errors", true), // 忽略证书错误
		//chromedp.Flag("remote-debugging-port", "9222"),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// 创建一个浏览器实例
	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	// 导航到指定的URL
	var buf []byte
	err := chromedp.Run(ctx, chromedp.Navigate(url), chromedp.Sleep(3*time.Second), chromedp.ActionFunc(func(ctx context.Context) error {
		// 获取页面截图
		var err error
		buf, err = page.CaptureScreenshot().WithQuality(quality).WithClip(&page.Viewport{X: 0, Y: 0, Width: 1024, Height: 768, Scale: 1}).Do(ctx)
		if err != nil {
			return err
		}
		return nil
	}))
	if err != nil {
		return err
	}

	// 将 URL 中的非法字符替换为下划线
	filename := strings.Map(func(r rune) rune {
		if r == '/' || r == ':' {
			return '_'
		}
		return r
	}, url)

	// 检查文件夹是否存在，如果不存在则创建
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	// 将截图保存到文件
	err = os.WriteFile(filepath.Join(dir, filename+".png"), buf, 0644)
	if err != nil {
		return err
	}

	return nil
}
