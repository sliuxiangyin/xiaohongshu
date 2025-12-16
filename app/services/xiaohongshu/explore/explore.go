package explore

import (
	"errors"
	"fmt"
	"strconv"
	"xiaohongshu/app/services/xiaohongshu/entity"
	"xiaohongshu/app/services/xiaohongshu/scripts"

	"github.com/playwright-community/playwright-go"
)

type FeedsInfo struct {
	Element entity.Element
	Index   int
	Title   entity.Element
	Cover   entity.Element
	User    entity.Element
	Avatar  entity.Element
	Likes   entity.Element
}

type Explore struct {
	locator     playwright.Locator
	allFeeds    []FeedsInfo
	pageFeeds   []FeedsInfo
	elementInfo scripts.ElementInfo
}

// hasFeedIndex 检查指定的索引是否已存在于allFeeds中
func (s *Explore) hasFeedIndex(index int) bool {
	for _, feed := range s.allFeeds {
		if feed.Index == index {
			return true
		}
	}
	return false
}

func NewExplore(page playwright.Page) *Explore {
	locator := page.Locator("#exploreFeeds")
	elementInfo, err := scripts.GetElementInfo(locator)
	if err != nil {
		fmt.Println(fmt.Sprintf("GetElementInfo err: %v", err))
	}
	return &Explore{
		locator:     locator,
		elementInfo: elementInfo,
		allFeeds:    make([]FeedsInfo, 0),
		pageFeeds:   make([]FeedsInfo, 0),
	}
}
func (s *Explore) Show() ([]FeedsInfo, error) {
	var err error
	s.pageFeeds, err = s.getExploreFeeds()
	s.allFeeds = append(s.allFeeds, s.pageFeeds...)
	return s.pageFeeds, err
}

// NextPage 下一页
func (s *Explore) NextPage() {
	feedLen := len(s.allFeeds)
	if feedLen != 0 {
		_ = scripts.SmoothScrollTo(s.locator, s.elementInfo.Element.VisibleHeight)
	}
}

func (s *Explore) getExploreFeeds() ([]FeedsInfo, error) {
	var elements []FeedsInfo
	locator := s.locator.Locator("section")

	sectionElements, err := locator.All()
	if err != nil {
		return nil, fmt.Errorf("failed to find section elements: %v", err)
	}
	if len(sectionElements) == 0 {
		return elements, nil // 返回空数组而不是错误
	}
	// 遍历每个元素，提取相关信息
	for _, element := range sectionElements {
		// 检查元素是否在视窗内可见
		isVisible := scripts.GetElementIsVisible(element)
		if !isVisible {
			continue
		}

		dataIndex, err := element.GetAttribute("data-index")
		if err != nil || dataIndex == "" {
			continue
		}

		dataIndexInt, err := strconv.ParseInt(dataIndex, 10, 64)
		if err != nil {
			continue
		}

		//判断 dataIndexInt 是否已经存在
		if s.hasFeedIndex(int(dataIndexInt)) {
			continue
		}

		var e FeedsInfo
		e.Element = entity.Element{
			Text:     "",
			Selector: element,
		}
		e.Index = int(dataIndexInt)
		// 获取封面图片链接
		coverImgElement := element.Locator("a.cover img")

		coverImg, err := coverImgElement.GetAttribute("src")
		if err != nil || coverImg == "" {
			continue
		}
		e.Cover = entity.Element{
			Text:     coverImg,
			Selector: coverImgElement,
		}

		// 获取标题
		titleElement := element.Locator(".footer .title")
		titleText, err := titleElement.TextContent()
		if err != nil || titleText == "" {
			continue
		}
		e.Title = entity.Element{
			Text:     titleText,
			Selector: titleElement,
		}

		// 获取用户名称和头像
		authorElement := element.Locator(".author-wrapper .author")

		authorName, err := authorElement.TextContent()
		if err != nil || authorName == "" {
			continue
		}
		e.User = entity.Element{
			Text:     authorName,
			Selector: authorElement,
		}

		authorAvatar := authorElement.Locator("img")
		avatarSrc, err := authorAvatar.GetAttribute("src")
		if err != nil || avatarSrc == "" {
			continue
		}
		e.Avatar = entity.Element{
			Text:     avatarSrc,
			Selector: authorAvatar,
		}
		// 获取点赞数
		likeElement := element.Locator(".like-wrapper .count")

		likeCount, err := likeElement.TextContent()
		if err != nil {
			continue
		}
		e.Likes = entity.Element{
			Text:     likeCount,
			Selector: likeElement,
		}
		elements = append(elements, e)
	}

	return elements, nil
}

func (s *Explore) GetFeed(index int) (FeedsInfo, error) {
	for _, feed := range s.pageFeeds {
		if feed.Index == index {
			return feed, nil
		}
	}
	return FeedsInfo{}, errors.New("not found")
}

// 刷新页面
func (s *Explore) RefreshPage() error {
	selector := s.locator.Locator(".floating-btn-sets .reload")
	s.pageFeeds = make([]FeedsInfo, 0)
	s.allFeeds = make([]FeedsInfo, 0)
	return selector.Click()

}
