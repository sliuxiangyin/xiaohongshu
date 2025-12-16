package scripts

import (
	"xiaohongshu/app/services/xiaohongshu/entity"

	"github.com/playwright-community/playwright-go"
)

func BuildElement(locator playwright.Locator, selector string, ars ...string) entity.Element {
	loc := locator.Locator(selector)
	count, err := loc.Count()
	if err != nil || count == 0 {
		return entity.Element{}
	}
	text := ""
	if len(ars) == 1 {
		text = ars[0]
	} else {
		text, _ = loc.TextContent()
	}
	return entity.Element{
		Text:     text,
		Selector: loc,
	}
}
