package EasyOnebot

import (
	"fmt"
	"time"

	"github.com/Miuzarte/EasyOnebot/api"
	"github.com/Miuzarte/EasyOnebot/api/napcat"
	"github.com/Miuzarte/EasyOnebot/message"
)

// 本文件是 ctx 在 Event 上下文上的一键操作。
//
// 能表达成 spec 类型的调用都走 [Bot.NapCat] (严格解码);
// 仍用 api.* 响应的两类是刻意的:
//   - get_msg / get_forward_msg: 生成类型里 message/messages 是 []interface{},
//     而 api 侧已经解成 message.SegmentArray, 对业务更有用
//   - get_cookies / get_csrf_token / get_credentials 等低频端点暂未迁移

// GetMsg 获取 Event.MessageId
func (c *Ctx) GetMsg() (*api.GetMsgResp, error) {
	resp, err, _ := c.Std.GetMsgSf(c.Event.MessageId)
	// resp.Message.TryAtoi()
	return resp, err
}

// GetForwardMsg 获取合并转发消息 Event.MessageId
//
// 生成类型里 messages 只是 []interface{}, 这里沿用 api 侧的 SegmentArray 解析。
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
	body := napcat.SendLikeJSONBody{
		UserID: fmt.Sprint(c.Event.UserId),
		Times:  napcat.BoxTo[int, napcat.SendLikeJSONBody_Times](times),
	}
	_, err := c.Bot.NapCat().SendLike(body)
	return err
}

// SetGroupKick 群组 Event.GroupId 踢 Event.UserId
func (c *Ctx) SetGroupKick(rejectAddRequest bool) error {
	body := napcat.SetGroupKickJSONBody{
		GroupID:          fmt.Sprint(c.Event.GroupId),
		UserID:           fmt.Sprint(c.Event.UserId),
		RejectAddRequest: new(napcat.BoxTo[bool, napcat.SetGroupKickJSONBody_RejectAddRequest](rejectAddRequest)),
	}
	_, err := c.Bot.NapCat().SetGroupKick(body)
	return err
}

// SetGroupBan 群组 Event.GroupId 单人禁言 Event.UserId
func (c *Ctx) SetGroupBan(duration time.Duration) error {
	body := napcat.SetGroupBanJSONBody{
		GroupID:  fmt.Sprint(c.Event.GroupId),
		UserID:   fmt.Sprint(c.Event.UserId),
		Duration: napcat.BoxTo[int64, napcat.SetGroupBanJSONBody_Duration](int64(duration / time.Second)),
	}
	_, err := c.Bot.NapCat().SetGroupBan(body)
	return err
}

// SetGroupWholeBan 群组 Event.GroupId 全员禁言
func (c *Ctx) SetGroupWholeBan(enable bool) error {
	body := napcat.SetGroupWholeBanJSONBody{
		GroupID: fmt.Sprint(c.Event.GroupId),
		Enable:  new(napcat.BoxTo[bool, napcat.SetGroupWholeBanJSONBody_Enable](enable)),
	}
	_, err := c.Bot.NapCat().SetGroupWholeBan(body)
	return err
}

// SetGroupAdmin 群组 Event.GroupId 设置管理员 Event.UserId
func (c *Ctx) SetGroupAdmin(enable bool) error {
	body := napcat.SetGroupAdminJSONBody{
		GroupID: fmt.Sprint(c.Event.GroupId),
		UserID:  fmt.Sprint(c.Event.UserId),
		Enable:  new(napcat.BoxTo[bool, napcat.SetGroupAdminJSONBody_Enable](enable)),
	}
	_, err := c.Bot.NapCat().SetGroupAdmin(body)
	return err
}

// SetGroupCard 设置 Event.GroupId Event.UserId 群名片 (群备注)
func (c *Ctx) SetGroupCard(card string) error {
	body := napcat.SetGroupCardJSONBody{
		GroupID: fmt.Sprint(c.Event.GroupId),
		UserID:  fmt.Sprint(c.Event.UserId),
		Card:    new(card),
	}
	_, err := c.Bot.NapCat().SetGroupCard(body)
	return err
}

// SetGroupName 设置群 Event.GroupId 名
func (c *Ctx) SetGroupName(name string) error {
	body := napcat.SetGroupNameJSONBody{
		GroupID:   fmt.Sprint(c.Event.GroupId),
		GroupName: name,
	}
	_, err := c.Bot.NapCat().SetGroupName(body)
	return err
}

// SetGroupLeave 退出群组 Event.GroupId
func (c *Ctx) SetGroupLeave(isDismiss bool) error {
	body := napcat.SetGroupLeaveJSONBody{
		GroupID:   fmt.Sprint(c.Event.GroupId),
		IsDismiss: new(napcat.BoxTo[bool, napcat.SetGroupLeaveJSONBody_IsDismiss](isDismiss)),
	}
	_, err := c.Bot.NapCat().SetGroupLeave(body)
	return err
}

// SetGroupSpecialTitle 设置群组 Event.GroupId Event.UserId 专属头衔
func (c *Ctx) SetGroupSpecialTitle(specialTitle string, duration time.Duration) error {
	body := napcat.SetGroupSpecialTitleJSONBody{
		GroupID:      fmt.Sprint(c.Event.GroupId),
		UserID:       fmt.Sprint(c.Event.UserId),
		SpecialTitle: specialTitle,
	}
	_, err := c.Bot.NapCat().SetGroupSpecialTitle(body)
	return err
}

