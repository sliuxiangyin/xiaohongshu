package scripts

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/playwright-community/playwright-go"
)

// ClassDOMObserver 结构体用于管理浏览器环境中指定类名元素的监听功能
type ClassDOMObserver struct {
	page          playwright.Page
	observers     map[string]string // observerId -> className
	mutex         sync.RWMutex
	className     string
	safeClassName string
	onAdd         func(string)
	onRemove      func(string)
}

// sanitizeClassName 将CSS类名转换为合法的JavaScript标识符
func sanitizeClassName(className string) string {
	// 将连字符和点替换为下划线
	sanitized := regexp.MustCompile(`[-.]`).ReplaceAllString(className, "_")
	// 移除非法字符，只保留字母、数字和下划线
	sanitized = regexp.MustCompile(`[^a-zA-Z0-9_]`).ReplaceAllString(sanitized, "")
	// 确保不以数字开头
	if len(sanitized) > 0 && regexp.MustCompile(`^[0-9]`).MatchString(sanitized) {
		sanitized = "_" + sanitized
	}
	// 如果结果为空，则使用默认名称
	if sanitized == "" {
		sanitized = "class"
	}
	return sanitized
}

// NewClassDOMObserver 创建一个新的ClassDOMObserver实例
func NewClassDOMObserver(page playwright.Page) *ClassDOMObserver {
	return &ClassDOMObserver{
		page:      page,
		observers: make(map[string]string),
	}
}

