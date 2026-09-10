package api

import (
	"github.com/Miuzarte/EasyOnebot/event"
	"github.com/Miuzarte/EasyOnebot/message"
)

/*
DeleteEssenceMsg 删除精华消息

https://lagrange-onebot.apifox.cn/236981678e0

参数:

	messageId: 消息 ID
*/
func DeleteEssenceMsg(messageId int) *Request {
	return NewReq("delete_essence_msg", map[string]any{
		"message_id": messageId,
	})
}

// FriendPoke 私聊戳一戳
func (c LgrCaller) DeleteEssenceMsg(messageId int) error {
	err := validate.msgIdNe0(messageId)
	if err != nil {
		return err
	}
	return c.PostReqNoResp(DeleteEssenceMsg(messageId))
}

// func DeleteMsg(messageId int) *Request
// std: [DeleteMsg] https://lagrange-onebot.apifox.cn/236981681e0

/*
FriendPoke 私聊戳一戳

https://lagrange-onebot.apifox.cn/236981714e0

参数:

	userId: 用户 Uin
*/
func FriendPoke(userId int) *Request {
	return NewReq("friend_poke", map[string]any{
		"user_id": userId,
	})
}

// FriendPoke 私聊戳一戳
func (c LgrCaller) FriendPoke(userId int) error {
	err := validate.idGt0(userId)
	if err != nil {
		return err
	}
	return c.PostReqNoResp(FriendPoke(userId))
}

/*
GetEssenceMsgList 获取精华消息列表

https://lagrange-onebot.apifox.cn/236981735e0

参数:

	groupId: 群 Uin
*/
func GetEssenceMsgList(groupId int) *Request {
	return NewReq("get_essence_msg_list", map[string]any{
		"group_id": groupId,
	})
}

type GetEssenceMsgListResp = map[string]any // too complex to parse

// FriendPoke 私聊戳一戳
func (c LgrCaller) GetEssenceMsgList(groupId int) (GetEssenceMsgListResp, error) {
	err := validate.idGt0(groupId)
	if err != nil {
		return nil, err
	}
	return DecodeResponseAssertion[map[string]any](c.PostReq(GetEssenceMsgList(groupId)))
}

// func GetForwardMsg(id string) *Request
// std: [GetForwardMsg] https://lagrange-onebot.apifox.cn/236981740e0

/*
GetFriendMsgHistory 获取好友历史聊天记录

https://lagrange-onebot.apifox.cn/236981743e0

参数:

	userId: 用户 Uin
	messageId: 消息 ID
	count(可选): 消息数量 <默认值: 20>
*/
func GetFriendMsgHistory(userId, messageId int, count ...int) *Request {
	if len(count) == 0 {
		return NewReq("get_friend_msg_history", map[string]any{
			"user_id":    userId,
			"message_id": messageId,
		})
	} else {
		return NewReq("get_friend_msg_history", map[string]any{
			"user_id":     userId,
			"message_id":  messageId,
			"message_num": count[0],
		})
	}
}

type validateGetFriendMsgHistory struct {
	UserId    int   `validate:"gt=0"`
	MessageId int   `validate:"ne=0"`
	Count     []int `validate:"omitempty,dive,gt=0"`
}

type GetFriendMsgHistoryResp struct {
	Messages []event.MessagePrivate `json:"messages" mapstructure:"messages"` // 获取的消息
}

// GetFriendMsgHistory 获取好友历史聊天记录
func (c LgrCaller) GetFriendMsgHistory(userId int, messageId int, count ...int) (*GetFriendMsgHistoryResp, error) {
	err := validate.Struct(&validateGetFriendMsgHistory{userId, messageId, count})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetFriendMsgHistoryResp](c.PostReq(GetFriendMsgHistory(userId, messageId, count...)))
}

