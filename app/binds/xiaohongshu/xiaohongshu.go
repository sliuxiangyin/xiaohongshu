package xiaohongshu

import (
	"context"
	"embed"
	"fmt"
	"time"
	"xiaohongshu/app/infra/app_context"
	"xiaohongshu/app/infra/eventbus"
	"xiaohongshu/app/services"
	"xiaohongshu/app/services/xiaohongshu/explore"
	"xiaohongshu/app/services/xiaohongshu/note"

	"github.com/playwright-community/playwright-go"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Xiaohongshu struct
type Xiaohongshu struct {
	appContext  *app_context.AppContext
	ctx         context.Context
	service     *services.XiaohongshuService
	page        playwright.Page
	explorePage *explore.Explore
	scriptPath  embed.FS
}

// NewXiaohongshu creates a new Xiaohongshu application struct
func NewXiaohongshu(appContext *app_context.AppContext, scriptPath embed.FS) *Xiaohongshu {
	return &Xiaohongshu{
		appContext: appContext,
		scriptPath: scriptPath,
	}
}

// Startup 初始化上下文
func (x *Xiaohongshu) Startup(ctx context.Context) {
	x.ctx = ctx
}

// OnDomReady 当DOM准备就绪时调用
func (x *Xiaohongshu) OnDomReady(ctx context.Context) {
	if x.service != nil {
		return
	}
	var err error
	x.service, err = services.NewXiaohongshuService(x.appContext.GetPlaywrightBrowser(), x.scriptPath)
	if err != nil {
		return
	}
	err = x.service.Start()
	if err != nil {
		return
	}
	x.page = x.service.GetPage()
	x.explorePage = explore.NewExplore(x.page)
	// 监听用户登录事件
	eventbus.GlobalEventBus.Subscribe("user-logged-in", func(userInfo interface{}) {
		runtime.EventsEmit(ctx, "user-logged-in", userInfo)
	})
}

// NextPage 下一页功能
func (x *Xiaohongshu) NextPage() error {
	// TODO: 实现下一页逻辑
	if x.page == nil {
		return fmt.Errorf("page is not initialized")
	}
	x.page.Evaluate(`console.log("window.__MediaCaptureController",window.__MediaCaptureController);`)
	x.page.Evaluate(`console.log("window.test",window.test );`)
	return nil
}

// Refresh 刷新功能
func (x *Xiaohongshu) Refresh() error {

	if x.page == nil {
		return fmt.Errorf("page is not initialized")
	}
	err := x.explorePage.RefreshPage()
	if err != nil {
		return err
	}
	return nil
}

// GetItems 获取列表项数据
func (x *Xiaohongshu) GetItems() ([]map[string]interface{}, error) {
	if x.page == nil {
		return nil, fmt.Errorf("page is not initialized")
	}

	feeds, err := x.explorePage.Show()
	if err != nil {
		return nil, fmt.Errorf("failed to get explore feeds: %v", err)
	}

	// 将 FeedsInfo 转换为 map[string]interface{} 以便前端使用
	items := make([]map[string]interface{}, 0, len(feeds))
	for _, feed := range feeds {
		item := map[string]interface{}{
			"index":         feed.Index,
			"title":         feed.Title.Text,
			"coverImageUrl": feed.Cover.Text,
			"username":      feed.User.Text,
			"avatarUrl":     feed.Avatar.Text,
		}
		items = append(items, item)
	}

	return items, nil
}

// OnItemClick 当列表项被点击时调用
func (x *Xiaohongshu) OnItemClick(index int) error {

	fmt.Println("OnItemClick", index)
	feed, err := x.explorePage.GetFeed(index)
	if err != nil {
		return err
	}

	if err != nil {
		fmt.Println(fmt.Sprintf("watch playwright error: %v", err))
		return err
	}
	err = feed.Element.Click()
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 1)
	newNote, err := note.NewNote(x.page).Show()
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to get new note: %v", err))
		return err
	}

	video := newNote.Video()
	if video == nil {
		fmt.Println(fmt.Sprintf("New video is nil"))
		return nil
	}
	err = video.MediaStart()
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to start video: %v", err))
		return err
	}
	err = video.ListenVideoState(func(b bool) {

		if video.IsMute() {
			video.ToggleVolume()
		}
		fmt.Println(fmt.Sprintf("New video listened: %+v", b))
		fmt.Println(fmt.Sprintf("IsMute: %v", video.IsMute()))
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to listen: %v", err))
		return err
	}

	err = video.ListenVideoAudio(func(audio note.VideoAudio) {
		fmt.Println(fmt.Sprintf("New audio listenedddddddd: %+v", audio.Buffer))
	})
	err = video.MediaStart()
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to start media: %v", err))
		return err
	}
	err = video.ListenVideoFrame(func(frame note.VideoFrame) {
		fmt.Println(fmt.Sprintf("New frame listeneddddddddddd: %+v", frame.Data))
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to listen: %v", err))
		return err
	}
	return nil
}
