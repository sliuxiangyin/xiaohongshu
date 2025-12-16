package scripts

import (
	"encoding/json"
	"fmt"

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


class MinimalClassDOMObserver {
  constructor(className) {
    this.className = className;
    this.onAdd = [];
    this.onRemove = [];
    
    this.observer = new MutationObserver((mutations) => {
      for (const mutation of mutations) {
        // 处理新增
        for (const node of mutation.addedNodes) {
          if (node.nodeType === 1) this.processElement(node, 'add');
        }
        // 处理移除
        for (const node of mutation.removedNodes) {
          if (node.nodeType === 1) this.processElement(node, 'remove');
        }
      }
    });
    
    this.observer.observe(document.body, {
      childList: true,
      subtree: true
    });
  }
  
  processElement(element, action) {
    // 检查元素本身
    if (element.classList && element.classList.contains(this.className)) {
      this.trigger(element, action);
    }
    
    // 检查后代
    if (element.querySelectorAll) {
      const children = element.querySelectorAll('.' + this.className);
      for (const child of children) {
        if (child !== element) { // 避免重复触发
          this.trigger(child, action);
        }
      }
    }
  }
  
  trigger(element, action) {
    const callbacks = action === 'add' ? this.onAdd : this.onRemove;
    for (const callback of callbacks) {
      callback(element);
    }
  }
  
  added(callback) {
    this.onAdd.push(callback);
    return this;
  }
  
  removed(callback) {
    this.onRemove.push(callback);
    return this;
  }
  
  stop() {
    this.observer.disconnect();
  }
}
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

// WatchClassElements 监控指定类名的div元素的添加和移除事件
func WatchClassElements(page playwright.Page, className string, onAddCallback, onRemoveCallback func(string) error) error {
	// 暴露添加回调函数到浏览器环境
	err := page.ExposeFunction("__onElementAdded", func(args ...interface{}) interface{} {
		if len(args) > 0 {
			if elementId, ok := args[0].(string); ok {
				if onAddCallback != nil {
					onAddCallback(elementId)
				}
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("暴露添加回调函数失败: %w", err)
	}

	// 暴露移除回调函数到浏览器环境
	err = page.ExposeFunction("__onElementRemoved", func(args ...interface{}) interface{} {
		if len(args) > 0 {
			if elementId, ok := args[0].(string); ok {
				if onRemoveCallback != nil {
					onRemoveCallback(elementId)
				}
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("暴露移除回调函数失败: %w", err)
	}

	// 在页面中执行监控脚本
	script := fmt.Sprintf(`
		(() => {
			// 确保MinimalClassDOMObserver已定义
			if (typeof MinimalClassDOMObserver !== 'undefined') {
				// 创建观察器实例
				const observer = new MinimalClassDOMObserver('%s');
				
				// 注册添加回调
				observer.added((element) => {
					// 为元素分配唯一ID（如果还没有的话）
					if (!element.id) {
						element.id = 'observed-' + Date.now() + '-' + Math.random().toString(36).substr(2, 9);
					}
					// 调用Go端的回调函数
					window.__onElementAdded(element.id);
				});
				
				// 注册移除回调
				observer.removed((element) => {
					// 调用Go端的回调函数
					window.__onElementRemoved(element.id || 'unknown');
				});
			} else {
				console.error('MinimalClassDOMObserver is not defined');
			}
		})();
	`, className)

	_, err = page.EvaluateHandle(script, nil)
	if err != nil {
		return fmt.Errorf("执行监控脚本失败: %w", err)
	}

	return nil
}

// // ExampleWatchClassElements 使用示例：监控特定类名元素的添加和移除
// func ExampleWatchClassElements(page playwright.Page) error {
// 	// 定义元素添加时的回调函数
// 	onAdd := func(elementId string) error {
// 		fmt.Printf("元素已添加，ID: %s\n", elementId)
// 		// 在这里可以添加您需要的处理逻辑
// 		return nil
// 	}

// 	// 定义元素移除时的回调函数
// 	onRemove := func(elementId string) error {
// 		fmt.Printf("元素已移除，ID: %s\n", elementId)
// 		// 在这里可以添加您需要的处理逻辑
// 		return nil
// 	}

// 	// 开始监控 "note-item" 类名的元素
// 	err := WatchClassElements(page, "note-item", onAdd, onRemove)
// 	if err != nil {
// 		return fmt.Errorf("设置元素监控失败: %w", err)
// 	}

// 	return nil
// }
