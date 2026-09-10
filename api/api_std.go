package api

import (
	"fmt"
	"time"

	"github.com/Miuzarte/EasyOnebot/event"
	"github.com/Miuzarte/EasyOnebot/message"
)

/*
SendPrivateMsg 发送私聊消息

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#send_private_msg-%E5%8F%91%E9%80%81%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF

参数:

	userId: 对方 QQ 号
	msg: 要发送的内容
	autoEscape: 消息内容是否作为纯文本发送（即不解析 CQ 码），只在 `message` 字段是字符串时有效
*/
func SendPrivateMsg(userId int, msg any, autoEscape ...bool) *Request {
	if len(autoEscape) == 0 {
		return NewReq("send_private_msg", map[string]any{
			"user_id": userId,
			"message": message.ConstraintMessage(msg),
		})
	} else {
		return NewReq("send_private_msg", map[string]any{
			"user_id":     userId,
			"message":     message.ConstraintMessage(msg),
			"auto_escape": autoEscape[0],
		})
	}
}

type validateSendPrivateMsg struct {
	UserId     int    `validate:"gt=0"`
	Message    any    `validate:"required"`
	AutoEscape []bool `validate:"omitempty"`
}

type SendPrivateMsgResp = SendAnyMsgResp

// SendPrivateMsg 发送私聊消息
func (c StdCaller) SendPrivateMsg(userId int, message any, autoEscape ...bool) (*SendPrivateMsgResp, error) {
	err := validate.Struct(&validateSendPrivateMsg{userId, message, autoEscape})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[SendPrivateMsgResp](c.PostReq(SendPrivateMsg(userId, message, autoEscape...)))
}

/*
SendGroupMsg 发送群消息

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#send_group_msg-%E5%8F%91%E9%80%81%E7%BE%A4%E6%B6%88%E6%81%AF

参数:

	groupId: 群号
	msg: 要发送的内容
	autoEscape: 消息内容是否作为纯文本发送（即不解析 CQ 码），只在 `message` 字段是字符串时有效
*/
func SendGroupMsg(groupId int, msg any, autoEscape ...bool) *Request {
	if len(autoEscape) == 0 {
		return NewReq("send_group_msg", map[string]any{
			"group_id": groupId,
			"message":  message.ConstraintMessage(msg),
		})
	} else {
		return NewReq("send_group_msg", map[string]any{
			"group_id":    groupId,
			"message":     message.ConstraintMessage(msg),
			"auto_escape": autoEscape[0],
		})
	}
}

type validateSendGroupMsg struct {
	GroupId    int    `validate:"gt=0"`
	Message    any    `validate:"required"`
	AutoEscape []bool `validate:"omitempty"`
}

type SendGroupMsgResp = SendAnyMsgResp

// SendGroupMsg 发送群消息
func (c StdCaller) SendGroupMsg(groupId int, msg any, autoEscape ...bool) (*SendGroupMsgResp, error) {
	err := validate.Struct(&validateSendGroupMsg{groupId, msg, autoEscape})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[SendGroupMsgResp](c.PostReq(SendGroupMsg(groupId, msg, autoEscape...)))
}

/*
SendMsg 发送消息

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#send_msg-%E5%8F%91%E9%80%81%E6%B6%88%E6%81%AF

参数:

	messageType: 消息类型，支持 `private`、`group`，分别对应私聊、群组，如不传入，则根据传入的 `*_id` 参数判断
	userId: 对方 QQ 号（消息类型为 `private` 时需要）
	groupId: 群号（消息类型为 `group` 时需要）
	msg: 要发送的内容
	autoEscape: 消息内容是否作为纯文本发送（即不解析 CQ 码），只在 `message` 字段是字符串时有效
*/
func SendMsg(messageType string, userId, groupId int, msg any, autoEscape ...bool) *Request {
	if len(autoEscape) == 0 {
		return NewReq("send_msg", map[string]any{
			"message_type": messageType,
			"user_id":      userId,
			"group_id":     groupId,
			"message":      message.ConstraintMessage(msg),
		})
	} else {
		return NewReq("send_msg", map[string]any{
			"message_type": messageType,
			"user_id":      userId,
			"group_id":     groupId,
			"message":      message.ConstraintMessage(msg),
			"auto_escape":  autoEscape[0],
		})
	}
}

type validateSendMsg struct {
	MessageType string `validate:"required,oneof=private group"`
	UserId      int    `validate:"required_if=MessageType private"`
	GroupId     int    `validate:"required_if=MessageType group"`
	Message     any    `validate:"required"`
	AutoEscape  []bool `validate:"omitempty"`
}

type SendAnyMsgResp struct {
	MessageId int `json:"message_id" mapstructure:"message_id"` // 消息 ID
}

// SendMsg 发送消息
func (c StdCaller) SendMsg(messageType string, userId, groupId int, msg any, autoEscape ...bool) (*SendAnyMsgResp, error) {
	err := validate.Struct(&validateSendMsg{messageType, userId, groupId, msg, autoEscape})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[SendAnyMsgResp](c.PostReq(SendMsg(messageType, userId, groupId, msg, autoEscape...)))
}

