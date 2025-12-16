package tests

import (
	"fmt"
	"testing"
	"xiaohongshu/app/infra/browser"
	"xiaohongshu/app/infra/eventbus"
	"xiaohongshu/app/services"
	"xiaohongshu/app/services/xiaohongshu/explore"
	"xiaohongshu/app/services/xiaohongshu/note"

	"github.com/playwright-community/playwright-go"
)

// TestXiaohongshuStartup 测试 Xiaohongshu 的 Startup 方法
func TestXiaohongshuStartup(t *testing.T) {

	newBrowser := browser.NewBrowser()
	err := newBrowser.Init()
	if err != nil {
		panic(err)
	}
	service, err := services.NewXiaohongshuService(newBrowser.GetBrowser())
	if err != nil {
		panic(err)
	}

	eventbus.GlobalEventBus.Subscribe("user-logged-in", func(userInfo interface{}) {
		fmt.Println(fmt.Printf("userInfo:%+v\n", userInfo))
	})
	err = service.Start()
	page := service.GetPage()
	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateDomcontentloaded,
	})
	if err != nil {
		panic(err)
	}
	newExplore := explore.NewExplore(service.GetPage())
	feeds, err := newExplore.Show()
	if err != nil && len(feeds) > 0 {
		return
	}
	page.Video()
	err = feeds[0].Element.Selector.Click()
	if err != nil {
		return
	}

	newNote := note.NewNote(service.GetPage())
	noteInfo, err := newNote.Show()
	if err != nil {
		return
	}

	fmt.Println(fmt.Sprintf("noteInfo:%+v\n", noteInfo))
	select {}
}
