package browser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"xiaohongshu/app/pkg/utils"

	"github.com/playwright-community/playwright-go"
)

// Browser 封装了playwright浏览器实例
type Browser struct {
	pw              *playwright.Playwright
	browser         playwright.Browser
	driverDirectory string
	initialized     bool
	mu              sync.Mutex
}

// NewBrowser 创建一个新的浏览器实例
func NewBrowser() *Browser {
	return &Browser{}
}

// Init 初始化浏览器，首先调用Download方法下载浏览器（如果尚未安装）
func (b *Browser) Init() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 调用Download方法下载浏览器（如果尚未安装）
	err := b.download()
	if err != nil {
		return err
	}
	err = b.launch()
	if err != nil {
		return err
	}
	return nil
}

// Download 下载
func (b *Browser) download() error {
	// 获取默认缓存目录
	cacheDirectory, err := utils.GetDefaultCacheDirectory()
	if err != nil {
		return fmt.Errorf("failed to get default cache directory: %w", err)
	}
	b.driverDirectory = filepath.Join(cacheDirectory, "ms-playwright-go")

	// 检查是否已安装（通过检测node可执行文件是否存在）
	nodeExecutable := getNodeExecutable(b.driverDirectory)
	if _, err := os.Stat(nodeExecutable); err == nil {
		// 已安装，跳过重复安装
		log.Println("Playwright is already installed, skipping installation")
		return nil
	}

	// 未安装，执行安装
	err = playwright.Install(&playwright.RunOptions{
		DriverDirectory: b.driverDirectory,
		Browsers:        []string{"chromium"},
		Verbose:         true,
	})

	if err != nil {
		return err
	}
	return nil
}

// LaunchExisting 启动已安装的Chrome浏览器（如果存在）
func (b *Browser) launch() error {
	// 确保浏览器实例只能被启动一次
	if b.initialized {
		return nil
	}
	// 初始化playwright
	pw, err := playwright.Run(&playwright.RunOptions{DriverDirectory: b.driverDirectory})
	if err != nil {
		return fmt.Errorf("failed to start playwright: %w", err)
	}
	b.pw = pw

	// 尝试启动已安装的Chrome浏览器
	//options := playwright.BrowserTypeLaunchOptions{
	//	Headless: playwright.Bool(false), // 使用有头模式
	//	Args: []string{
	//		"--no-sandbox",
	//		"--disable-blink-features=AutomationControlled",
	//		"--disable-extensions",
	//		"--disable-plugins",
	//		"--disable-plugins-discovery",
	//		"--disable-web-security",
	//		"--disable-features=IsolateOrigins,site-per-process",
	//	},
	//}
	b.browser, err = pw.Chromium.ConnectOverCDP("http://localhost:9222")
	if err != nil {
		return err
	}
	// 启动Chrome浏览器
	// b.browser, err = pw.Chromium.Launch(options)
	// if err != nil {
	// return fmt.Errorf("failed to launch chrome browser: %w", err)
	// }

	// 标记浏览器已初始化
	b.initialized = true

	log.Println("Chrome browser launched successfully")
	return nil
}

// Close 关闭浏览器和playwright实例
func (b *Browser) Close() error {
	if b.browser != nil {
		if err := b.browser.Close(); err != nil {
			return fmt.Errorf("failed to close browser: %w", err)
		}
	}

	if b.pw != nil {
		if err := b.pw.Stop(); err != nil {
			return fmt.Errorf("failed to stop playwright: %w", err)
		}
	}

	return nil
}

// GetBrowser 返回底层的playwright浏览器实例
func (b *Browser) GetBrowser() playwright.Browser {
	return b.browser
}

// 获取node可执行文件路径
func getNodeExecutable(driverDirectory string) string {
	node := "node"
	if runtime.GOOS == "windows" {
		node = "node.exe"
	}
	return filepath.Join(driverDirectory, node)
}
