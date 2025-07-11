package global

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	once        sync.Once
	allocCtx    context.Context
	allocCancel context.CancelFunc
	browserCtx  context.Context
)

func StartScreenshotWorkers(workers int) {
	if len(SsaveImgURLs) == 0 {
		return
	}
	if len(SsaveImgURLs) < workers {
		workers = len(SsaveImgURLs)
	}
	InitBrowser()
	defer ShutdownBrowser()

	_ = os.MkdirAll(SsaveIMGDIR, 0755)

	total := len(SsaveImgURLs)
	var finished int32 = 0
	var wg sync.WaitGroup
	taskChan := make(chan string, total)

	var lastStatus atomic.Value // 显示“✔ https://...”或“✘ https://...”

	printProgress := func(done, total int32, status string) {
		barWidth := 40
		percent := float64(done) / float64(total)
		doneBlocks := int(percent * float64(barWidth))
		bar := strings.Repeat("█", doneBlocks) + strings.Repeat("░", barWidth-doneBlocks)

		// 截断状态内容，最多显示50个字符，避免粘连或终端混乱
		truncate := func(s string, max int) string {
			if len(s) <= max {
				return s
			}
			return s[:max] + "..."
		}
		shortStatus := truncate(status, 50)

		// 输出进度并清除行尾（使用 ANSI 的 \033[K）
		fmt.Printf("\r📸 截图进度 [%s] %d/%d (%.1f%%) %s\033[K",
			bar, done, total, percent*100, shortStatus)
	}

	// 刷新进度条
	stopChan := make(chan struct{})
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				status, _ := lastStatus.Load().(string)
				printProgress(atomic.LoadInt32(&finished), int32(total), status)
			case <-stopChan:
				fmt.Printf("\r\033[2K") // 直接清除进度条这一整行
				//
				//status, _ := lastStatus.Load().(string)
				//printProgress(atomic.LoadInt32(&finished), int32(total), status)
				//time.Sleep(100 * time.Millisecond)
				//fmt.Printf("\r\033[2K\n") //清除整行 + 换行
				return

			}
		}
	}()

	// 启动 worker
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range taskChan {
				err := CaptureScreenshot(url, 90, SsaveIMGDIR)
				if err != nil {
					lastStatus.Store(fmt.Sprintf("✘ %s", url))
				} else {
					lastStatus.Store(fmt.Sprintf("✔ %s", url))
				}
				atomic.AddInt32(&finished, 1)
			}
		}()
	}

	// 启动任务
	saveImgMu.Lock()
	for _, url := range SsaveImgURLs {
		taskChan <- url
	}
	saveImgMu.Unlock()
	close(taskChan)

	wg.Wait()
	close(stopChan)

	couunt, err := CountDirFiles(SsaveIMGDIR)
	if couunt == 0 && err != nil {
		return
	}
	fmt.Printf("[*] Web扫描截图保存目录：%v 当前共计截图数量：%v\n", SsaveIMGDIR, couunt)
}

// InitBrowser 初始化共享 Chrome 实例
func InitBrowser() {
	once.Do(func() {

		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
			chromedp.Flag("ignore-certificate-errors", true),
			chromedp.WindowSize(1920, 1080),
		)

		allocCtx, allocCancel = chromedp.NewExecAllocator(context.Background(), opts...)
		browserCtx, _ = chromedp.NewContext(allocCtx)

		// 启动浏览器连接
		_ = chromedp.Run(browserCtx)
	})
}

// ShutdownBrowser 关闭共享 Chrome 实例
func ShutdownBrowser() {
	if allocCancel != nil {
		allocCancel()
	}
}

// GetBrowserContext 返回共享上下文
func GetBrowserContext() context.Context {
	return browserCtx
}

// CaptureScreenshot 截图任务，保存为 PNG 文件
func CaptureScreenshot(url string, quality int64, dir string) error {
	// 创建新标签页（共享 Chrome 实例）
	ctx, cancel := chromedp.NewContext(GetBrowserContext())
	defer cancel()

	// 设置超时时间
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// 导航 + 截图
	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second), // 等待页面渲染
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      0,
					Y:      0,
					Width:  1920,
					Height: 1080,
					Scale:  1,
				}).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return err
	}

	// 生成合法文件名
	filename := strings.Map(func(r rune) rune {
		if r == '/' || r == ':' || r == '?' || r == '&' {
			return '_'
		}
		return r
	}, url)

	// 保存 PNG 文件
	output := filepath.Join(dir, filename+".png")
	return os.WriteFile(output, buf, 0644)
}