/*
DeleteMsg 撤回消息

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#delete_msg-%E6%92%A4%E5%9B%9E%E6%B6%88%E6%81%AF

参数:

	messageId: 消息 ID
*/
func DeleteMsg(messageId int) *Request {
	return NewReq("delete_msg", map[string]any{
		"message_id": messageId,
	})
}

// DeleteMsg 撤回消息
func (c StdCaller) DeleteMsg(messageId int) error {
	err := validate.msgIdNe0(messageId)
	if err != nil {
		return err
	}
	return c.PostReqNoResp(DeleteMsg(messageId))
}

// DeleteMsgSf 撤回消息的 singleflight 版本
func (c StdCaller) DeleteMsgSf(messageId int) (error, bool) {
	_, err, shared := sfg.Do(fmt.Sprintf("delete_msg:%d", messageId), func() (any, error) {
		return nil, c.DeleteMsg(messageId)
	})
	return err, shared
}

/*
GetMsg 获取消息

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_msg-%E8%8E%B7%E5%8F%96%E6%B6%88%E6%81%AF

参数:

	messageId: 消息 ID
*/
func GetMsg(messageId int) *Request {
	return NewReq("get_msg", map[string]any{
		"message_id": messageId,
	})
}

type GetMsgResp struct {
	Time        int                  `json:"time" mapstructure:"time"`                 // 发送时间
	MessageType string               `json:"message_type" mapstructure:"message_type"` // 消息类型，同 `消息事件`
	MessageId   int                  `json:"message_id" mapstructure:"message_id"`     // 消息 ID
	RealId      int                  `json:"real_id" mapstructure:"real_id"`           // 消息真实 ID
	Sender      event.GroupSender    `json:"sender" mapstructure:"sender"`             // 发送人信息，同 `消息事件`
	Message     message.SegmentArray `json:"message" mapstructure:"message"`           // 消息内容
}

// GetMsg 获取消息
func (c StdCaller) GetMsg(messageId int) (*GetMsgResp, error) {
	err := validate.msgIdNe0(messageId)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetMsgResp](c.PostReq(GetMsg(messageId)))
}

// GetMsgSf 获取消息的 singleflight 版本
func (c StdCaller) GetMsgSf(messageId int) (*GetMsgResp, error, bool) {
	resp, err, shared := sfg.Do(fmt.Sprintf("get_msg:%d", messageId), func() (any, error) {
		result, err := c.GetMsg(messageId)
		if err == nil {
			result.Message.TryAtoi()
		}
		return result, err
	})
	return resp.(*GetMsgResp), err, shared
}

/*
GetForwardMsg 获取合并转发消息

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_forward_msg-%E8%8E%B7%E5%8F%96%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91%E6%B6%88%E6%81%AF

参数:

	id: 合并转发 ID
*/
func GetForwardMsg(id string) *Request {
	return NewReq("get_forward_msg", map[string]any{"id": id})
}

type GetForwardMsgResp struct {
	Message message.SegmentArray `json:"message" mapstructure:"message"`
}

// GetForwardMsg 获取合并转发消息
func (c StdCaller) GetForwardMsg(id string) (*GetForwardMsgResp, error) {
	err := validate.require(id)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetForwardMsgResp](c.PostReq(GetForwardMsg(id)))
}

// GetForwardMsgSf 获取合并转发消息的 singleflight 版本
func (c StdCaller) GetForwardMsgSf(id string) (*GetForwardMsgResp, error, bool) {
	resp, err, shared := sfg.Do(fmt.Sprintf("get_forward_msg:%s", id), func() (any, error) {
		return c.GetForwardMsg(id)
	})
	return resp.(*GetForwardMsgResp), err, shared
}

/*
SendLike 发送好友赞

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#send_like-%E5%8F%91%E9%80%81%E5%A5%BD%E5%8F%8B%E8%B5%9E

参数:

	userId: 对方 QQ 号
	times: 赞的次数，每个好友每天最多 10 次
*/
func SendLike(userId int, times int) *Request {
	return NewReq("send_like", map[string]any{
		"user_id": userId,
		"times":   times,
	})
}

type validateSendLike struct {
	UserId int `validate:"gt=0"`
	// Times int `validate:"gte=0,lte=10"`
	Times int `validate:"gte=0,lte=20"` // SVIP 20
}

// SendLike 发送好友赞
func (c StdCaller) SendLike(userId int, times int) error {
	err := validate.Struct(&validateSendLike{userId, times})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SendLike(userId, times))
}

/*
SetGroupKick 群组踢人

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_kick-%E7%BE%A4%E7%BB%84%E8%B8%A2%E4%BA%BA

参数:

	groupId: 群号
	userId: 要踢的 QQ 号
	rejectAddRequest: 拒绝此人的加群请求
*/
func SetGroupKick(groupId, userId int, rejectAddRequest bool) *Request {
	return NewReq("set_group_kick", map[string]any{
		"group_id":           groupId,
		"user_id":            userId,
		"reject_add_request": rejectAddRequest,
	})
}

type validateSetGroupKick struct {
	GroupId int `validate:"gt=0"`
	UserId  int `validate:"gt=0"`
}

// SetGroupKick 群组踢人
func (c StdCaller) SetGroupKick(groupId, userId int, rejectAddRequest bool) error {
	err := validate.Struct(&validateSetGroupKick{groupId, userId})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupKick(groupId, userId, rejectAddRequest))
}

