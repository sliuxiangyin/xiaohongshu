package explore

import (
	"strings"
	"xiaohongshu/app/services/xiaohongshu/entity"

	"github.com/playwright-community/playwright-go"
)

type ChannelInfo struct {
	entity.Element
	active bool // 元素是否具有"active"类名
}
type Channel struct {
	locator playwright.Locator
	info    []ChannelInfo
}

func NewChannel(page playwright.Page) *Channel {
	return &Channel{
		locator: page.Locator("#channel-container"),
		info:    make([]ChannelInfo, 0),
	}
}
func (c *Channel) Show() ([]ChannelInfo, error) {
	locator := c.locator.Locator(".content-container .channel")
	err := locator.WaitFor(playwright.LocatorWaitForOptions{Timeout: playwright.Float(3000)})
	if err != nil {
		return nil, err
	}
	elements, err := locator.All()
	if err != nil {
		return nil, err
	}
	result := make([]ChannelInfo, 0, len(elements))

	for _, element := range elements {
		// 提取元素的文本内容
		text, err := element.TextContent()
		if err != nil {
			continue
		}

		// 检查元素是否含有"active"类名
		className, err := element.GetAttribute("class")
		active := false
		if err == nil && className != "" {
			// 检查class属性中是否包含"active"类名
			active = strings.Contains(className, "active")
		}

		// 创建ChannelInfo实例并添加到结果中
		channelInfo := ChannelInfo{
			Element: entity.Element{
				Text:     text,
				Selector: element, // 直接使用当前element作为选择器
			},
			active: active,
		}
		result = append(result, channelInfo)
	}
	c.info = result
	return nil, err
}
