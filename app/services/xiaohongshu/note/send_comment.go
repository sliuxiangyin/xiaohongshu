package note

import "github.com/playwright-community/playwright-go"

// SendComment 发送评论
type SendComment struct {
	locator playwright.Locator
}

func NewSendComment(locator playwright.Locator) *SendComment {
	//engage-bar active
	return &SendComment{
		locator: locator,
	}
}

func (s *SendComment) Click() {
	s.locator.Locator("")
}