/*
SetGroupBan 群组单人禁言

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_ban-%E7%BE%A4%E7%BB%84%E5%8D%95%E4%BA%BA%E7%A6%81%E8%A8%80

参数:

	groupId: 群号
	userId: 要禁言的 QQ 号
	duration: 禁言时长，单位秒，0 表示取消禁言
*/
func SetGroupBan(groupId, userId int, duration time.Duration) *Request {
	return NewReq("set_group_ban", map[string]any{
		"group_id": groupId,
		"user_id":  userId,
		"duration": duration / time.Second,
	})
}

type validateSetGroupBan struct {
	GroupId  int           `validate:"gt=0"`
	UserId   int           `validate:"gt=0"`
	Duration time.Duration `validate:"gte=0"`
}

// SetGroupBan 群组单人禁言
func (c StdCaller) SetGroupBan(groupId, userId int, duration time.Duration) error {
	err := validate.Struct(&validateSetGroupBan{groupId, userId, duration})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupBan(groupId, userId, duration))
}

/*
SetGroupAnonymousBan 群组匿名用户禁言

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_anonymous_ban-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D%E7%94%A8%E6%88%B7%E7%A6%81%E8%A8%80

参数:

	groupId: 群号
	anonymous: 可选，要禁言的匿名用户对象（群消息上报的 `anonymous` 字段）
	flag: 可选，要禁言的匿名用户的 flag（需从群消息上报的数据中获得）
	duration: 禁言时长，单位秒，无法取消匿名用户禁言

	上面的 `anonymous` 和 `anonymous_flag` 两者任选其一传入即可，若都传入，则使用 `anonymous`。
*/
func SetGroupAnonymousBan(groupId int, anonymous *event.Anonymous, flag string, duration time.Duration) *Request {
	return NewReq("set_group_anonymous_ban", map[string]any{
		"group_id":  groupId,
		"anonymous": anonymous,
		"flag":      flag,
		"duration":  duration / time.Second,
	})
}

type validateSetGroupAnonymousBan struct {
	GroupId   int              `validate:"gt=0"`
	Anonymous *event.Anonymous `validate:"required_without=Flag"`
	Flag      string           `validate:"required_without=Anonymous"`
	Duration  time.Duration    `validate:"gte=0"`
}

// SetGroupAnonymousBan 群组匿名用户禁言
func (c StdCaller) SetGroupAnonymousBan(groupId int, anonymous *event.Anonymous, flag string, duration time.Duration) error {
	err := validate.Struct(&validateSetGroupAnonymousBan{groupId, anonymous, flag, duration})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupAnonymousBan(groupId, anonymous, flag, duration))
}

/*
SetGroupWholeBan 群组全员禁言

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80

参数:

	groupId: 群号
	enable: 是否禁言
*/
func SetGroupWholeBan(groupId int, enable bool) *Request {
	return NewReq("set_group_whole_ban", map[string]any{
		"group_id": groupId,
		"enable":   enable,
	})
}

// SetGroupWholeBan 群组全员禁言
func (c StdCaller) SetGroupWholeBan(groupId int, enable bool) error {
	err := validate.idGt0(groupId)
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupWholeBan(groupId, enable))
}

/*
SetGroupAdmin 群组设置管理员

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80

参数:

	groupId: 群号
	userId: 要设置管理员的 QQ 号
	enable: true 为设置，false 为取消
*/
func SetGroupAdmin(groupId, userId int, enable bool) *Request {
	return NewReq("set_group_admin", map[string]any{
		"group_id": groupId,
		"user_id":  userId,
		"enable":   enable,
	})
}

type validateSetGroupAdmin struct {
	GroupId int `validate:"gt=0"`
	UserId  int `validate:"gt=0"`
}

// SetGroupAdmin 群组设置管理员
func (c StdCaller) SetGroupAdmin(groupId, userId int, enable bool) error {
	err := validate.Struct(&validateSetGroupAdmin{groupId, userId})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupAdmin(groupId, userId, enable))
}

/*
SetGroupAnonymous 群组匿名

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80

参数:

	groupId: 群号
	enable: 是否允许匿名聊天
*/
func SetGroupAnonymous(groupId int, enable bool) *Request {
	return NewReq("set_group_anonymous", map[string]any{
		"group_id": groupId,
		"enable":   enable,
	})
}

// SetGroupAnonymous 群组匿名
func (c StdCaller) SetGroupAnonymous(groupId int, enable bool) error {
	err := validate.idGt0(groupId)
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupAnonymous(groupId, enable))
}

/*
SetGroupCard 设置群名片（群备注）

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80

参数:

	groupId: 群号
	userId: 要设置的 QQ 号
	card: 群名片内容，不填或空字符串表示删除群名片
*/
func SetGroupCard(groupId, userId int, card string) *Request {
	return NewReq("set_group_card", map[string]any{
		"group_id": groupId,
		"user_id":  userId,
		"card":     card,
	})
}

type validateSetGroupCard struct {
	GroupId int `validate:"gt=0"`
	UserId  int `validate:"gt=0"`
}

// SetGroupCard 设置群名片（群备注）
func (c StdCaller) SetGroupCard(groupId, userId int, card string) error {
	err := validate.Struct(&validateSetGroupCard{groupId, userId})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupCard(groupId, userId, card))
}

