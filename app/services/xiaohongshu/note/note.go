package note

import (
	"xiaohongshu/app/services/xiaohongshu/entity"
	"xiaohongshu/app/services/xiaohongshu/scripts"

	"github.com/playwright-community/playwright-go"
	"github.com/spf13/cast"
)

// Note 浏览详情
type SwiperItem struct {
	imag   string
	active bool
	index  int
}
type Swiper struct {
	item []SwiperItem
	entity.Element
}

type NoteInfo struct {
	noteType string
	//图片
	swiper     Swiper
	swiperPrev entity.Element
	swiperNext entity.Element
	//视频
	video *Video
	//发布作者
	authorElement entity.Element
	//关注按钮
	followElement entity.Element
	//点赞 .engage-bar .like-lottie
	likeElement entity.Element
	//收藏 .engage-bar .collect-wrapper
	collectElement entity.Element
	//内容
	title entity.Element
	//内容
	desc entity.Element
	//时间地址
	dateAddress  entity.Element
	commentCount entity.Element
	//评论
	comment *Comment
}

func (n NoteInfo) Video() *Video {
	return n.video
}

type Note struct {
	locator playwright.Locator
}

func NewNote(page playwright.Page) *Note {
	locator := page.Locator("#noteContainer")
	return &Note{locator: locator}
}

func (n *Note) Show() (NoteInfo, error) {
	var info NoteInfo

	attribute, err := n.locator.GetAttribute("data-type")
	if err != nil {
		return info, err
	}
	info.noteType = attribute
	if info.noteType == "video" {
		info.video = NewVideo(n.locator.Locator(".player-container "))

	} else {
		info.swiper = n.swiperHandler(n.locator.Locator(".swiper-slide:not(.swiper-slide-duplicate-active)"))
	}
	//用户头像区域
	info.authorElement = scripts.BuildElement(n.locator, ".interaction-container .author-container .author-wrapper .info .name")
	//关注区域
	info.followElement = scripts.BuildElement(n.locator, ".interaction-container .author-container .author-wrapper  .note-detail-follow-btn", "关注")
	info.likeElement = scripts.BuildElement(n.locator, ".engage-bar .like-lottie")
	info.collectElement = scripts.BuildElement(n.locator, ".engage-bar .collect-wrapper")

	//标题
	info.title = scripts.BuildElement(n.locator, ".interaction-container .note-scroller .note-content .title")

	info.desc = scripts.BuildElement(n.locator, ".interaction-container .note-scroller .note-content .desc")

	info.dateAddress = scripts.BuildElement(n.locator, ".interaction-container .note-scroller .note-content .bottom-container .date")

	info.commentCount = scripts.BuildElement(n.locator, ".interaction-container .note-scroller .comments-el .comments-container .total")
	info.comment = NewComment(n.locator)
	return info, nil
}

func (n *Note) swiperHandler(slideLocator playwright.Locator) Swiper {
	slides, err := slideLocator.All()
	if err != nil {
		slides = make([]playwright.Locator, 0)
	}
	items := make([]SwiperItem, len(slides))
	for _, slide := range slides {
		index, err := slide.GetAttribute("data-index")
		if err != nil {
			continue
		}
		imaSrc, err := slide.Locator("img").GetAttribute("src")
		if err != nil {
			continue
		}
		items = append(items, SwiperItem{
			imag:   imaSrc,
			active: false,
			index:  cast.ToInt(index),
		})
	}
	return Swiper{
		item: items,
		Element: entity.Element{
			Text:     "show swiper",
			Selector: slideLocator,
		},
	}
}