// SetFriendAddRequest 处理加好友请求 Event.Flag
func (c *Ctx) SetFriendAddRequest(approve bool, remark string) error {
	body := napcat.SetFriendAddRequestJSONBody{
		Flag:    c.Event.Flag,
		Approve: new(napcat.BoxTo[bool, napcat.SetFriendAddRequestJSONBody_Approve](approve)),
		Remark:  new(remark),
	}
	_, err := c.Bot.NapCat().SetFriendAddRequest(body)
	return err
}

// SetGroupAddRequest 处理加群请求 / 邀请 Event.Flag Event.SubType
func (c *Ctx) SetGroupAddRequest(approve bool, reason string) error {
	body := napcat.SetGroupAddRequestJSONBody{
		Flag:    c.Event.Flag,
		Approve: new(napcat.BoxTo[bool, napcat.SetGroupAddRequestJSONBody_Approve](approve)),
		Reason:  new(reason),
	}
	_, err := c.Bot.NapCat().SetGroupAddRequest(body)
	return err
}

// GetLoginInfo 获取登录号信息
func (c *Ctx) GetLoginInfo() (napcat.OB11User, error) {
	return c.Bot.NapCat().GetLoginInfo()
}

// GetStrangerInfo 获取陌生人 Event.UserId 信息
func (c *Ctx) GetStrangerInfo(noCache bool) (napcat.GetStrangerInfoData, error) {
	body := napcat.GetStrangerInfoJSONBody{
		UserID:  fmt.Sprint(c.Event.UserId),
		NoCache: napcat.BoxTo[bool, napcat.GetStrangerInfoJSONBody_NoCache](noCache),
	}
	return c.Bot.NapCat().GetStrangerInfo(body)
}

// GetFriendList 获取好友列表
func (c *Ctx) GetFriendList() (napcat.GetFriendListData, error) {
	return c.Bot.NapCat().GetFriendList()
}

// GetGroupInfo 获取群 Event.GroupId 信息
func (c *Ctx) GetGroupInfo(noCache bool) (napcat.OB11Group, error) {
	_ = noCache // NapCat 的 get_group_info 没有 no_cache 参数
	body := napcat.GetGroupInfoJSONBody{GroupID: fmt.Sprint(c.Event.GroupId)}
	return c.Bot.NapCat().GetGroupInfo(body)
}

// GetGroupList 获取群列表
func (c *Ctx) GetGroupList() (napcat.GetGroupListData, error) {
	return c.Bot.NapCat().GetGroupList()
}

// GetGroupMemberInfo 获取群 Event.GroupId 成员 Event.UserId 信息
func (c *Ctx) GetGroupMemberInfo(noCache bool) (napcat.OB11GroupMember, error) {
	body := napcat.GetGroupMemberInfoJSONBody{
		GroupID: fmt.Sprint(c.Event.GroupId),
		UserID:  fmt.Sprint(c.Event.UserId),
	}
	if noCache {
		body.NoCache = new(napcat.BoxTo[bool, napcat.GetGroupMemberInfoJSONBody_NoCache](true))
	}
	return c.Bot.NapCat().GetGroupMemberInfo(body)
}

// GetGroupMemberList 获取群 Event.GroupId 成员列表
func (c *Ctx) GetGroupMemberList() (napcat.GetGroupMemberListData, error) {
	return c.Bot.NapCat().GetGroupMemberList(napcat.GetGroupMemberListJSONBody{GroupID: fmt.Sprint(c.Event.GroupId)})
}

// GetGroupHonorInfo 获取群 Event.GroupId 荣誉信息
func (c *Ctx) GetGroupHonorInfo(typ string) (napcat.GetGroupHonorInfoData, error) {
	body := napcat.GetGroupHonorInfoJSONBody{
		GroupID: fmt.Sprint(c.Event.GroupId),
		Type:    new(napcat.BoxTo[string, napcat.GetGroupHonorInfoJSONBodyType](typ)),
	}
	return c.Bot.NapCat().GetGroupHonorInfo(body)
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

// GetRecord 获取语音 (未实现)
func (c *Ctx) GetRecord(outFormat string) (*api.GetRecordResp, error) {
	panic("todo: not implemented")
}

// GetImage 获取图片 (未实现)
func (c *Ctx) GetImage() (*api.GetImageResp, error) {
	panic("todo: not implemented")
}

// CanSendImage 是否可以发送图片
func (c *Ctx) CanSendImage() (napcat.CanSendImageData, error) {
	return c.Bot.NapCat().CanSendImage(napcat.CanSendImageJSONBody{})
}

// CanSendRecord 是否可以发送语音
func (c *Ctx) CanSendRecord() (napcat.CanSendRecordData, error) {
	return c.Bot.NapCat().CanSendRecord(napcat.CanSendRecordJSONBody{})
}

// GetStatus 获取运行状态
func (c *Ctx) GetStatus() (napcat.GetStatusData, error) {
	return c.Bot.NapCat().GetStatus(napcat.GetStatusJSONBody{})
}

// GetVersionInfo 获取版本信息
func (c *Ctx) GetVersionInfo() (napcat.GetVersionInfoData, error) {
	return c.Bot.NapCat().GetVersionInfo(napcat.GetVersionInfoJSONBody{})
}

// SetRestart 重启 OneBot 实现
func (c *Ctx) SetRestart(delay time.Duration) error {
	return c.Std.SetRestart(delay)
}

// CleanCache 清理缓存
func (c *Ctx) CleanCache() error {
	_, err := c.Bot.NapCat().CleanCache(napcat.CleanCacheJSONBody{})
	return err
}