/*
SetGroupName 设置群名

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80

参数:

	groupId: 群号
	name: 新群名
*/
func SetGroupName(groupId int, name string) *Request {
	return NewReq("set_group_name", map[string]any{
		"group_id": groupId,
		"name":     name,
		// "groupName": groupName, // lagrange
	})
}

type validateSetGroupName struct {
	GroupId int    `validate:"gt=0"`
	Name    string `validate:"required"`
}

// SetGroupName 设置群名
func (c StdCaller) SetGroupName(groupId int, name string) error {
	err := validate.Struct(&validateSetGroupName{groupId, name})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupName(groupId, name))
}

/*
SetGroupLeave 退出群组

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80

参数:

	groupId: 群号
	isDismiss: 是否解散，如果登录号是群主，则仅在此项为 true 时能够解散
*/
func SetGroupLeave(groupId int, isDismiss bool) *Request {
	return NewReq("set_group_leave", map[string]any{
		"group_id":   groupId,
		"is_dismiss": isDismiss,
	})
}

// SetGroupLeave 退出群组
func (c StdCaller) SetGroupLeave(groupId int, isDismiss bool) error {
	err := validate.idGt0(groupId)
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupLeave(groupId, isDismiss))
}

/*
SetGroupSpecialTitle 设置群组专属头衔

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80

参数:

	groupId: 群号
	userId: 要设置的 QQ 号
	specialTitle: 专属头衔，不填或空字符串表示删除专属头衔
	duration: 专属头衔有效期，单位秒，-1 表示永久，不过此项似乎没有效果，可能是只有某些特殊的时间长度有效，有待测试
*/
func SetGroupSpecialTitle(groupId, userId int, specialTitle string, duration time.Duration) *Request {
	if duration == -1 {
		return NewReq("set_group_special_title", map[string]any{
			"group_id":      groupId,
			"user_id":       userId,
			"special_title": specialTitle,
			"duration":      -1,
		})
	} else {
		return NewReq("set_group_special_title", map[string]any{
			"group_id":      groupId,
			"user_id":       userId,
			"special_title": specialTitle,
			"duration":      duration / time.Second,
		})
	}
}

type validateSetGroupSpecialTitle struct {
	GroupId int `validate:"gt=0"`
	UserId  int `validate:"gt=0"`
}

// SetGroupSpecialTitle 设置群组专属头衔
func (c StdCaller) SetGroupSpecialTitle(groupId, userId int, specialTitle string, duration time.Duration) error {
	err := validate.Struct(&validateSetGroupSpecialTitle{groupId, userId})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupSpecialTitle(groupId, userId, specialTitle, duration))
}

/*
SetFriendAddRequest 处理加好友请求

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80

参数:

	flag: 加好友请求的 flag（需从上报的数据中获得）
	approve: 是否同意请求
	remark: 添加后的好友备注（仅在同意时有效）
*/
func SetFriendAddRequest(flag string, approve bool, remark string) *Request {
	return NewReq("set_friend_add_request", map[string]any{
		"flag":    flag,
		"approve": approve,
		"remark":  remark,
	})
}

// SetFriendAddRequest 处理加好友请求
func (c StdCaller) SetFriendAddRequest(flag string, approve bool, remark string) error {
	err := validate.require(flag)
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetFriendAddRequest(flag, approve, remark))
}

/*
SetGroupAddRequest 处理加群请求／邀请

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E7%BE%A4%E8%AF%B7%E6%B1%82%E9%82%80%E8%AF%B7

参数:

	flag: 加群请求的 flag（需从上报的数据中获得）
	subType: `add` 或 `invite`，请求类型（需要和上报消息中的 `sub_type` 字段相符）
	approve: 是否同意请求／邀请
	reason: 拒绝理由（仅在拒绝时有效）
*/
func SetGroupAddRequest(flag, subType string, approve bool, reason string) *Request {
	return NewReq("set_group_add_request", map[string]any{
		"flag":     flag,
		"sub_type": subType,
		"approve":  approve,
		"reason":   reason,
	})
}

type validateSetGroupAddRequest struct {
	Flag    string `validate:"required"`
	SubType string `validate:"required,oneof=add invite"`
}

// SetGroupAddRequest 处理加群请求／邀请
func (c StdCaller) SetGroupAddRequest(flag, subType string, approve bool, reason string) error {
	err := validate.Struct(&validateSetGroupAddRequest{flag, subType})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupAddRequest(flag, subType, approve, reason))
}

/*
GetLoginInfo 获取登录号信息

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_login_info-%E8%8E%B7%E5%8F%96%E7%99%BB%E5%BD%95%E5%8F%B7%E4%BF%A1%E6%81%AF
*/
func GetLoginInfo() *Request {
	return NewReq("get_login_info", nil)
}

type GetLoginInfoResp = LoginInfo

type LoginInfo struct {
	UserId   int    `json:"user_id" mapstructure:"user_id"`   // QQ 号
	Nickname string `json:"nickname" mapstructure:"nickname"` // QQ 昵称
}

// GetLoginInfo 获取登录号信息
func (c StdCaller) GetLoginInfo() (*GetLoginInfoResp, error) {
	return DecodeResponse[GetLoginInfoResp](c.PostReq(GetLoginInfo()))
}

