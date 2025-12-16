package entities

import (
	"encoding/json"
)

// FlexibleUserInfo 定义灵活的用户信息结构体，可以处理缺失的字段
type FlexibleUserInfo struct {
	RedId    *string `json:"red_id,omitempty"`
	Nickname *string `json:"nickname,omitempty"`
	Desc     *string `json:"desc,omitempty"`
	Gender   *int    `json:"gender,omitempty"`
	Images   *string `json:"images,omitempty"`
	Imageb   *string `json:"imageb,omitempty"`
	UserId   *string `json:"user_id,omitempty"`
	Guest    *bool   `json:"guest,omitempty"`
}

// ToUserInfo 将FlexibleUserInfo转换为标准的UserInfo
func (f *FlexibleUserInfo) ToUserInfo() UserInfo {
	userInfo := UserInfo{}

	if f.RedId != nil {
		userInfo.RedId = *f.RedId
	}

	if f.Nickname != nil {
		userInfo.Nickname = *f.Nickname
	}

	if f.Desc != nil {
		userInfo.Desc = *f.Desc
	}

	if f.Gender != nil {
		userInfo.Gender = *f.Gender
	}

	if f.Images != nil {
		userInfo.Images = *f.Images
	}

	if f.Imageb != nil {
		userInfo.Imageb = *f.Imageb
	}

	if f.UserId != nil {
		userInfo.UserId = *f.UserId
	}

	if f.Guest != nil {
		userInfo.Guest = *f.Guest
	}

	return userInfo
}

// FlexibleApiResponse 定义灵活的API响应结构体
type FlexibleApiResponse struct {
	Code    int              `json:"code"`
	Success bool             `json:"success"`
	Msg     string           `json:"msg"`
	Data    FlexibleUserInfo `json:"data"`
}

// ParseFlexibleApiResponse 从字节数据解析灵活的API响应
func ParseFlexibleApiResponse(data []byte) (*FlexibleApiResponse, error) {
	var response FlexibleApiResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
