package scripts

import (
	"encoding/json"

	"github.com/playwright-community/playwright-go"
)

type ElementInfo struct {
	Scroll  ScrollInfo    `json:"scroll"`
	Element ElementDetail `json:"element"`
}

// ScrollInfo 描述滚动信息
type ScrollInfo struct {
	ScrollY float64 `json:"scrollY"`
	ScrollX float64 `json:"scrollX"`
}

// ElementDetail 描述元素详细信息
type ElementDetail struct {
	VisibleHeight            float64 `json:"visibleHeight"`
	ActualHeight             float64 `json:"actualHeight"`
	ContentHeight            float64 `json:"contentHeight"`
	ViewportTop              float64 `json:"viewportTop"`
	ViewportBottom           float64 `json:"viewportBottom"`
	AbsoluteTop              float64 `json:"absoluteTop"`
	AbsoluteBottom           float64 `json:"absoluteBottom"`
	IsFullyVisible           bool    `json:"isFullyVisible"`
	IsPartiallyVisible       bool    `json:"isPartiallyVisible"`
	VisibleRatio             float64 `json:"visibleRatio"`
	TopFromViewportTop       float64 `json:"topFromViewportTop"`
	BottomFromViewportBottom float64 `json:"bottomFromViewportBottom"`
}

const ToolJs = `

function smoothScrollTo(element, distance, duration = 1000) {
    const target = typeof element === 'string' ? document.querySelector(element) : element;
    if (!target) return;

    const start = window.pageYOffset || document.documentElement.scrollTop;
    const targetPosition = start + distance;
    let startTime = null;

    function animation(currentTime) {
        if (startTime === null) startTime = currentTime;
        const timeElapsed = currentTime - startTime;
        const progress = Math.min(timeElapsed / duration, 1);
        
        // 使用缓动函数
        const easeProgress = easeInOutCubic(progress);
        const currentPosition = start + distance * easeProgress;
        
        window.scrollTo(0, currentPosition);
        
        if (timeElapsed < duration) {
            requestAnimationFrame(animation);
        }
    }
    
    function easeInOutCubic(t) {
        return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
    }
    
    requestAnimationFrame(animation);
}


// 获取当前滚动距离
function getScrollDistance() {
    return {
        scrollY: window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop,
        scrollX: window.pageXOffset || document.documentElement.scrollLeft || document.body.scrollLeft
    };
}
// 获取元素可见高度（视口内可见部分）+ 滚动距离
function getVisibleHeight(element) {
    const rect = element.getBoundingClientRect();
    const windowHeight = window.innerHeight || document.documentElement.clientHeight;
    const scrollY = getScrollDistance().scrollY;

    const visibleTop = Math.max(rect.top, 0);
    const visibleBottom = Math.min(rect.bottom, windowHeight);
    const visibleHeight = Math.max(visibleBottom - visibleTop, 0);

    // 元素在文档中的实际位置（加上滚动距离）
    const absoluteTop = rect.top + scrollY;
    const absoluteBottom = rect.bottom + scrollY;

    return {
        visibleHeight: visibleHeight,
        visibleTop: visibleTop,
        visibleBottom: visibleBottom,
        absoluteTop: absoluteTop,
        absoluteBottom: absoluteBottom,
        scrollY: scrollY,
        isFullyVisible: rect.top >= 0 && rect.bottom <= windowHeight,
        isPartiallyVisible: rect.top < windowHeight && rect.bottom >= 0
    };
}
// 获取元素实际高度（包括padding、border）+ 滚动距离
function getActualHeight(element) {
    const scrollY = getScrollDistance().scrollY;
    const rect = element.getBoundingClientRect();

    return {
        actualHeight: element.offsetHeight,
        top: rect.top + scrollY,
        bottom: rect.bottom + scrollY,
        scrollY: scrollY
    };
}
// 获取元素完整信息
function getElementInfo(element) {
    const scrollInfo = getScrollDistance();
    const rect = element.getBoundingClientRect();
    const windowHeight = window.innerHeight;

    const visibleTop = Math.max(rect.top, 0);
    const visibleBottom = Math.min(rect.bottom, windowHeight);
    const visibleHeight = Math.max(visibleBottom - visibleTop, 0);

    return {
        scroll: scrollInfo,
        element: {
            // 高度信息
            visibleHeight: visibleHeight,
            actualHeight: element.offsetHeight,
            contentHeight: element.scrollHeight,

            // 相对视口位置
            viewportTop: rect.top,
            viewportBottom: rect.bottom,

            // 绝对位置（文档中的位置）
            absoluteTop: rect.top + scrollInfo.scrollY,
            absoluteBottom: rect.bottom + scrollInfo.scrollY,

            // 可见性状态
            isFullyVisible: rect.top >= 0 && rect.bottom <= windowHeight,
            isPartiallyVisible: rect.top < windowHeight && rect.bottom >= 0,
            visibleRatio: visibleHeight / element.offsetHeight,

            // 边界信息
            topFromViewportTop: rect.top,
            bottomFromViewportBottom: windowHeight - rect.bottom
        }
    };
}
window.smoothScrollTo=smoothScrollTo;
window.getScrollDistance=getScrollDistance;
window.getVisibleHeight=getVisibleHeight;
window.getActualHeight=getActualHeight;
window.getElementInfo=getElementInfo;



`

func GetElementInfo(selector playwright.Locator) (ElementInfo, error) {
	var element ElementInfo
	// 获取元素的信息并转化为结构体
	evaluate, err := selector.EvaluateHandle("(element) => {let info= getElementInfo(element); console.log(info);return info;}", nil)
	if err != nil {
		return element, err
	}
	value, err := evaluate.JSONValue()
	if err != nil {
		return ElementInfo{}, err
	}
	elementInfoJSON, err := json.Marshal(value)
	if err != nil {
		return element, err
	}
	err = json.Unmarshal(elementInfoJSON, &element)
	return element, err

}

func GetElementIsVisible(selector playwright.Locator) bool {
	isVisible, err := selector.EvaluateHandle("(element) => { const info = getElementInfo(element); return info.element.isPartiallyVisible;}", nil)
	if err != nil {
		return false
	}
	value, err := isVisible.JSONValue()
	if err != nil || !value.(bool) {
		return false
	}
	return true
}

// SmoothScrollTo 在浏览器中对指定元素执行平滑滚动操作
func SmoothScrollTo(element playwright.Locator, distance float64) error {
	_, err := element.EvaluateHandle(`(element,distance) => {
		window.smoothScrollTo(element, distance);
	}`, distance)
	return err
}

func GetVideoFrame(element playwright.Locator, bindFunc string) error {
	script := `() => {
  const video = document.querySelector('video');
  if (!video) return;

  const canvas = document.createElement('canvas');
  const ctx = canvas.getContext('2d');

  function capture() {
    if (video.readyState >= 2) {
      canvas.width = video.videoWidth;
      canvas.height = video.videoHeight;

      ctx.drawImage(video, 0, 0);
      const frame = ctx.getImageData(0, 0, canvas.width, canvas.height);

      // RGBA Uint8ClampedArray
      window.__onVideoFrame({
        width: frame.width,
        height: frame.height,
        data: Array.from(frame.data), // 或 base64
        ts: video.currentTime,
      });
    }
    requestAnimationFrame(capture);
  }
  capture();
}();
`

	_, err := element.EvaluateHandle(script, bindFunc)

	return err
}
