package entity

import "github.com/playwright-community/playwright-go"

// Element ElementInfo 描述DOM元素信息
type Element struct {
	Text     string             // 元素的文本内容
	Selector playwright.Locator // 元素的选择器
}

func (e *Element) Click() error {
	err := e.Selector.Click()
	return err
}
