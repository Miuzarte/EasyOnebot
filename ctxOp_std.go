package EasyOnebot

import (
	"time"

	"github.com/Miuzarte/EasyOnebot/api"
	"github.com/Miuzarte/EasyOnebot/message"
)

// SendPrivateMsg 发送私聊消息至 Event.UserId
func (c *Ctx) SendPrivateMsg(msg any, autoEscape ...bool) (*api.SendAnyMsgResp, error) {
	return c.Std.SendPrivateMsg(c.Event.UserId, msg, autoEscape...)
}

// SendGroupMsg 发送群聊消息至 Event.GroupId
func (c *Ctx) SendGroupMsg(msg any, autoEscape ...bool) (*api.SendAnyMsgResp, error) {
	return c.Std.SendGroupMsg(c.Event.GroupId, msg, autoEscape...)
}

// SendMsg 发送消息至 Event.MessageType, Event.UserId, Event.GroupId
func (c *Ctx) SendMsg(msg any, autoEscape ...bool) (*api.SendAnyMsgResp, error) {
	return c.Std.SendMsg(c.Event.MessageType, c.Event.UserId, c.Event.GroupId, msg, autoEscape...)
}

// DeleteMsg 撤回 Event.MessageId
func (c *Ctx) DeleteMsg() error {
	err, _ := c.Std.DeleteMsgSf(c.Event.MessageId)
	return err
}

// GetMsg 获取 Event.MessageId
func (c *Ctx) GetMsg() (*api.GetMsgResp, error) {
	resp, err, _ := c.Std.GetMsgSf(c.Event.MessageId)
	// resp.Message.TryAtoi()
	return resp, err
}

// GetForwardMsg 获取合并转发消息 Event.MessageId
func (c *Ctx) GetForwardMsg() (*api.GetForwardMsgResp, error) {
	seg := c.ParsedSegments.GetFirstType(message.TYPE_FORWARD)
	if seg == nil {
		return nil, ErrNoForwardSeg
	}
	id, ok := seg.Data["id"].(string)
	if !ok {
		return nil, ErrNoForwardId
	}
	return c.Std.GetForwardMsg(id)
}

// SendLike 发送好友赞至 Event.UserId
func (c *Ctx) SendLike(times int) error {
	return c.Std.SendLike(c.Event.UserId, times)
}

// SetGroupKick 群组 Event.GroupId 踢 Event.UserId
func (c *Ctx) SetGroupKick(rejectAddRequest bool) error {
	return c.Std.SetGroupKick(c.Event.GroupId, c.Event.UserId, rejectAddRequest)
}

// SetGroupBan 群组 Event.GroupId 单人禁言 Event.UserId
func (c *Ctx) SetGroupBan(duration time.Duration) error {
	return c.Std.SetGroupBan(c.Event.GroupId, c.Event.UserId, duration)
}

// SetGroupAnonymousBan 群组 Event.GroupId 匿名用户禁言
func (c *Ctx) SetGroupAnonymousBan(duration time.Duration) error {
	return c.Std.SetGroupAnonymousBan(c.Event.GroupId, c.Event.Anonymous, c.Event.Flag, duration)
}

// SetGroupWholeBan 群组 Event.GroupId 全员禁言
func (c *Ctx) SetGroupWholeBan(enable bool) error {
	return c.Std.SetGroupWholeBan(c.Event.GroupId, enable)
}

// SetGroupAdmin 群组 Event.GroupId 设置管理员 Event.UserId
func (c *Ctx) SetGroupAdmin(enable bool) error {
	return c.Std.SetGroupAdmin(c.Event.GroupId, c.Event.UserId, enable)
}

// SetGroupAnonymous 群组 Event.GroupId 匿名
func (c *Ctx) SetGroupAnonymous(enable bool) error {
	return c.Std.SetGroupAnonymous(c.Event.GroupId, enable)
}

// SetGroupCard 设置 Event.GroupId Event.UserId 群名片（群备注）
func (c *Ctx) SetGroupCard(card string) error {
	return c.Std.SetGroupCard(c.Event.GroupId, c.Event.UserId, card)
}

// SetGroupName 设置群 Event.GroupId 名
func (c *Ctx) SetGroupName(name string) error {
	return c.Std.SetGroupName(c.Event.GroupId, name)
}

// SetGroupLeave 退出群组 Event.GroupId
func (c *Ctx) SetGroupLeave(isDismiss bool) error {
	return c.Std.SetGroupLeave(c.Event.GroupId, isDismiss)
}

