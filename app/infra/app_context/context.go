package app_context

import (
	"context"
	"log"
	"path"
	"sync"
	"xiaohongshu/app/infra/browser"
	"xiaohongshu/app/infra/db"
	"xiaohongshu/app/pkg/utils"

	"github.com/playwright-community/playwright-go"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type AppContext struct {
	environmentInfo runtime.EnvironmentInfo
	rootPath        string
	browser         *browser.Browser
	errors          []error
	initComplete    chan struct{}
	initOnce        sync.Once
	// 添加用于控制持续发送事件的字段
	stopSending    chan struct{}
	sendingStopped chan struct{}
}

func NewContext() *AppContext {
	return &AppContext{
		browser:        browser.NewBrowser(),
		errors:         make([]error, 0),
		initComplete:   make(chan struct{}),
		stopSending:    make(chan struct{}),
		sendingStopped: make(chan struct{}),
	}
}

func (c *AppContext) OnStartup(ctx context.Context) {
	c.environmentInfo = runtime.Environment(ctx)
	var err error
	c.rootPath, err = utils.GetPath(c.environmentInfo.BuildType)
	if err != nil {
		c.errors = append(c.errors, err)
		close(c.initComplete) // 确保即使出错也能解除阻塞
		return
	}
	go func() {
		defer close(c.initComplete) // 确保在函数结束时关闭通道
		// 收集浏览器下载过程中的错误
		if err := c.browser.Init(); err != nil {
			c.errors = append(c.errors, err)
			log.Printf("Failed to download browser: %v", err)
			return
		}

		// 收集数据库初始化过程中的错误
		if _, err := db.Init(path.Join(c.rootPath, "app.db")); err != nil {
			c.errors = append(c.errors, err)
			log.Printf("Failed to initialize database: %v", err)
			return
		}
	}()
}

func (c *AppContext) OnDomReady(ctx context.Context) {
	runtime.LogPrint(ctx, "OnDomReady start")
	// 等待初始化完成
	<-c.initComplete
	runtime.LogPrint(ctx, "OnDomReady end")

	runtime.EventsOn(ctx, "startReady", func(optionalData ...interface{}) {
		errorMessages := make([]string, len(c.errors))
		for i, err := range c.errors {
			errorMessages[i] = err.Error()
		}
		runtime.EventsEmit(ctx, "initialization-complete", errorMessages)

	})
}

func (c *AppContext) OnMount(ctx context.Context) {
	runtime.EventsEmit(ctx, "on-mount", nil)
}

// GetErrors 返回收集到的所有错误
func (c *AppContext) GetErrors() []error {
	return c.errors
}

func (c *AppContext) GetRootPath() string {
	return c.rootPath
}

func (c *AppContext) GetPlaywrightBrowser() playwright.Browser {
	return c.browser.GetBrowser()
}
