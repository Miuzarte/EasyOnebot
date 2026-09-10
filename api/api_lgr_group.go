package api

import "strconv"

/*
DelGroupNotice 删除群公告

https://lagrange-onebot.apifox.cn/236975967e0

参数:

	groupId: 群 Uin
	noticeId: 公告 ID
*/
func DelGroupNotice(groupId int, noticeId string) *Request {
	return NewReq("_del_group_notice", map[string]any{
		"group_id":  groupId,
		"notice_id": noticeId,
	})
}

type validateDelGroupNotice struct {
	GroupId  int    `validate:"gt=0"`
	NoticeId string `validate:"required"`
}

// DelGroupNotice 删除群公告
func (c LgrCaller) DelGroupNotice(groupId int, noticeId string) error {
	err := validate.Struct(&validateDelGroupNotice{groupId, noticeId})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(DelGroupNotice(groupId, noticeId))
}

/*
GetAiRecord 获取群 Ai 语音

https://lagrange-onebot.apifox.cn/236976005e0

参数:

	character: 语音声色
	groupId: 群 Uin
	text: 语音文本
	chatType(可选): 语音类型 <enum: 1, 2> <默认值: 1>
*/
func GetAiRecord(character string, groupId int, text string, chatType ...int) *Request {
	if len(chatType) == 0 {
		return NewReq("get_ai_record", map[string]any{
			"character": character,
			"group_id":  groupId,
			"text":      text,
		})
	} else {
		return NewReq("get_ai_record", map[string]any{
			"character": character,
			"group_id":  groupId,
			"text":      text,
			"chat_type": chatType[0],
		})
	}
}

type validateGetAiRecord struct {
	Character string `validate:"required"`
	GroupId   int    `validate:"gt=0"`
	Text      string `validate:"required"`
	ChatType  []int  `validate:"omitempty,dive,oneof=1 2"`
}

type GetAiRecordResp = string // 语音 Url

// GetAiRecord 获取群 Ai 语音
func (c LgrCaller) GetAiRecord(character string, groupId int, text string, chatType ...int) (GetAiRecordResp, error) {
	err := validate.Struct(&validateGetAiRecord{character, groupId, text, chatType})
	if err != nil {
		return "", err
	}
	return DecodeResponseAssertion[string](c.PostReq(GetAiRecord(character, groupId, text, chatType...)))
}

// func GetGroupHonorInfo(groupId int, typ string) *Request
// std: [GetGroupHonorInfo] https://lagrange-onebot.apifox.cn/236976055e0

/*
GetGroupNotice 获取群公告

https://lagrange-onebot.apifox.cn/236976194e0

参数:

	groupId: 群 Uin
*/
func GetGroupNotice(groupId int) *Request {
	return NewReq("_get_group_notice", map[string]any{
		"group_id": groupId,
	})
}

type GetGroupNoticeResp []struct {
	Message     NoticeMsg `json:"message" mapstructure:"message"`           // 消息
	NoticeID    string    `json:"notice_id" mapstructure:"notice_id"`       // 公告 ID
	PublishTime int       `json:"publish_time" mapstructure:"publish_time"` // 发布时间
	SenderID    int       `json:"sender_id" mapstructure:"sender_id"`       // 发布者 Uin
}

type NoticeMsg struct {
	Images []NoticeImg `json:"images" mapstructure:"images"` // 图片信息列表
	Text   string      `json:"text" mapstructure:"text"`     // 文本
}

type NoticeImg struct {
	Height string `json:"height" mapstructure:"height"` // 图片高
	ID     string `json:"id" mapstructure:"id"`         // 图片 ID
	Width  string `json:"width" mapstructure:"width"`   // 图片宽
}

// GetGroupNotice 获取群公告
func (c LgrCaller) GetGroupNotice(groupId int) (GetGroupNoticeResp, error) {
	err := validate.idGt0(groupId)
	if err != nil {
		return nil, err
	}
	return DecodeResponseAssertion[GetGroupNoticeResp](c.PostReq(GetGroupNotice(groupId)))
}

// func SetGroupAdmin(groupId, userId int, enable bool) *Request
// std: [SetGroupAdmin] https://lagrange-onebot.apifox.cn/236978274e0
// func SetGroupBan(userId, groupId, duration int) *Request
// std: [SetGroupBan] https://lagrange-onebot.apifox.cn/236978290e0

/*
SetGroupBotStatus 设置群Bot发言状态

https://lagrange-onebot.apifox.cn/236980735e0

参数:

	groupId: 群 Uin
	botId: 机器人 Uin
	enable: 是否开启
*/
func SetGroupBotStatus(groupId, botId int, enable bool) *Request {
	return NewReq("set_group_bot_status", map[string]any{
		"group_id": groupId,
		"bot_id":   botId,
		"enable":   enable,
	})
}

type validateSetGroupBotStatus struct {
	GroupId int `validate:"gt=0"`
	BotId   int `validate:"gt=0"`
}

type SetGroupBotStatusResp = int // 机器人 ID < >=0 >

// SetGroupBotStatus 设置群Bot发言状态
func (c LgrCaller) SetGroupBotStatus(groupId, botId int, enable bool) (SetGroupBotStatusResp, error) {
	err := validate.Struct(&validateSetGroupBotStatus{groupId, botId})
	if err != nil {
		return 0, err
	}
	return DecodeResponseAssertion[int](c.PostReq(SetGroupBotStatus(groupId, botId, enable)))
}