// SetGroupSpecialTitle 设置群组 Event.GroupId c.Event.UserId 专属头衔
func (c *Ctx) SetGroupSpecialTitle(specialTitle string, duration time.Duration) error {
	return c.Std.SetGroupSpecialTitle(c.Event.GroupId, c.Event.UserId, specialTitle, duration)
}

// SetFriendAddRequest 处理加好友请求 Event.Flag
func (c *Ctx) SetFriendAddRequest(approve bool, remark string) error {
	return c.Std.SetFriendAddRequest(c.Event.Flag, approve, remark)
}

// SetGroupAddRequest 处理加群请求／邀请 Event.Flag Event.SubType
func (c *Ctx) SetGroupAddRequest(approve bool, reason string) error {
	return c.Std.SetGroupAddRequest(c.Event.Flag, c.Event.SubType, approve, reason)
}

// GetLoginInfo 获取登录号信息
func (c *Ctx) GetLoginInfo() (*api.GetLoginInfoResp, error) {
	return c.Std.GetLoginInfo()
}

// GetStrangerInfo 获取陌生人 Event.UserId 信息
func (c *Ctx) GetStrangerInfo(noCache bool) (*api.GetStrangerInfoResp, error) {
	return c.Std.GetStrangerInfo(c.Event.UserId, noCache)
}

// GetFriendList 获取好友列表
func (c *Ctx) GetFriendList() (api.GetFriendListResp, error) {
	return c.Std.GetFriendList()
}

// GetGroupInfo 获取群 Event.GroupId 信息
func (c *Ctx) GetGroupInfo(noCache bool) (*api.GetGroupInfoResp, error) {
	return c.Std.GetGroupInfo(c.Event.GroupId, noCache)
}

// GetGroupList 获取群列表
func (c *Ctx) GetGroupList() (api.GetGroupListResp, error) {
	return c.Std.GetGroupList()
}

// GetGroupMemberInfo 获取群 Event.GroupId 成员 Event.UserId 信息
func (c *Ctx) GetGroupMemberInfo(noCache bool) (*api.GetGroupMemberInfoResp, error) {
	return c.Std.GetGroupMemberInfo(c.Event.GroupId, c.Event.UserId, noCache)
}

// GetGroupMemberList 获取群 Event.GroupId 成员列表
func (c *Ctx) GetGroupMemberList() (api.GetGroupMemberListResp, error) {
	return c.Std.GetGroupMemberList(c.Event.GroupId)
}

// GetGroupHonorInfo 获取群 Event.GroupId 荣誉信息
func (c *Ctx) GetGroupHonorInfo(typ string) (*api.GetGroupHonorInfoResp, error) {
	return c.Std.GetGroupHonorInfo(c.Event.GroupId, typ)
}

// GetCookies 获取 Cookies
func (c *Ctx) GetCookies(domain string) (*api.GetCookiesResp, error) {
	return c.Std.GetCookies(domain)
}

// GetCsrfToken 获取 CSRF Token
func (c *Ctx) GetCsrfToken() (*api.GetCsrfTokenResp, error) {
	return c.Std.GetCsrfToken()
}

// GetCredentials 获取 QQ 相关接口凭证
func (c *Ctx) GetCredentials(domain string) (*api.GetCredentialsResp, error) {
	return c.Std.GetCredentials(domain)
}

// GetRecord 获取语音
func (c *Ctx) GetRecord(outFormat string) (*api.GetRecordResp, error) {
	panic("todo: not implemented")
	// var file string // todo: parse from c.Event.RawMessage
	// return c.Std.GetRecord(file, outFormat)
}

// GetImage 获取图片
func (c *Ctx) GetImage() (*api.GetImageResp, error) {
	panic("todo: not implemented")
	// var file string // todo: parse from c.Event.RawMessage
	// return c.Std.GetImage(file)
}

// CanSendImage 是否可以发送图片
func (c *Ctx) CanSendImage() (*api.CanSendImageResp, error) {
	return c.Std.CanSendImage()
}

// CanSendRecord 是否可以发送语音
func (c *Ctx) CanSendRecord() (*api.CanSendRecordResp, error) {
	return c.Std.CanSendRecord()
}

// GetStatus 获取运行状态
func (c *Ctx) GetStatus() (*api.GetStatusResp, error) {
	return c.Std.GetStatus()
}

// GetVersionInfo 获取版本信息
func (c *Ctx) GetVersionInfo() (*api.GetVersionInfoResp, error) {
	return c.Std.GetVersionInfo()
}

// SetRestart 重启 OneBot 实现
func (c *Ctx) SetRestart(delay time.Duration) error {
	return c.Std.SetRestart(delay)
}

// CleanCache 清理缓存
func (c *Ctx) CleanCache() error {
	return c.Std.CleanCache()
}
