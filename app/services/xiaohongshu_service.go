package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"xiaohongshu/app/entities"
	"xiaohongshu/app/infra/eventbus"
	"xiaohongshu/app/pkg/utils"
	scripts2 "xiaohongshu/app/services/xiaohongshu/scripts"

	"github.com/playwright-community/playwright-go"
)

type XiaohongshuService struct {
	browser           playwright.Browser
	page              playwright.Page
	accountCookiePath *string
	cookiePath        string
}

func NewXiaohongshuService(browser playwright.Browser) (*XiaohongshuService, error) {
	directory, err := utils.GetDefaultCacheDirectory()
	if err != nil {
		return nil, err
	}
	cookiePath := path.Join(directory, "xiaohongshu.cookies")

	// 检查 cookie 文件是否存在
	var accountCookiePath *string
	if _, err := os.Stat(cookiePath); err == nil {
		// 文件存在，将 cookiePath 赋值给 accountCookiePath 字段
		accountCookiePath = &cookiePath
	} else {
		// 文件不存在，将 nil 赋值给 accountCookiePath 字段
		accountCookiePath = nil
	}
	scripts2.InitEventBus()
	return &XiaohongshuService{
		browser:           browser,
		accountCookiePath: accountCookiePath,
		cookiePath:        cookiePath,
	}, nil
}

func (s *XiaohongshuService) Start() error {
	var err error
	context, err := s.browser.NewContext(playwright.BrowserNewContextOptions{
		StorageStatePath: s.accountCookiePath,
	})
	if err != nil {
		return err
	}
	s.page, err = context.NewPage()
	if err != nil {
		return err
	}
	// 添加反检测脚本，隐藏webdriver标志
	scriptContent := `
			delete navigator.__proto__.webdriver;
			window.chrome = {runtime: {}};
			window.test = "添加反检测脚本，隐藏webdriver标志";
			Object.defineProperty(navigator, 'languages', {
				get: () => ['en-US', 'en']
			});
			Object.defineProperty(navigator, 'plugins', {
				get: () => [1, 2, 3, 4, 5]
			});
		`
	_ = s.page.AddInitScript(playwright.Script{
		Content: &scriptContent,
	})
	js := scripts2.ToolJs
	_ = s.page.AddInitScript(playwright.Script{
		Content: &js,
	})
	// 设置视口大小，模拟真实浏览器
	err = s.page.SetViewportSize(1366, 768)
	if err != nil {
		return err
	}
	_ = scripts2.InjectMediaCaptureScript(s.page)
	// 导航到页面
	_, err = s.page.Goto("https://www.xiaohongshu.com")
	if err != nil {
		fmt.Println(fmt.Sprintf("GotoURL err: %s", err))
		return err
	}

	s.page.OnResponse(s.onResponse)
	return err
}

func (s *XiaohongshuService) GetPage() playwright.Page {
	return s.page
}

func (s *XiaohongshuService) onResponse(response playwright.Response) {
	// 检查URL是否包含v2/user/me
	go s.me(response)
}

func (s *XiaohongshuService) me(response playwright.Response) {
	if strings.Contains(response.URL(), "v2/user/me") {
		// 检查响应状态是否成功
		// 读取响应体
		go func() {
			body, err := response.Body()
			if err != nil {
				fmt.Printf("获取响应体失败: %v\n", err)
				return
			}
			// 解析响应内容
			var apiResponse entities.ApiResponse
			err = json.Unmarshal(body, &apiResponse)
			if err != nil {
				return
			}
			if apiResponse.Code != 0 {
				fmt.Printf("API返回错误码: %d, 消息: %s\n", apiResponse.Code, apiResponse.Msg)
				return
			}
			context := s.page.Context()
			_, err = context.StorageState(s.cookiePath)
			if err != nil {
				return
			}
			// 检查API响应是否成功
			if apiResponse.Success {
				// 通过event_bus发送用户信息
				eventbus.GlobalEventBus.Publish("user-logged-in", apiResponse.Data)
			} else {
				fmt.Printf("API调用不成功: %+v\n", apiResponse)
			}
		}()

	}
}