/*
GetStrangerInfo 获取陌生人信息

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_stranger_info-%E8%8E%B7%E5%8F%96%E9%99%8C%E7%94%9F%E4%BA%BA%E4%BF%A1%E6%81%AF

参数:

	userId: QQ 号
	noCache: 是否不使用缓存（使用缓存可能更新不及时，但响应更快）
*/
func GetStrangerInfo(userId int, noCache bool) *Request {
	return NewReq("get_stranger_info", map[string]any{
		"user_id":  userId,
		"no_cache": noCache,
	})
}

type GetStrangerInfoResp = StrangerInfo

type StrangerInfo struct {
	UserId   int    `json:"user_id" mapstructure:"user_id"`   // QQ 号
	Nickname string `json:"nickname" mapstructure:"nickname"` // 昵称
	Sex      string `json:"sex" mapstructure:"sex"`           // 性别，`male` 或 `female` 或 `unknown`
	Age      int    `json:"age" mapstructure:"age"`           // 年龄
}

// GetStrangerInfo 获取陌生人信息
func (c StdCaller) GetStrangerInfo(userId int, noCache bool) (*GetStrangerInfoResp, error) {
	err := validate.idGt0(userId)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetStrangerInfoResp](c.PostReq(GetStrangerInfo(userId, noCache)))
}

/*
GetFriendList 获取好友列表

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_friend_list-%E8%8E%B7%E5%8F%96%E5%A5%BD%E5%8F%8B%E5%88%97%E8%A1%A8
*/
func GetFriendList() *Request {
	return NewReq("get_friend_list", nil)
}

type GetFriendListResp = []FriendInfo

type FriendInfo struct {
	UserId   int    `json:"user_id" mapstructure:"user_id"`   // QQ 号
	Nickname string `json:"nickname" mapstructure:"nickname"` // 昵称
	Remark   string `json:"remark" mapstructure:"remark"`     // 备注名
}

// GetFriendList 获取好友列表
func (c StdCaller) GetFriendList() (GetFriendListResp, error) {
	return DecodeResponseSlice[GetFriendListResp](c.PostReq(GetFriendList()))
}

/*
GetGroupInfo 获取群信息

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E4%BF%A1%E6%81%AF

参数:

	groupId: 群号
	noCache: 是否不使用缓存（使用缓存可能更新不及时，但响应更快）
*/
func GetGroupInfo(groupId int, noCache ...bool) *Request {
	if len(noCache) == 0 {
		return NewReq("get_group_info", map[string]any{
			"group_id": groupId,
		})
	} else {
		return NewReq("get_group_info", map[string]any{
			"group_id": groupId,
			"no_cache": noCache[0],
		})
	}
}

type GetGroupInfoResp = GroupInfo

type GroupInfo struct {
	GroupId        int    `json:"group_id" mapstructure:"group_id"`                 // 群号
	GroupName      string `json:"group_name" mapstructure:"group_name"`             // 群名称
	MemberCount    int    `json:"member_count" mapstructure:"member_count"`         // 成员数
	MaxMemberCount int    `json:"max_member_count" mapstructure:"max_member_count"` // 最大成员数量（群容量）
}

// GetGroupInfo 获取群信息
func (c StdCaller) GetGroupInfo(groupId int, noCache bool) (*GetGroupInfoResp, error) {
	err := validate.idGt0(groupId)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetGroupInfoResp](c.PostReq(GetGroupInfo(groupId, noCache)))
}

/*
GetGroupList 获取群列表

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%88%97%E8%A1%A8
*/
func GetGroupList(noCache ...bool) *Request {
	if len(noCache) == 0 {
		return NewReq("get_group_list", nil)
	} else {
		return NewReq("get_group_list", map[string]any{
			"no_cache": noCache[0],
		})
	}
}

type GetGroupListResp = []GroupInfo

// GetGroupList 获取群列表
func (c StdCaller) GetGroupList() (GetGroupListResp, error) {
	return DecodeResponseSlice[GetGroupListResp](c.PostReq(GetGroupList()))
}

/*
GetGroupMemberInfo 获取群成员信息

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_member_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E4%BF%A1%E6%81%AF

参数:

	groupId: 群号
	userId: QQ 号
	noCache: 是否不使用缓存（使用缓存可能更新不及时，但响应更快）
*/
func GetGroupMemberInfo(groupId, userId int, noCache ...bool) *Request {
	if len(noCache) == 0 {
		return NewReq("get_group_member_info", map[string]any{
			"group_id": groupId,
			"user_id":  userId,
		})
	} else {
		return NewReq("get_group_member_info", map[string]any{
			"group_id": groupId,
			"user_id":  userId,
			"no_cache": noCache[0],
		})
	}
}

type GetGroupMemberInfoResp = GroupMemberInfo