/*
GetGroupMsgHistory 获取群历史聊天记录

https://lagrange-onebot.apifox.cn/236981748e0

参数:

	userId: 用户 Uin
	messageId: 消息 ID
	count(可选): 消息数量 <默认值: 20>
*/
func GetGroupMsgHistory(userId, messageId int, count ...int) *Request {
	if len(count) == 0 {
		return NewReq("get_group_msg_history", map[string]any{
			"user_id":    userId,
			"message_id": messageId,
		})
	} else {
		return NewReq("get_group_msg_history", map[string]any{
			"user_id":     userId,
			"message_id":  messageId,
			"message_num": count[0],
		})
	}
}

type validateGetGroupMsgHistory struct {
	GroupId   int   `validate:"gt=0"`
	MessageId int   `validate:"ne=0"`
	Count     []int `validate:"omitempty,dive,gt=0"`
}

type GetGroupMsgHistoryResp struct {
	Messages []event.MessageGroup `json:"messages" mapstructure:"message"` // 获取的消息
}

// GetGroupMsgHistory 获取群历史聊天记录
func (c LgrCaller) GetGroupMsgHistory(groupId int, messageId int, count ...int) (*GetGroupMsgHistoryResp, error) {
	err := validate.Struct(&validateGetGroupMsgHistory{groupId, messageId, count})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetGroupMsgHistoryResp](c.PostReq(GetGroupMsgHistory(groupId, messageId, count...)))
}

// func GetMsg() *Request
// std: [GetMsg] https://lagrange-onebot.apifox.cn/236981756e0

// func GetMusicArk() *Request
// [Deprecated] https://lagrange-onebot.apifox.cn/236981813e0

/*
GroupPoke 群里戳一戳

https://lagrange-onebot.apifox.cn/236981842e0

参数:

	groupId: 群 Uin
	userId: 用户 Uin
*/
func GroupPoke(groupId, userId int) *Request {
	return NewReq("group_poke", map[string]any{
		"group_id": groupId,
		"user_id":  userId,
	})
}

type validateGroupPoke struct {
	GroupId int `validate:"gt=0"`
	UserId  int `validate:"gt=0"`
}

// GroupPoke 群里戳一戳
func (c LgrCaller) GroupPoke(groupId, userId int) error {
	err := validate.Struct(&validateGroupPoke{groupId, userId})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(GroupPoke(groupId, userId))
}

/*
MarkMsgAsRead 标记消息为已读

https://lagrange-onebot.apifox.cn/236981846e0

参数:

	messageId: 消息 ID
*/
func MarkMsgAsRead(messageId int) *Request {
	return NewReq("mark_msg_as_read", map[string]any{
		"message_id": messageId,
	})
}

// GroupPoke 群里戳一戳
func (c LgrCaller) MarkMsgAsRead(messageId int) error {
	err := validate.msgIdNe0(messageId)
	if err != nil {
		return err
	}
	return c.PostReqNoResp(MarkMsgAsRead(messageId))
}

/*
SendForwardMsg_Lgr 构造合并转发消息 [有异常]

https://lagrange-onebot.apifox.cn/236981861e0

获取的 Res Id 是属于群的, 在私聊中发送会导致图片等资源无法加载

参数:

	messages: Node segment array
*/
func SendForwardMsg_Lgr(messages message.SegmentArray) *Request {
	return NewReq("send_forward_msg", map[string]any{
		"messages": messages,
	})
}

type SendForwardMsg_LgrResp = string // Res Id

// SendForwardMsg 构造合并转发消息 [有异常]
func (c LgrCaller) SendForwardMsg(messages message.SegmentArray) (SendForwardMsg_LgrResp, error) {
	err := validate.require(messages)
	if err != nil {
		return "", err
	}
	return DecodeResponseAssertion[string](c.PostReq(SendForwardMsg_Lgr(messages)))
}

