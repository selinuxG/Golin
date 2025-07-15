package global

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
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

var ScreenshotCtx context.Context
var ScreenshotCancel context.CancelFunc

func init() {
	if SaveIMG {
		ScreenshotCtx, ScreenshotCancel = context.WithCancel(context.Background())
		// 监听 Ctrl+C
		go func() {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)
			<-sigChan
			CancelScreenshot()
		}()
	}
}

func StartScreenshotWorkers(workers int) {
	if len(SsaveImgURLs) == 0 { //不判断是否开启是因为漏洞截图不受状态影响
		return
	}
	if len(SsaveImgURLs) < workers {
		workers = len(SsaveImgURLs)
	}

	_, err := DetectChromePath()
	if err != nil {
		return
	}

	InitBrowser()
	defer ShutdownBrowser()

	_ = os.MkdirAll(SsaveIMGDIR, 0755)

	total := len(SsaveImgURLs)
	var finished int32 = 0
	var wg sync.WaitGroup
	taskChan := make(chan string, total)

	var lastStatus atomic.Value

	printProgress := func(done, total int32, status string) {
		barWidth := 40
		percent := float64(done) / float64(total)
		doneBlocks := int(percent * float64(barWidth))
		bar := strings.Repeat("█", doneBlocks) + strings.Repeat("░", barWidth-doneBlocks)

		truncate := func(s string, max int) string {
			if len(s) <= max {
				return s
			}
			return s[:max] + "..."
		}
		shortStatus := truncate(status, 50)

		fmt.Printf("\r[-] 📸 截图进度 [%s] %d/%d (%.1f%%) %s(可随时CTRL+C取消此项)\033[K",
			bar, done, total, percent*100, shortStatus)
	}

	stopChan := make(chan struct{})
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				status, _ := lastStatus.Load().(string)
				printProgress(atomic.LoadInt32(&finished), int32(total), status)
			case <-ScreenshotCtx.Done():
				//fmt.Printf("\r\033[2K[!] 已中断截图任务\n")
				return
			case <-stopChan:
				fmt.Printf("\r\033[2K")
				return
			}
		}
	}()

	// worker
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ScreenshotCtx.Done():
					return
				case url, ok := <-taskChan:
					if !ok {
						return
					}
					err = CaptureScreenshot(url, 90, SsaveIMGDIR)
					if err != nil {
						lastStatus.Store(fmt.Sprintf("✘ %s", url))
					} else {
						lastStatus.Store(fmt.Sprintf("✔ %s", url))
					}
					atomic.AddInt32(&finished, 1)
				}
			}
		}()
	}

	// 分发任务
	saveImgMu.Lock()
	for _, url := range SsaveImgURLs {
		select {
		case <-ScreenshotCtx.Done():
			break
		default:
			taskChan <- url
		}
	}
	saveImgMu.Unlock()
	close(taskChan)

	wg.Wait()
	close(stopChan)
	count, err := CountDirFiles(SsaveIMGDIR)
	if count == 0 && err != nil {
		return
	}
	if ScreenshotCtx.Err() != nil {
		fmt.Printf("[!] 截图任务被取消，跳过剩余任务")
		fmt.Printf("\033[2K\r[*] Web扫描截图保存目录：%v 当前共计截图数量：%v\n", SsaveIMGDIR, count)
		return
	}
	fmt.Printf("\033[2K\r[*] Web扫描截图保存目录：%v 当前共计截图数量：%v\n", SsaveIMGDIR, count)
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
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
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

func DetectChromePath() (string, error) {
	var locations []string
	switch runtime.GOOS {
	case "darwin":
		locations = []string{
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		}
	case "windows":
		locations = []string{
			"chrome", "chrome.exe",
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\Google\Chrome\Application\chrome.exe`),
			filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\Chromium\Application\chrome.exe`),
		}
	default:
		locations = []string{
			"headless_shell", "headless-shell", "chromium", "chromium-browser",
			"google-chrome", "google-chrome-stable", "google-chrome-beta", "google-chrome-unstable",
			"/usr/bin/google-chrome", "/usr/local/bin/chrome", "/snap/bin/chromium", "chrome",
		}
	}

	for _, name := range locations {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("[✘] 未找到 Chrome 可执行文件，请安装 Google Chrome 或 Chromium")
}

// CancelScreenshot 中断截图任务
func CancelScreenshot() {
	if ScreenshotCancel != nil {
		fmt.Printf("\r[!] 用户按下 Ctrl+C,已中断截图任务,请等待已下发任务结束%s", strings.Repeat(" ", 50))
		ScreenshotCancel()
	}
}