type GroupMemberInfo struct {
	GroupId        int    `json:"group_id" mapstructure:"group_id"`               // 群号
	UserId         int    `json:"user_id" mapstructure:"user_id"`                 // QQ 号
	Nickname       string `json:"nickname" mapstructure:"nickname"`               // 昵称
	Card           string `json:"card" mapstructure:"card"`                       // 群名片／备注
	Sex            string `json:"sex" mapstructure:"sex"`                         // 性别，`male` 或 `female` 或 `unknown`
	Age            int    `json:"age" mapstructure:"age"`                         // 年龄
	Area           string `json:"area" mapstructure:"area"`                       // 地区
	JoinTime       int    `json:"join_time" mapstructure:"join_time"`             // 加群时间戳
	LastSentTime   int    `json:"last_sent_time" mapstructure:"last_sent_time"`   // 最后发言时间戳
	Level          string `json:"level" mapstructure:"level"`                     // 成员等级
	Role           string `json:"role" mapstructure:"role"`                       // 角色，`owner` 或 `admin` 或 `member`
	Unfriendly     bool   `json:"unfriendly" mapstructure:"unfriendly"`           // 是否不良记录成员
	CardChangeable bool   `json:"card_changeable" mapstructure:"card_changeable"` // 是否允许修改群名片
}

type validateGetGroupMemberInfo struct {
	GroupId int `validate:"gt=0"`
	UserId  int `validate:"gt=0"`
}

// GetGroupMemberInfo 获取群成员信息
func (c StdCaller) GetGroupMemberInfo(groupId, userId int, noCache bool) (*GetGroupMemberInfoResp, error) {
	err := validate.Struct(&validateGetGroupMemberInfo{groupId, userId})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetGroupMemberInfoResp](c.PostReq(GetGroupMemberInfo(groupId, userId, noCache)))
}

/*
GetGroupMemberList 获取群成员列表

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_member_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E5%88%97%E8%A1%A8

参数:

	groupId: 群号

响应内容为 JSON 数组，每个元素的内容和上面的 `get_group_member_info` 接口相同，
但对于同一个群组的同一个成员，获取列表时和获取单独的成员信息时，某些字段可能有所不同，
例如 `area`、`title` 等字段在获取列表时无法获得，具体应以单独的成员信息为准。
*/
func GetGroupMemberList(groupId int) *Request {
	return NewReq("get_group_member_list", map[string]any{
		"group_id": groupId,
	})
}

type GetGroupMemberListResp = []GroupMemberInfo

// GetGroupMemberList 获取群成员列表
func (c StdCaller) GetGroupMemberList(groupId int) (GetGroupMemberListResp, error) {
	err := validate.idGt0(groupId)
	if err != nil {
		return nil, err
	}
	return DecodeResponseSlice[GetGroupMemberListResp](c.PostReq(GetGroupMemberList(groupId)))
}

/*
GetGroupHonorInfo 获取群荣誉信息

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_honor_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E8%8D%A3%E8%AA%89%E4%BF%A1%E6%81%AF

参数:

	groupId: 群号
	typ: 要获取的群荣誉类型，可传入 `talkative` `performer` `legend` `strong_newbie` `emotion` 以分别获取单个类型的群荣誉数据，或传入 `all` 获取所有数据
*/
func GetGroupHonorInfo(groupId int, typ string) *Request {
	return NewReq("get_group_honor_info", map[string]any{
		"group_id": groupId,
		"type":     typ,
	})
}

type GetGroupHonorInfoResp = GroupHonorInfo

type GroupHonorInfo struct {
	GroupId          int              `json:"group_id" mapstructure:"group_id"`                     // 群号
	CurrentTalkative CurrentTalkative `json:"current_talkative" mapstructure:"current_talkative"`   // 当前龙王，仅 `type` 为 `talkative` 或 `all` 时有数据
	TalkativeList    []Honor          `json:"talkative_list" mapstructure:"talkative_list"`         // 历史龙王，仅 `type` 为 `talkative` 或 `all` 时有数据
	PerformerList    []Honor          `json:"performer_list" mapstructure:"performer_list"`         // 群聊之火，仅 `type` 为 `performer` 或 `all` 时有数据
	LegendList       []Honor          `json:"legend_list" mapstructure:"legend_list"`               // 群聊炽焰，仅 `type` 为 `legend` 或 `all` 时有数据
	StrongNewbieList []Honor          `json:"strong_newbie_list" mapstructure:"strong_newbie_list"` // 冒尖小春笋，仅 `type` 为 `strong_newbie` 或 `all` 时有数据
	EmotionList      []Honor          `json:"emotion_list" mapstructure:"emotion_list"`             // 快乐之源，仅 `type` 为 `emotion` 或 `all` 时有数据
}

type CurrentTalkative struct {
	UserId   int    `json:"user_id" mapstructure:"user_id"`     // QQ 号
	Nickname string `json:"nickname" mapstructure:"nickname"`   // 昵称
	Avatar   string `json:"avatar" mapstructure:"avatar"`       // 头像 URL
	DayCount int    `json:"day_count" mapstructure:"day_count"` // 持续天数
}

type Honor struct {
	UserId      int    `json:"user_id" mapstructure:"user_id"`         // QQ 号
	Nickname    string `json:"nickname" mapstructure:"nickname"`       // 昵称
	Avatar      string `json:"avatar" mapstructure:"avatar"`           // 头像 URL
	Description string `json:"description" mapstructure:"description"` // 荣誉描述
}

type validateGetGroupHonorInfo struct {
	GroupId int    `validate:"gt=0"`
	Type    string `validate:"required,oneof=all talkative performer legend strong_newbie emotion"`
}

