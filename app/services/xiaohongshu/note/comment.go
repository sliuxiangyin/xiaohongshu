package note

import (
	"xiaohongshu/app/services/xiaohongshu/entity"
	"xiaohongshu/app/services/xiaohongshu/scripts"

	"github.com/playwright-community/playwright-go"
)

type CommentInfo struct {
	author      entity.Element
	content     entity.Element
	imgs        []entity.Element
	dateAddress entity.Element
	like        entity.Element
	reply       entity.Element
	subComment  []CommentInfo
	showMore    entity.Element
}

type Comment struct {
	locator playwright.Locator // 元素的选择器
}

func NewComment(locator playwright.Locator) *Comment {
	return &Comment{
		locator: locator,
	}
}
func (c Comment) Show() ([]CommentInfo, error) {
	// 获取所有父评论
	parentCommentsLocator := c.locator.Locator(".comments-el .comments-container .list-container .parent-comment")
	parentComments, err := parentCommentsLocator.All()
	if err != nil {
		return nil, err
	}

	comments := make([]CommentInfo, len(parentComments))
	for i, commentLocator := range parentComments {
		comment := CommentInfo{}
		// 用户信息名称
		comment.author = scripts.BuildElement(commentLocator, ".right .author-wrapper .author")
		// 内容
		comment.content = scripts.BuildElement(commentLocator, ".right .content")
		// 评论图片
		imgLocators, err := commentLocator.Locator(".right .comment-picture img").All()
		if err == nil {
			imgElements := make([]entity.Element, len(imgLocators))
			for j, imgLocator := range imgLocators {
				text, _ := imgLocator.GetAttribute("src")
				imgElements[j] = entity.Element{
					Text:     text,
					Selector: imgLocator,
				}
			}
			comment.imgs = imgElements
		}
		// 评论日期和地址
		comment.dateAddress = scripts.BuildElement(commentLocator, ".right .info .date")
		// 点赞(数量)
		comment.like = scripts.BuildElement(commentLocator, ".right .info .interactions .like-wrapper")
		// 回复(数量)
		comment.reply = scripts.BuildElement(commentLocator, ".right .info .interactions .reply")
		// 子评论处理
		subCommentsLocator := commentLocator.Locator(".reply-container .list-container .comment-item-sub")
		subComments, err := subCommentsLocator.All()
		if err == nil {
			subCommentElements := make([]CommentInfo, len(subComments))
			for j, subCommentLocator := range subComments {
				subComment := CommentInfo{}

				// 用户信息名称
				subComment.author = scripts.BuildElement(subCommentLocator, ".right .author-wrapper .author")

				// 内容
				subComment.content = scripts.BuildElement(subCommentLocator, ".right .content")
				// 评论图片
				subImgLocators, err := subCommentLocator.Locator(".right .comment-picture img").All()
				if err == nil {
					subImgElements := make([]entity.Element, len(subImgLocators))
					for k, subImgLocator := range subImgLocators {
						text, _ := subImgLocator.GetAttribute("src")
						subImgElements[k] = entity.Element{
							Text:     text,
							Selector: subImgLocator,
						}
					}
					subComment.imgs = subImgElements
				}
				// 评论日期和地址
				subComment.dateAddress = scripts.BuildElement(subCommentLocator, ".right .info .date")

				// 点赞(数量)
				subComment.like = scripts.BuildElement(subCommentLocator, ".right .info .interactions .like-wrapper")

				// 回复(数量)
				subComment.reply = scripts.BuildElement(subCommentLocator, ".right .info .interactions .reply")

				subCommentElements[j] = subComment
			}
			comment.subComment = subCommentElements
		}
		// 显示更多 展开
		comment.showMore = scripts.BuildElement(commentLocator, ".reply-container .show-more")
		comments[i] = comment
	}

	return comments, nil
}

//#noteContainer
// attr data-type : default  video

// 左侧banner .media-container
// 获取banner  document.querySelectorAll(".swiper-slide:not(.swiper-slide-duplicate)")
// 获取bannerImage  .swiper-slide:not(.swiper-slide-duplicate) .note-slider-img  img

//视频类型 .video-player-media

// 右侧内容区域 .interaction-container
// 关注：.interaction-container .author-container .author-wrapper  .note-detail-follow-btn
// 发布用户信息：.interaction-container .author-container .author-wrapper  .author-container
// 标题：.interaction-container .note-scroller .note-content .title
// 介绍：.interaction-container .note-scroller .note-content .desc
// 日期和地址：.interaction-container .note-scroller .note-content .bottom-container .date
// 举报按钮：.interaction-container .note-scroller .note-content .bottom-container .notedetail-menu

//评论区域 .interaction-container .note-scroller .comments-el .comments-container
//评论数量 .interaction-container .note-scroller .comments-el .comments-container .total

//评论内容 .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment
//用户信息名称 .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment  .right .author-wrapper .author
//内容 .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment  .right  .content
//评论图片 .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment  .right  .comment-picture img
//评论图片 .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment  .right  .comment-picture img
//评论日期和地址 .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment  .right  .info .date
//点赞(数量) .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment  .right  .info .interactions .like-wrapper
//回复(数量) .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment  .right  .info .interactions .reply

//子评论列表 .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment  .reply-container  .list-container .comment-item-sub
//显示更多 展开 .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment  .reply-container  .show-more

//用户信息名称 .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment   .reply-container .list-container .comment-item-sub .right .author-wrapper .author
//内容 .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment  .reply-container .list-container .comment-item-sub  .right  .content
//评论图片 .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment  .reply-container .list-container .comment-item-sub  .right  .comment-picture img
//评论日期和地址 .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment  .reply-container .list-container .comment-item-sub  .right  .info .date
//点赞(数量) .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment   .reply-container .list-container .comment-item-sub  .right  .info .interactions .like-wrapper
//回复(数量) .interaction-container .note-scroller .comments-el .comments-container .list-container .parent-comment  .reply-container .list-container .comment-item-sub   .right  .info .interactions .reply