/*
SendGroupBotCallback 调用群机器人回调

https://lagrange-onebot.apifox.cn/236980749e0

参数:

	groupId: 群 Uin
	botId: 机器人 Uin
	enable: 是否开启
	data1(可选): 数据1
	data2(可选): 数据2
	data3...?
*/
func SendGroupBotCallback(groupId, botId int, datas ...string) *Request {
	req := NewReq("send_group_bot_callback", map[string]any{
		"group_id": groupId,
		"bot_id":   botId,
	})
	for i, data := range datas {
		req.Params["data_"+strconv.Itoa(i+1)] = data
	}
	return req
}

type SendGroupBotCallbackResp = int // 机器人 Uin < >=0 >

type validateSendGroupBotCallback struct {
	GroupId int      `validate:"gt=0"`
	BotId   int      `validate:"gt=0"`
	Datas   []string `validate:"omitempty"`
}

// SendGroupBotCallback 调用群机器人回调
func (c LgrCaller) SendGroupBotCallback(groupId, botId int, datas ...string) (SetGroupBotStatusResp, error) {
	err := validate.Struct(&validateSendGroupBotCallback{groupId, botId, datas})
	if err != nil {
		return 0, err
	}
	return DecodeResponseAssertion[int](c.PostReq(SendGroupBotCallback(groupId, botId, datas...)))
}

// func SetGroupCard(userId, groupId int, card string) *Request
// std: [SetGroupCard] https://lagrange-onebot.apifox.cn/236980775e0
// func SetGroupKick(userId, groupId int, rejectAddRequest bool) *Request
// std: [SetGroupKick] https://lagrange-onebot.apifox.cn/236980790e0
// func SetGroupLeave(groupId int, isDismiss bool) *Request
// std: [SetGroupLeave] https://lagrange-onebot.apifox.cn/236980810e0

/*
SendGroupNotice 发送群公告

https://lagrange-onebot.apifox.cn/236980823e0

参数:

	groupId: 群 Uin
	content: 公告内容
	image: 公告 image 链接, 支持 http/https/file/base64
*/
func SendGroupNotice(groupId int, content string, image any) *Request {
	return NewReq("_send_group_notice", map[string]any{
		"group_id": groupId,
		"content":  content,
		"image":    image,
	})
}

// SendGroupNotice 发送群公告
type validateSendGroupNotice struct {
	GroupId int    `validate:"gt=0"`
	Content string `validate:"required"`
	Image   any    `validate:"required"`
}

type SendGroupNoticeResp = string // 公告 ID

func (c LgrCaller) SendGroupNotice(groupId int, content string, image any) (SendGroupNoticeResp, error) {
	err := validate.Struct(&validateSendGroupNotice{groupId, content, image})
	if err != nil {
		return "", err
	}
	return DecodeResponseAssertion[string](c.PostReq(SendGroupNotice(groupId, content, image)))
}

// func SetGroupName(groupId int, groupName string) *Request
// std: [SetGroupName] https://lagrange-onebot.apifox.cn/236980841e0

/*
SetGroupPortrait 设置群头像

https://lagrange-onebot.apifox.cn/236980850e0

参数:

	groupId: 群 Uin
	file: file 链接, 支持 http/https/file/base64
*/
func SetGroupPortrait(groupId int, file any) *Request {
	return NewReq("set_group_portrait", map[string]any{
		"group_id": groupId,
		"file":     file,
	})
}

// SetGroupPortrait 设置群头像
type validateSetGroupPortrait struct {
	GroupId int `validate:"gt=0"`
	File    any `validate:"required"`
}

func (c LgrCaller) SetGroupPortrait(groupId int, file any) error {
	err := validate.Struct(&validateSetGroupPortrait{groupId, file})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupPortrait(groupId, file))
}

/*
SetGroupReaction 表情回复操作

https://lagrange-onebot.apifox.cn/236981369e0

参数:

	groupId: 群 Uin
	messageId: 消息 ID
	code: 表情代码
	isAdd: 是否是添加
*/
func SetGroupReaction(groupId, messageId int, code string, isAdd bool) *Request {
	return NewReq("set_group_reaction", map[string]any{
		"group_id":   groupId,
		"message_id": messageId,
		"code":       code,
		"is_add":     isAdd,
	})
}

// SetGroupReaction 表情回复操作
type validateSetGroupReaction struct {
	GroupId   int    `validate:"gt=0"`
	MessageId int    `validate:"ne=0"`
	Code      string `validate:"required"`
}

func (c LgrCaller) SetGroupReaction(groupId, messageId int, code string, isAdd bool) error {
	err := validate.Struct(&validateSetGroupReaction{groupId, messageId, code})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupReaction(groupId, messageId, code, isAdd))
}

// func SetGroupSpecialTitle(groupId, userId int, specialTitle string, duration int) *Request
// std: [SetGroupSpecialTitle] https://lagrange-onebot.apifox.cn/236981401e0
// func SetGroupWholeBan(groupId int, enable string) *Request
// std: [SetGroupWholeBan] https://lagrange-onebot.apifox.cn/236981414e0
