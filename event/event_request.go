package event

// https://github.com/botuniverse/onebot-11/blob/master/event/request.md

// RequestBase 请求事件
type RequestBase struct {
	Base // [TYPE_L1_REQUEST] "request"

	// "friend", "group" // 请求类型
	RequestType string `json:"request_type" mapstructure:"request_type"`
	// 发送请求的 QQ 号
	UserId int `json:"user_id" mapstructure:"user_id"`
	// 验证信息
	Comment string `json:"comment" mapstructure:"comment"`
	// 请求 flag，在调用处理请求的 API 时需要传入
	Flag string `json:"flag" mapstructure:"flag"`
}

// RequestFriend 加好友请求
type RequestFriend struct {
	RequestBase // [TYPE_L2_REQUEST_FRIEND] "friend"
}

// RequestGroup 加群请求／邀请
type RequestGroup struct {
	RequestBase // [TYPE_L2_REQUEST_GROUP] "group"

	// "add", "invite" // 请求子类型，分别表示加群请求、邀请登录号入群
	SubType string `json:"sub_type" mapstructure:"sub_type"`
	// 群号
	GroupId int `json:"group_id" mapstructure:"group_id"`

	RequestGroup_Nc
}

type RequestGroup_Nc struct {
	InvitorId int `json:"invitor_id" mapstructure:"invitor_id"`
}
