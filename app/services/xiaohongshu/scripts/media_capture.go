package scripts

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
	"github.com/spf13/cast"
)

// MediaCaptureScript 媒体捕获控制器的JavaScript代码
const MediaCaptureScript = `
(() => {
  if (window.__MediaCaptureController) return;
  class MediaCaptureController {
    constructor() {
      this.videoMap = new Map(); // key: videoEl, value: { audioCtx, src, processor, canvas, ctx, rafId }
    }

    async start(videoEl = null) {
      if (!videoEl) videoEl = document.querySelector('video');
      if (!videoEl) throw new Error('no video element');
      if (this.videoMap.has(videoEl)) return; // 已经在采集

      const video = videoEl;

      // ================== AUDIO ==================
      video.muted = false;
      video.volume = 1.0;
      video.crossOrigin = 'anonymous';

      const audioCtx = new AudioContext({ sampleRate: 48000 });
      await audioCtx.resume();

      const src = audioCtx.createMediaElementSource(video);
      const processor = audioCtx.createScriptProcessor(4096, 2, 2);

      src.connect(processor);
      processor.connect(audioCtx.destination);

      processor.onaudioprocess = e => {
        if (!this.videoMap.has(video)) return;
        if (video.paused) return;

        const input = e.inputBuffer;
        const frames = input.length;
        const ch0 = input.getChannelData(0);

        // energy gate
        let energy = 0;
        for (let i = 0; i < ch0.length; i++) {
          energy += Math.abs(ch0[i]);
          if (energy > 0.001) break;
        }
        if (energy <= 0.001) return;

        const pcm = new Float32Array(frames);
        pcm.set(ch0);

        window?.__onVideoAudio?.({
          sampleRate: audioCtx.sampleRate,
          channels: 1,
          frames,
          buffer: pcm.buffer,
          ts: audioCtx.currentTime,
        });
      };

      // ================== VIDEO ==================
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');
      let rafId = null;

      const captureFrame = () => {
        if (!this.videoMap.has(video)) return;

        if (video.readyState >= 2 && !video.paused) {
          const w = video.videoWidth;
          const h = video.videoHeight;

          if (w && h) {
            canvas.width = w;
            canvas.height = h;

            ctx.drawImage(video, 0, 0, w, h);
            const img = ctx.getImageData(0, 0, w, h);
		    console.log(img);
            window?.__onVideoFrame?.({
              width: w,
              height: h,
              data: img.data.buffer,
              ts: video.currentTime,
            });
          }
        }

        rafId = requestAnimationFrame(captureFrame);
      };

      captureFrame();

      // 保存状态
      this.videoMap.set(video, { audioCtx, src, processor, canvas, ctx, rafId });
    }

    stop(videoEl) {
      if (!videoEl) return;
      const state = this.videoMap.get(videoEl);
      if (!state) return;

      const { audioCtx, src, processor, rafId } = state;

      if (rafId) cancelAnimationFrame(rafId);

      if (processor) {
        processor.disconnect();
        processor.onaudioprocess = null;
      }

      if (src) src.disconnect();
      if (audioCtx && audioCtx.state !== 'closed') audioCtx.suspend();

      this.videoMap.delete(videoEl);
    }

    destroy(videoEl) {
      if (!videoEl) return;

      const state = this.videoMap.get(videoEl);
      if (!state) return;

      const { audioCtx, src, processor, rafId } = state;

      if (rafId) cancelAnimationFrame(rafId);

      if (processor) {
        processor.disconnect();
        processor.onaudioprocess = null;
      }

      if (src) src.disconnect();
      if (audioCtx && audioCtx.state !== 'closed') audioCtx.close();

      this.videoMap.delete(videoEl);
    }

    stopAll() {
      for (const videoEl of this.videoMap.keys()) {
        this.stop(videoEl);
      }
    }

    destroyAll() {
      for (const videoEl of Array.from(this.videoMap.keys())) {
        this.destroy(videoEl);
      }
    }
  }

  window.__MediaCaptureController = new MediaCaptureController();
})();
`

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
		console.log("调用MediaCaptureController.start");
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
func MediaListenVideoState(element playwright.Locator, handler func(bool)) error {
	page, err := element.Page()
	if err != nil {
		return err
	}
	// 暴露回调函数给浏览器
	err = page.ExposeFunction("__onVideoStateChange", func(args ...interface{}) interface{} {
		if len(args) > 0 {
			if isPlaying, ok := args[0].(bool); ok {
				handler(isPlaying)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	_, err = element.EvaluateHandle(`
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
func InjectMediaCaptureScript(page playwright.Page) error {
	// 首先获取页面实例
	var err error
	// 暴露__onVideoFrame回调函数到浏览器环境
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
	page.Video()
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
				fmt.Println(fmt.Sprintf("__onVideoAudio: %v", audio))
				fmt.Println("media:video:audio", GetEventBus().HasCallback("media:video:audio"))
				// 通过事件总线发送音频数据
				GetEventBus().Publish("media:video:audio", audio)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	// 注入媒体捕获脚本
	script := MediaCaptureScript
	err = page.AddInitScript(playwright.Script{
		Content: &script,
	})
	return err
}