/*
SendGroupAiRecord 发送群 Ai 语音

https://lagrange-onebot.apifox.cn/236981906e0

参数:

	character: 语音声色
	groupId: 群 Uin
	text: 语音文本
	chatType(可选): 语音类型 <enum: 1, 2> <默认值: 1>
*/
func SendGroupAiRecord(character string, groupId int, text string, chatType ...int) *Request {
	if len(chatType) == 0 {
		return NewReq("send_group_ai_record", map[string]any{
			"character": character,
			"group_id":  groupId,
			"text":      text,
		})
	} else {
		return NewReq("send_group_ai_record", map[string]any{
			"character": character,
			"group_id":  groupId,
			"text":      text,
			"chat_type": chatType[0],
		})
	}
}

/*
SendGroupForwardMsg 发送群聊合并转发消息

https://lagrange-onebot.apifox.cn/236981913e0

参数:

	groupId: 群 Uin
	messages: Node segment array
*/
func SendGroupForwardMsg(groupId int, messages message.SegmentArray) *Request {
	return NewReq("send_group_forward_msg", map[string]any{
		// "group_id": strconv.Itoa(groupId),
		"group_id": groupId, // ok
		"messages": messages,
	})
}

type validateSendGroupForwardMsg struct {
	GroupId  int                  `validate:"gt=0"`
	Messages message.SegmentArray `validate:"required,min=1"`
}

type SendAnyForwardMsgResp struct {
	MessageId int    `json:"message_id" mapstructure:"message_id"` // 消息 ID
	ForwardId string `json:"forward_id" mapstructure:"forward_id"` // 转发消息 ID
}

type SendGroupForwardMsgResp = SendAnyForwardMsgResp

// SendGroupForwardMsg 发送群聊合并转发消息
func (c LgrCaller) SendGroupForwardMsg(groupId int, nodes message.SegmentArray) (*SendAnyForwardMsgResp, error) {
	err := validate.Struct(&validateSendGroupForwardMsg{groupId, nodes})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[SendAnyForwardMsgResp](c.PostReq(SendGroupForwardMsg(groupId, nodes)))
}

// func SendGroupMsg(groupId int, message any, autoEscape ...bool) *Request
// std: [SendGroupMsg] https://lagrange-onebot.apifox.cn/236981917e0

// func SendMsg(messageType string, userId, groupId int, message any, autoEscape ...bool) *Request
// std: [SendMsg] https://lagrange-onebot.apifox.cn/236981924e0

/*
SendPrivateForwardMsg 发送私聊合并转发消息

https://lagrange-onebot.apifox.cn/236981928e0

参数:

	userId: 用户 Uin
	message: Node segment array
*/
func SendPrivateForwardMsg(userId int, nodes message.SegmentArray) *Request {
	return NewReq("send_private_forward_msg", map[string]any{
		// "user_id": strconv.Itoa(userId),
		"user_id":  userId,
		"messages": nodes,
	})
}

type validateSendPrivateForwardMsg struct {
	UserId   int                  `validate:"gt=0"`
	Messages message.SegmentArray `validate:"required,min=1"`
}

type SendPrivateForwardMsgResp = SendAnyForwardMsgResp

// SendPrivateForwardMsg 发送私聊合并转发消息
func (c LgrCaller) SendPrivateForwardMsg(userId int, nodes message.SegmentArray) (*SendAnyForwardMsgResp, error) {
	err := validate.Struct(&validateSendPrivateForwardMsg{userId, nodes})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[SendAnyForwardMsgResp](c.PostReq(SendPrivateForwardMsg(userId, nodes)))
}

// func SendPrivateMsg(userId int, message any) *Request
// std: [SendPrivateMsg] https://lagrange-onebot.apifox.cn/236981933e0

/*
SetEssenceMsg 设置精华消息

https://lagrange-onebot.apifox.cn/236981937e0

参数:

	messageId: 消息 ID
*/
func SetEssenceMsg(messageId int) *Request {
	return NewReq("set_essence_msg", map[string]any{
		"message_id": messageId,
	})
}

func (c LgrCaller) SetEssenceMsg(messageId int) error {
	err := validate.msgIdNe0(messageId)
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetEssenceMsg(messageId))
}
