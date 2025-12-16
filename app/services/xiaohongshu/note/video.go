package note

import (
	"fmt"
	"xiaohongshu/app/services/xiaohongshu/scripts"

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

// VideoState 视频状态枚举
type VideoState int

const (
	HAVE_NOTHING      VideoState = iota // 0 = 没有关于音频/视频是否就绪的信息
	HAVE_METADATA                       // 1 = 关于音频/视频就绪的元数据
	HAVE_CURRENT_DATA                   // 2 = 关于当前播放位置的数据是可用的
	HAVE_FUTURE_DATA                    // 3 = 当前及至少下一帧的数据是可用的
	HAVE_ENOUGH_DATA                    // 4 = 可用数据足以开始播放
)

type Video struct {
	locator      playwright.Locator
	videoElement playwright.Locator
}

func NewVideo(locator playwright.Locator) *Video {
	return &Video{
		locator:      locator,
		videoElement: locator.Locator("video"),
	}
}

//player-el xgplayer xgplayer-pc xhsplayer-skin-default xgplayer-pause xgplayer-volume-muted
//player-el xgplayer xgplayer-pc xhsplayer-skin-default xgplayer-pause xgplayer-volume-large

// IsPlayable  是否在播放
func (v *Video) IsPlayable() bool {
	// 使用JavaScript检查视频是否正在播放
	result, err := v.videoElement.Evaluate(`(video) => !video.paused && !video.ended && video.readyState > 2`, nil)
	if err != nil {
		return false
	}

	if isPlaying, ok := result.(bool); ok {
		return isPlaying
	}
	return false
}

// IsMute  是否在静音
func (v *Video) IsMute() bool {
	count, err := v.locator.Locator(".xgplayer-volume-muted").Count()
	if err != nil {
		return false
	}
	return count != 0
}

func (v *Video) TogglePlay() bool {
	_ = v.videoElement.Click()
	return v.IsPlayable()
}

func (v *Video) ToggleVolume() bool {
	_ = v.locator.Locator(".xgplayer-volume .xgplayer-icon").Click()
	return v.IsMute()
}

// GetVideoState 获取视频当前的加载/播放状态
func (v *Video) GetVideoState() (VideoState, error) {
	result, err := v.videoElement.Evaluate(`(video) => video.readyState`, nil)
	if err != nil {
		return HAVE_NOTHING, err
	}

	readyState := cast.ToInt(result)
	return VideoState(readyState), nil
}

// ListenVideoState 监听视频播放状态变化
func (v *Video) ListenVideoState(handler func(bool)) error {
	err := scripts.MediaListenVideoState(v.videoElement, handler)
	if err != nil {
		return err
	}
	return nil
}

func (v *Video) RemoveVideoState() error {
	return scripts.MediaRemoveVideoStateListener(v.videoElement)
}

// MediaStart 启动媒体捕获
func (v *Video) MediaStart() error {
	// 启动媒体捕获
	return scripts.MediaStart(v.videoElement)
}

// MediaStop 停止媒体捕获
func (v *Video) MediaStop() error {
	return scripts.MediaStop(v.videoElement)
}

// MediaDestroy 销毁媒体捕获
func (v *Video) MediaDestroy() error {
	return scripts.MediaDestroy(v.videoElement)
}

// MediaStopAll 停止所有媒体捕获
func (v *Video) MediaStopAll() error {
	return scripts.MediaStopAll(v.videoElement)
}

// MediaDestroyAll 销毁所有媒体捕获
func (v *Video) MediaDestroyAll() error {
	return scripts.MediaDestroyAll(v.videoElement)
}

// ListenVideoFrame 订阅视频帧数据事件
func (v *Video) ListenVideoFrame(handler func(VideoFrame)) error {
	err := scripts.GetEventBus().Subscribe("media:video:frame", func(frame interface{}) {
		if videoFrame, ok := frame.(VideoFrame); ok {
			handler(videoFrame)
		}
	})
	return err
}

// ListenVideoAudio 订阅视频音频数据事件
func (v *Video) ListenVideoAudio(handler func(VideoAudio)) error {
	err := scripts.GetEventBus().Subscribe("media:video:audio", func(audio interface{}) {

		fmt.Println("订阅视频音频数据事件 media:video:audio ", audio)
		if videoAudio, ok := audio.(VideoAudio); ok {
			handler(videoAudio)
		}
	})

	return err
}
