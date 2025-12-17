package scripts

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
	"github.com/spf13/cast"
)

// VideoFrame 视频帧数据结构
type VideoFrame struct {
	Width  int64       `json:"width"`
	Height int64       `json:"height"`
	Data   interface{} `json:"data"`
	Ts     float64     `json:"ts"`
}

// VideoAudio 视频音频数据结构
type VideoAudio struct {
	SampleRate int64       `json:"sampleRate"`
	Channels   int64       `json:"channels"`
	Frames     int64       `json:"frames"`
	Buffer     interface{} `json:"buffer"`
	Ts         float64     `json:"ts"`
}

// MediaStart 启动媒体捕获
func MediaStart(element playwright.Locator) error {
	// 首先检查MediaCaptureController是否存在
	exists, err := element.Evaluate(`(videoEl) => {
		return window.__MediaCaptureController !== undefined;
	}`, nil)
	if err != nil {
		return fmt.Errorf("检查MediaCaptureController失败: %w", err)
	}

	if !exists.(bool) {
		return fmt.Errorf("MediaCaptureController未定义，请确保脚本已正确注入")
	}

	_, err = element.EvaluateHandle(`(videoEl) => {
		window.__MediaCaptureController.start(videoEl);
	}`, nil)
	return err
}

// MediaStop 停止媒体捕获
func MediaStop(element playwright.Locator) error {
	_, err := element.EvaluateHandle(`(element) => {
		window.__MediaCaptureController.stop(element);
	}`, nil)
	return err
}

// MediaDestroy 销毁媒体捕获
func MediaDestroy(element playwright.Locator) error {
	_, err := element.EvaluateHandle(`(element) => {
window.__MediaCaptureController.destroy(element);
	}`, nil)
	return err
}

// MediaStopAll 停止所有媒体捕获
func MediaStopAll(element playwright.Locator) error {

	_, err := element.EvaluateHandle(`(element) => {
window.__MediaCaptureController.stopAll();
	}`, nil)
	return err
}

// MediaDestroyAll 销毁所有媒体捕获
func MediaDestroyAll(element playwright.Locator) error {

	_, err := element.EvaluateHandle(`(element) => {
window.__MediaCaptureController.destroyAll(element);
	}`, nil)
	return err
}

// MediaListenVideoState  监听视频播放状态变化
func MediaListenVideoState(element playwright.Locator) error {

	_, err := element.EvaluateHandle(`
		(video) => {
			// 确保只添加一次事件监听器
			if (video.__listening) return;
			video.__listening = true;
			
			const notifyState = () => {
				const isPlaying = !video.paused && !video.ended && video.readyState > 2;
				window.__onVideoStateChange?.(isPlaying);
			};
			
			// 监听相关的事件
			video.addEventListener('play', notifyState);
			video.addEventListener('pause', notifyState);
			video.addEventListener('ended', notifyState);
			video.addEventListener('loadstart', notifyState);
			video.addEventListener('canplay', notifyState);
			video.addEventListener('canplaythrough', notifyState);
			video.addEventListener('waiting', notifyState);
			
			// 初始状态通知
			notifyState();
		}
	`, nil)
	return err
}

// MediaRemoveVideoStateListener 移除视频状态监听
func MediaRemoveVideoStateListener(element playwright.Locator) error {
	page, err := element.Page()
	if err != nil {
		return err
	}
	// 移除JavaScript中的事件监听器
	_, err = element.Evaluate(`
		(video) => {
			if (!video.__listening) return;
			video.__listening = false;
			
			const notifyState = () => {
				const isPlaying = !video.paused && !video.ended && video.readyState > 2;
				window.__onVideoStateChange?.(isPlaying);
			};
			
			video.removeEventListener('play', notifyState);
			video.removeEventListener('pause', notifyState);
			video.removeEventListener('ended', notifyState);
			video.removeEventListener('loadstart', notifyState);
			video.removeEventListener('canplay', notifyState);
			video.removeEventListener('canplaythrough', notifyState);
			video.removeEventListener('waiting', notifyState);
		}
	`, nil)
	if err != nil {
		return err
	}
	// 移除暴露的函数
	_, err = page.Evaluate("delete window.__onVideoStateChange", nil)
	return err
}

// InjectMediaCaptureScript 注入媒体捕获脚本到页面
func InjectMediaCaptureScript(page playwright.Page, scriptContent string) error {
	var err error
	// 暴露回调函数给浏览器
	err = page.ExposeFunction("__onVideoStateChange", func(args ...interface{}) interface{} {
		if len(args) > 0 {
			if isPlaying, ok := args[0].(bool); ok {
				GetEventBus().Publish("media:video:state", isPlaying)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	// 首先获取页面实例
	err = page.ExposeFunction("__onVideoFrame", func(args ...interface{}) interface{} {
		if len(args) > 0 {
			if data, ok := args[0].(map[string]interface{}); ok {
				frame := VideoFrame{
					Width:  cast.ToInt64(data["width"]),
					Height: cast.ToInt64(data["height"]),
					Data:   data["data"],
					Ts:     cast.ToFloat64(data["ts"]),
				}
				// 通过事件总线发送视频帧数据
				fmt.Println("media:video:frame", GetEventBus().HasCallback("media:video:frame"))
				GetEventBus().Publish("media:video:frame", frame)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	// 暴露__onVideoAudio回调函数到浏览器环境
	err = page.ExposeFunction("__onVideoAudio", func(args ...interface{}) interface{} {
		if len(args) > 0 {
			if data, ok := args[0].(map[string]interface{}); ok {
				audio := VideoAudio{
					SampleRate: cast.ToInt64(data["sampleRate"]),
					Channels:   cast.ToInt64(data["channels"]),
					Frames:     cast.ToInt64(data["frames"]),
					Buffer:     data["buffer"],
					Ts:         cast.ToFloat64(data["ts"]),
				}
				fmt.Println(fmt.Sprintf("__onVideoAudio: %+v", audio))
				// 通过事件总线发送音频数据
				GetEventBus().Publish("media:video:audio", audio)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = page.AddInitScript(playwright.Script{
		Content: &scriptContent,
	})
	return err
}