// GetGroupHonorInfo 获取群荣誉信息
func (c StdCaller) GetGroupHonorInfo(groupId int, typ string) (*GetGroupHonorInfoResp, error) {
	err := validate.Struct(&validateGetGroupHonorInfo{groupId, typ})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetGroupHonorInfoResp](c.PostReq(GetGroupHonorInfo(groupId, typ)))
}

/*
GetCookies 获取 Cookies

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_cookies-%E8%8E%B7%E5%8F%96-cookies

见 NapCat 端点 [get_cookies] (docs/onebot-napcat-endpoints.md)

参数:

	domain: 需要获取 cookies 的域名
*/
func GetCookies(domain string) *Request {
	return NewReq("get_cookies", map[string]any{
		"domain": domain,
	})
}

type GetCookiesResp struct {
	Cookies string `json:"cookies" mapstructure:"cookies"` // Cookies
}

// GetCookies 获取 Cookies
func (c StdCaller) GetCookies(domain string) (*GetCookiesResp, error) {
	err := validate.require(domain)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetCookiesResp](c.PostReq(GetCookies(domain)))
}

/*
GetCsrfToken 获取 CSRF Token

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_csrf_token-%E8%8E%B7%E5%8F%96-csrf-token

见 NapCat 端点 [get_csrf_token] (docs/onebot-napcat-endpoints.md)
*/
func GetCsrfToken() *Request {
	return NewReq("get_csrf_token", nil)
}

type GetCsrfTokenResp struct {
	Token int `json:"token" mapstructure:"token"` // CSRF Token
	// Token string // lagrange
}

// GetCsrfToken 获取 CSRF Token
func (c StdCaller) GetCsrfToken() (*GetCsrfTokenResp, error) {
	return DecodeResponse[GetCsrfTokenResp](c.PostReq(GetCsrfToken()))
}

/*
GetCredentials 获取 QQ 相关接口凭证

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_credentials-%E8%8E%B7%E5%8F%96-qq-%E7%9B%B8%E5%85%B3%E6%8E%A5%E5%8F%A3%E5%87%AD%E8%AF%81

即上面两个接口的合并。

见 NapCat 端点 [get_credentials] (docs/onebot-napcat-endpoints.md)

参数:

	domain: 需要获取 cookies 的域名
*/
func GetCredentials(domain string) *Request {
	return NewReq("get_credentials", map[string]any{
		"domain": domain,
	})
}

type GetCredentialsResp struct {
	Cookies   string `json:"cookies" mapstructure:"cookies"`       // Cookies
	CsrfToken int    `json:"csrf_token" mapstructure:"csrf_token"` // CSRF Token
	// CsrfToken string // lagrange
}

// GetCredentials 获取 QQ 相关接口凭证
func (c StdCaller) GetCredentials(domain string) (*GetCredentialsResp, error) {
	err := validate.require(domain)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetCredentialsResp](c.PostReq(GetCredentials(domain)))
}

/*
GetRecord 获取语音

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_record-%E8%8E%B7%E5%8F%96%E8%AF%AD%E9%9F%B3

提示：要使用此接口，通常需要安装 ffmpeg，请参考 OneBot 实现的相关说明。

参数:

	file: 收到的语音文件名（消息段的 `file` 参数），如 `0B38145AA44505000B38145AA4450500.silk`
	outFormat: 要转换到的格式，目前支持 `mp3`、`amr`、`wma`、`m4a`、`spx`、`ogg`、`wav`、`flac`
*/
func GetRecord(file, outFormat string) *Request {
	return NewReq("get_record", map[string]any{
		"file":       file,
		"out_format": outFormat,
	})
}

type GetRecordResp struct {
	File string `json:"file" mapstructure:"file"` // 转换后的语音文件路径，如 `/home/somebody/cqhttp/data/record/0B38145AA44505000B38145AA4450500.mp3`
}

type validateGetRecord struct {
	File      string `validate:"required"`
	OutFormat string `validate:"required,oneof=mp3 amr wma m4a spx ogg wav flac"`
}

// GetRecord 获取语音
func (c StdCaller) GetRecord(file, outFormat string) (*GetRecordResp, error) {
	err := validate.Struct(&validateGetRecord{file, outFormat})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetRecordResp](c.PostReq(GetRecord(file, outFormat)))
}

/*
GetImage 获取图片

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_image-%E8%8E%B7%E5%8F%96%E5%9B%BE%E7%89%87

参数:

	file: 收到的图片文件名（消息段的 `file` 参数），如 `6B4DE3DFD1BD271E3297859D41C530F5.jpg`
*/
func GetImage(file string) *Request {
	return NewReq("get_image", map[string]any{
		"file": file,
	})
}

type GetImageResp struct {
	File string `json:"file" mapstructure:"file"` // 下载后的图片文件路径，如 `/home/somebody/cqhttp/data/image/6B4DE3DFD1BD271E3297859D41C530F5.jpg`
}

// GetImage 获取图片
func (c StdCaller) GetImage(file string) (*GetImageResp, error) {
	err := validate.require(file)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetImageResp](c.PostReq(GetImage(file)))
}

