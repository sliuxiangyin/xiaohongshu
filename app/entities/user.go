package entities

// UserInfo 定义用户信息结构体
type UserInfo struct {
	RedId    string `json:"red_id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Desc     string `json:"desc,omitempty"`
	Gender   int    `json:"gender,omitempty"`
	Images   string `json:"images,omitempty"`
	Imageb   string `json:"imageb,omitempty"`
	UserId   string `json:"user_id,omitempty"`
	Guest    bool   `json:"guest,omitempty"`
}

// ApiResponse 定义API响应结构体
type ApiResponse struct {
	Code    int      `json:"code"`
	Success bool     `json:"success"`
	Msg     string   `json:"msg"`
	Data    UserInfo `json:"data"`
}