func (cdo *ClassDOMObserver) Start(className string) error {
	cdo.mutex.Lock()
	defer cdo.mutex.Unlock()
	cdo.className = className
	// 对类名进行安全处理，确保生成合法的JavaScript标识符
	cdo.safeClassName = sanitizeClassName(className)

	// 暴露添加回调函数到浏览器环境
	err := cdo.page.ExposeFunction("__onClassElementAdded_"+cdo.safeClassName, func(args ...interface{}) interface{} {
		if len(args) > 0 {
			if elementId, ok := args[0].(string); ok {
				if cdo.onAdd != nil {
					cdo.onAdd(elementId)
				}
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("暴露添加回调函数失败: %w", err)
	}

	// 暴露移除回调函数到浏览器环境
	err = cdo.page.ExposeFunction("__onClassElementRemoved_"+cdo.safeClassName, func(args ...interface{}) interface{} {
		if len(args) > 0 {
			if elementId, ok := args[0].(string); ok {
				if cdo.onRemove != nil {
					cdo.onRemove(elementId)
				}
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("暴露移除回调函数失败: %w", err)
	}
	return nil
}

// Observe 启动对某个特定CSS类名元素的监听
// className: 要监听的CSS类名
// onAdd: 元素被添加时的回调函数
// onRemove: 元素被移除时的回调函数
func (cdo *ClassDOMObserver) Observe() (string, error) {
	cdo.mutex.Lock()
	defer cdo.mutex.Unlock()
	// 创建观察器实例（假设MinimalClassDOMObserver已在其他地方注入）
	script := fmt.Sprintf(`
		(function() {
			if (typeof MinimalClassDOMObserver !== 'undefined') {
				// 为每个观察器创建一个唯一的ID
				const observerId = 'observer_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
				// 创建观察器实例
				const observer = new MinimalClassDOMObserver({
					className: '%s',
					onAdd: function(element) {
						// 为元素分配唯一ID（如果还没有的话）
						if (!element.id) {
							element.id = 'observed-' + Date.now() + '-' + Math.random().toString(36).substr(2, 9);
						}
						console.log("element.id",window.__onClassElementAdded)
						// 调用Go端的回调函数，传递元素ID
						window.__onClassElementAdded_%s(element.id);
					},
					onRemove: function(element) {
					console.log("element.id",window.__onClassElementRemoved)
						// 调用Go端的回调函数，传递元素ID
						const elementId = element.id || 'unknown';
						window.__onClassElementRemoved_%s(elementId);
					}
				});

				// 启动观察器
				observer.start();

				// 将观察器存储在全局变量中以便后续操作
				if (!window.__classDOMObservers) {
					window.__classDOMObservers = {};
				}
				window.__classDOMObservers[observerId] = observer;

				return observerId;
			} else {
				console.error('MinimalClassDOMObserver is not defined');
				return null;
			}
		})();
	`, cdo.className, cdo.safeClassName, cdo.safeClassName)

	result, err := cdo.page.Evaluate(script)
	if err != nil {
		return "", fmt.Errorf("执行观察脚本失败: %w", err)
	}
	if result == nil {
		return "", fmt.Errorf("无法创建观察器，可能是因为MinimalClassDOMObserver未定义")
	}
	observerId, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("观察器ID不是字符串类型")
	}
	// 存储观察器信息
	cdo.observers[observerId] = cdo.className
	return observerId, nil
}

func (cdo *ClassDOMObserver) OnAdd(callback func(string2 string)) error {
	cdo.onAdd = callback
	return nil
}
func (cdo *ClassDOMObserver) OnRemove(callback func(string2 string)) error {
	cdo.onRemove = callback
	return nil
}

// Unobserve 停止并移除对某个已存在监听任务的观察器
func (cdo *ClassDOMObserver) Unobserve(observerId string) error {
	cdo.mutex.Lock()
	defer cdo.mutex.Unlock()

	_, exists := cdo.observers[observerId]
	if !exists {
		return fmt.Errorf("未找到观察器ID: %s", observerId)
	}

	// 执行JavaScript代码停止观察器
	script := fmt.Sprintf(`
		(function() {
			if (typeof __classDOMObservers !== 'undefined' && __classDOMObservers['%s']) {
				// 销毁观察器
				__classDOMObservers['%s'].destroy();
				// 从全局变量中移除
				delete __classDOMObservers['%s'];
				return true;
			} else {
				console.error('Observer not found: %s');
				return false;
			}
		})();
	`, observerId, observerId, observerId, observerId)

	_, err := cdo.page.Evaluate(script)
	if err != nil {
		return fmt.Errorf("执行停止观察脚本失败: %w", err)
	}

	// 清理Go端的资源
	delete(cdo.observers, observerId)

	// 移除暴露的函数

	return nil
}

// UnobserveAll 停止所有观察器
func (cdo *ClassDOMObserver) UnobserveAll() error {
	cdo.mutex.Lock()
	defer cdo.mutex.Unlock()
	// 执行JavaScript代码停止所有观察器
	script := `
		(function() {
			if (typeof __classDOMObservers !== 'undefined') {
				// 销毁所有观察器
				for (const observerId in __classDOMObservers) {
					__classDOMObservers[observerId].destroy();
					delete __classDOMObservers[observerId];
				}
				return true;
			} else {
				console.error('__classDOMObservers is not defined');
				return false;
			}
		})();
	`

	_, err := cdo.page.Evaluate(script)
	if err != nil {
		return fmt.Errorf("执行停止所有观察脚本失败: %w", err)
	}

	// 清理Go端的所有资源
	for observerId, _ := range cdo.observers {
		// 对类名进行安全处理，确保生成合法的JavaScript标识符
		//safeClassName := sanitizeClassName(className)
		// 移除暴露的函数
		//cdo.page.Evaluate(fmt.Sprintf("delete window.__onClassElementAdded_%s", safeClassName), nil)
		//cdo.page.Evaluate(fmt.Sprintf("delete window.__onClassElementRemoved_%s", safeClassName), nil)
		delete(cdo.observers, observerId)
	}

	return nil
}

// CheckState 检查指定观察器的状态
func (cdo *ClassDOMObserver) CheckState(observerId string) (bool, error) {
	cdo.mutex.RLock()
	defer cdo.mutex.RUnlock()

	_, exists := cdo.observers[observerId]
	if !exists {
		return false, fmt.Errorf("未找到观察器ID: %s", observerId)
	}

	// 执行JavaScript代码检查状态
	script := fmt.Sprintf(`
		(function() {
			if (typeof __classDOMObservers !== 'undefined' && __classDOMObservers['%s']) {
				// 检查观察器是否存在且未被销毁
				return !!__classDOMObservers['%s'];
			} else {
				return false;
			}
		})();
	`, observerId, observerId)

	result, err := cdo.page.Evaluate(script)
	if err != nil {
		return false, fmt.Errorf("执行状态检查脚本失败: %w", err)
	}

	state, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("状态检查结果不是布尔类型")
	}

	return state, nil
}

// GetActiveObservers 获取所有活动的观察器
func (cdo *ClassDOMObserver) GetActiveObservers() ([]string, error) {
	// 执行JavaScript代码获取所有活动观察器
	script := `
		(function() {
			if (typeof __classDOMObservers !== 'undefined') {
				return Object.keys(__classDOMObservers);
			} else {
				return [];
			}
		})();
	`

	result, err := cdo.page.Evaluate(script)
	if err != nil {
		return nil, fmt.Errorf("执行获取活动观察器脚本失败: %w", err)
	}

	observers, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("活动观察器列表格式不正确")
	}

	// 转换为字符串切片
	activeObservers := make([]string, len(observers))
	for i, obs := range observers {
		if obsStr, ok := obs.(string); ok {
			activeObservers[i] = obsStr
		}
	}

	return activeObservers, nil
}

// CreateClassDOMObserver 创建并返回一个新的ClassDOMObserver实例的便捷函数
func CreateClassDOMObserver(page playwright.Page) *ClassDOMObserver {
	return NewClassDOMObserver(page)
}