/*
CanSendImage 是否可以发送图片

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#can_send_image-%E6%A3%80%E6%9F%A5%E6%98%AF%E5%90%A6%E5%8F%AF%E4%BB%A5%E5%8F%91%E9%80%81%E5%9B%BE%E7%89%87
*/
func CanSendImage() *Request {
	return NewReq("can_send_image", nil)
}

type CanSendImageResp struct {
	Yes bool `json:"yes" mapstructure:"yes"` // 是或否
}

// CanSendImage 是否可以发送图片
func (c StdCaller) CanSendImage() (*CanSendImageResp, error) {
	return DecodeResponse[CanSendImageResp](c.PostReq(CanSendImage()))
}

/*
CanSendRecord 是否可以发送语音

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#can_send_record-%E6%A3%80%E6%9F%A5%E6%98%AF%E5%90%A6%E5%8F%AF%E4%BB%A5%E5%8F%91%E9%80%81%E8%AF%AD%E9%9F%B3
*/
func CanSendRecord() *Request {
	return NewReq("can_send_record", nil)
}

type CanSendRecordResp struct {
	Yes bool `json:"yes" mapstructure:"yes"` // 是或否
}

// CanSendRecord 是否可以发送语音
func (c StdCaller) CanSendRecord() (*CanSendRecordResp, error) {
	return DecodeResponse[CanSendRecordResp](c.PostReq(CanSendRecord()))
}

/*
GetStatus 获取运行状态

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_status-%E8%8E%B7%E5%8F%96%E8%BF%90%E8%A1%8C%E7%8A%B6%E6%80%81

通常情况下建议只使用 `online` 和 `good` 这两个字段来判断运行状态，因为根据 OneBot 实现的不同，其它字段可能完全不同。
*/
func GetStatus() *Request {
	return NewReq("get_status", nil)
}

type GetStatusResp struct {
	Online         bool `json:"online" mapstructure:"online"`                   // 当前 QQ 在线，`null` 表示无法查询到在线状态
	Good           bool `json:"good" mapstructure:"good"`                       // 状态符合预期，意味着各模块正常运行、功能正常，且 QQ 在线
	AppEnabled     bool `json:"app_enabled" mapstructure:"app_enabled"`         // lagrange extension
	AppGood        bool `json:"app_good" mapstructure:"app_good"`               // lagrange extension
	AppInitialized bool `json:"app_initialized" mapstructure:"app_initialized"` // lagrange extension
	PluginsGood    bool `json:"plugins_good" mapstructure:"plugins_good"`       // lagrange extension
	Memory         int  `json:"memory" mapstructure:"memory"`                   // lagrange extension
}

// GetStatus 获取运行状态
func (c StdCaller) GetStatus() (*GetStatusResp, error) {
	return DecodeResponse[GetStatusResp](c.PostReq(GetStatus()))
}

/*
GetVersionInfo 获取版本信息

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_version_info-%E8%8E%B7%E5%8F%96%E7%89%88%E6%9C%AC%E4%BF%A1%E6%81%AF
*/
func GetVersionInfo() *Request {
	return NewReq("get_version_info", nil)
}

type GetVersionInfoResp = VersionInfo

type VersionInfo struct {
	AppName         string `json:"app_name" mapstructure:"app_name"`                 // 应用标识，如 `mirai-native`
	AppVersion      string `json:"app_version" mapstructure:"app_version"`           // 应用版本，如 `1.2.3`
	ProtocolVersion string `json:"protocol_version" mapstructure:"protocol_version"` // OneBot 标准版本，如 `v11`

	NtProtocol string `json:"nt_protocol" mapstructure:"nt_protocol"` // lagrange extension
}

// GetVersionInfo 获取版本信息
func (c StdCaller) GetVersionInfo() (*GetVersionInfoResp, error) {
	return DecodeResponse[GetVersionInfoResp](c.PostReq(GetVersionInfo()))
}

/*
SetRestart 重启 OneBot 实现

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_restart-%E9%87%8D%E5%90%AF-onebot-%E5%AE%9E%E7%8E%B0

由于重启 OneBot 实现同时需要重启 API 服务，这意味着当前的 API 请求会被中断，因此需要异步地重启，接口返回的 `status` 是 `async`。

见 NapCat 端点 [set_restart] (docs/onebot-napcat-endpoints.md)

参数:

	delay: 要延迟的毫秒数，如果默认情况下无法重启，可以尝试设置延迟为 2000 左右
*/
func SetRestart(delay time.Duration) *Request {
	return NewReq("set_restart", map[string]any{
		"delay": delay / time.Millisecond,
	})
}

type validateSetRestart struct {
	Delay time.Duration `validate:"gte=0"`
}

// SetRestart 重启 OneBot 实现
func (c StdCaller) SetRestart(delay time.Duration) error {
	err := validate.Struct(&validateSetRestart{delay})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetRestart(delay))
}

/*
CleanCache 清理缓存

https://github.com/botuniverse/onebot-11/blob/master/api/public.md#clean_cache-%E6%B8%85%E7%90%86%E7%BC%93%E5%AD%98

用于清理积攒了太多的缓存文件。
*/
func CleanCache() *Request {
	return NewReq("clean_cache", nil)
}

// CleanCache 清理缓存
func (c StdCaller) CleanCache() error {
	return c.PostReqNoResp(CleanCache())
}
