package napcat

import (
	"encoding/json"
	"fmt"
)

// 本文件是 EasyOnebot 实际用到的 NapCat 端点的类型化封装。
//
// 方法名与 spec 的 operationId 对应, 想扩展只需要照着加一个方法;
// 参数若是可空的 (生成代码里是 *T), 用 Go 1.26+ 的 new(expr) 取指针, 例如 new("123")。

// BoxMessage 把任意消息输入 (string / 消息段 / 消息段数组) 转成生成代码里的消息联合类型。
//
// 生成代码没有导出联合类型的构造函数, 所以用一次 JSON 往返; 这条路与 NapCat
// 实际收发的 wire 格式完全一致, 不会引入额外的转换规则。
func BoxMessage(v any) (OB11MessageMixType, error) {
	var m OB11MessageMixType
	b, err := json.Marshal(v)
	if err != nil {
		return m, fmt.Errorf("napcat: 编码消息失败: %w", err)
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return m, fmt.Errorf("napcat: 消息不是合法的 OneBot 消息: %w", err)
	}
	return m, nil
}

// BoxTo 把普通值转成生成代码里的联合类型 (spec 里很多字段是 T|string 这种 anyOf)。
//
// 生成代码没有导出联合类型的构造函数, 统一用一次 JSON 往返构造, 与 BoxMessage 同路。
func BoxTo[T any, U any](v T) U {
	var out U
	b, err := json.Marshal(v)
	if err != nil {
		return out
	}
	_ = json.Unmarshal(b, &out)
	return out
}

// ---- 消息 ----

// SendGroupMsg 发送群聊消息
func (c Caller) SendGroupMsg(body SendGroupMsgJSONBody) (SendGroupMsgData, error) {
	return Call[SendGroupMsgJSONBody, SendGroupMsgData](c, "SendGroupMsg", body)
}

// SendPrivateMsg 发送私聊消息
func (c Caller) SendPrivateMsg(body SendPrivateMsgJSONBody) (SendPrivateMsgData, error) {
	return Call[SendPrivateMsgJSONBody, SendPrivateMsgData](c, "SendPrivateMsg", body)
}

// SendGroupForwardMsg 发送群合并转发
func (c Caller) SendGroupForwardMsg(body SendGroupForwardMsgJSONBody) (SendGroupForwardMsgData, error) {
	return Call[SendGroupForwardMsgJSONBody, SendGroupForwardMsgData](c, "SendGroupForwardMsg", body)
}

// SendPrivateForwardMsg 发送私聊合并转发
func (c Caller) SendPrivateForwardMsg(body SendPrivateForwardMsgJSONBody) (SendPrivateForwardMsgData, error) {
	return Call[SendPrivateForwardMsgJSONBody, SendPrivateForwardMsgData](c, "SendPrivateForwardMsg", body)
}

// DeleteMsg 撤回消息
func (c Caller) DeleteMsg(body DeleteMsgJSONBody) (DeleteMsgData, error) {
	return Call[DeleteMsgJSONBody, DeleteMsgData](c, "DeleteMsg", body)
}

// GetMsg 获取消息
func (c Caller) GetMsg(body GetMsgJSONBody) (GetMsgData, error) {
	return Call[GetMsgJSONBody, GetMsgData](c, "GetMsg", body)
}

// GetForwardMsg 获取合并转发内容
func (c Caller) GetForwardMsg(body GetForwardMsgJSONBody) (GetForwardMsgData, error) {
	return Call[GetForwardMsgJSONBody, GetForwardMsgData](c, "GetForwardMsg", body)
}

// ---- 信息 ----

// GetLoginInfo 获取登录号信息
func (c Caller) GetLoginInfo() (OB11User, error) {
	return Call[GetLoginInfoJSONBody, OB11User](c, "GetLoginInfo", GetLoginInfoJSONBody{})
}

// GetGroupList 获取群列表
func (c Caller) GetGroupList() (GetGroupListData, error) {
	return Call[GetGroupListJSONBody, GetGroupListData](c, "GetGroupList", GetGroupListJSONBody{})
}

// GetFriendList 获取好友列表
func (c Caller) GetFriendList() (GetFriendListData, error) {
	return Call[GetFriendListJSONBody, GetFriendListData](c, "GetFriendList", GetFriendListJSONBody{})
}

// GetGroupMemberInfo 获取群成员信息
func (c Caller) GetGroupMemberInfo(body GetGroupMemberInfoJSONBody) (OB11GroupMember, error) {
	return Call[GetGroupMemberInfoJSONBody, OB11GroupMember](c, "GetGroupMemberInfo", body)
}

// ---- 交互 ----

// FriendPoke 好友戳一戳
func (c Caller) FriendPoke(body FriendPokeJSONBody) (FriendPokeData, error) {
	return Call[FriendPokeJSONBody, FriendPokeData](c, "FriendPoke", body)
}

// GroupPoke 群戳一戳
func (c Caller) GroupPoke(body GroupPokeJSONBody) (GroupPokeData, error) {
	return Call[GroupPokeJSONBody, GroupPokeData](c, "GroupPoke", body)
}

// ---- 群管理 ----

// SetGroupBan 群组禁言
func (c Caller) SetGroupBan(body SetGroupBanJSONBody) (SetGroupBanData, error) {
	return Call[SetGroupBanJSONBody, SetGroupBanData](c, "SetGroupBan", body)
}

// SetGroupKick 群组踢人
func (c Caller) SetGroupKick(body SetGroupKickJSONBody) (SetGroupKickData, error) {
	return Call[SetGroupKickJSONBody, SetGroupKickData](c, "SetGroupKick", body)
}

// SetGroupAdmin 设置群管理员
func (c Caller) SetGroupAdmin(body SetGroupAdminJSONBody) (SetGroupAdminData, error) {
	return Call[SetGroupAdminJSONBody, SetGroupAdminData](c, "SetGroupAdmin", body)
}

// SetGroupWholeBan 全员禁言
func (c Caller) SetGroupWholeBan(body SetGroupWholeBanJSONBody) (SetGroupWholeBanData, error) {
	return Call[SetGroupWholeBanJSONBody, SetGroupWholeBanData](c, "SetGroupWholeBan", body)
}

// SetGroupCard 设置群名片
func (c Caller) SetGroupCard(body SetGroupCardJSONBody) (SetGroupCardData, error) {
	return Call[SetGroupCardJSONBody, SetGroupCardData](c, "SetGroupCard", body)
}

// SetGroupName 设置群名称
func (c Caller) SetGroupName(body SetGroupNameJSONBody) (SetGroupNameData, error) {
	return Call[SetGroupNameJSONBody, SetGroupNameData](c, "SetGroupName", body)
}

// SetGroupLeave 退出群组
func (c Caller) SetGroupLeave(body SetGroupLeaveJSONBody) (SetGroupLeaveData, error) {
	return Call[SetGroupLeaveJSONBody, SetGroupLeaveData](c, "SetGroupLeave", body)
}

// SetGroupSpecialTitle 设置专属头衔
func (c Caller) SetGroupSpecialTitle(body SetGroupSpecialTitleJSONBody) (SetGroupSpecialTitleData, error) {
	return Call[SetGroupSpecialTitleJSONBody, SetGroupSpecialTitleData](c, "SetGroupSpecialTitle", body)
}

// GetGroupMemberList 获取群成员列表
func (c Caller) GetGroupMemberList(body GetGroupMemberListJSONBody) (GetGroupMemberListData, error) {
	return Call[GetGroupMemberListJSONBody, GetGroupMemberListData](c, "GetGroupMemberList", body)
}

// GetGroupInfo 获取群信息
func (c Caller) GetGroupInfo(body GetGroupInfoJSONBody) (OB11Group, error) {
	return Call[GetGroupInfoJSONBody, OB11Group](c, "GetGroupInfo", body)
}

// ---- 系统 / 能力 ----

// GetStatus 获取运行状态
func (c Caller) GetStatus(body GetStatusJSONBody) (GetStatusData, error) {
	return Call[GetStatusJSONBody, GetStatusData](c, "GetStatus", body)
}

// GetVersionInfo 获取版本信息
func (c Caller) GetVersionInfo(body GetVersionInfoJSONBody) (GetVersionInfoData, error) {
	return Call[GetVersionInfoJSONBody, GetVersionInfoData](c, "GetVersionInfo", body)
}

// CanSendImage 是否可以发送图片
func (c Caller) CanSendImage(body CanSendImageJSONBody) (CanSendImageData, error) {
	return Call[CanSendImageJSONBody, CanSendImageData](c, "CanSendImage", body)
}

// CanSendRecord 是否可以发送语音
func (c Caller) CanSendRecord(body CanSendRecordJSONBody) (CanSendRecordData, error) {
	return Call[CanSendRecordJSONBody, CanSendRecordData](c, "CanSendRecord", body)
}

// ---- 其它常用 ----

// GetStrangerInfo 获取陌生人信息
func (c Caller) GetStrangerInfo(body GetStrangerInfoJSONBody) (GetStrangerInfoData, error) {
	return Call[GetStrangerInfoJSONBody, GetStrangerInfoData](c, "GetStrangerInfo", body)
}

// GetImage 获取图片
func (c Caller) GetImage(body GetImageJSONBody) (GetImageData, error) {
	return Call[GetImageJSONBody, GetImageData](c, "GetImage", body)
}

// GetRecord 获取语音
func (c Caller) GetRecord(body GetRecordJSONBody) (GetRecordData, error) {
	return Call[GetRecordJSONBody, GetRecordData](c, "GetRecord", body)
}

// SendLike 点赞
func (c Caller) SendLike(body SendLikeJSONBody) (SendLikeData, error) {
	return Call[SendLikeJSONBody, SendLikeData](c, "SendLike", body)
}

// CleanCache 清理缓存
func (c Caller) CleanCache(body CleanCacheJSONBody) (CleanCacheData, error) {
	return Call[CleanCacheJSONBody, CleanCacheData](c, "CleanCache", body)
}

// ---- 文件 / 历史 / 系统 ----

// GetGroupRootFiles 获取群根目录文件列表
func (c Caller) GetGroupRootFiles(body GetGroupRootFilesJSONBody) (GetGroupRootFilesData, error) {
	return Call[GetGroupRootFilesJSONBody, GetGroupRootFilesData](c, "GetGroupRootFiles", body)
}

// GetGroupFilesByFolder 获取群文件夹文件列表
func (c Caller) GetGroupFilesByFolder(body GetGroupFilesByFolderJSONBody) (GetGroupFilesByFolderData, error) {
	return Call[GetGroupFilesByFolderJSONBody, GetGroupFilesByFolderData](c, "GetGroupFilesByFolder", body)
}

// GetGroupFileURL 获取群文件URL
func (c Caller) GetGroupFileURL(body GetGroupFileURLJSONBody) (GetGroupFileURLData, error) {
	return Call[GetGroupFileURLJSONBody, GetGroupFileURLData](c, "GetGroupFileURL", body)
}

// UploadGroupFile 上传群文件
func (c Caller) UploadGroupFile(body UploadGroupFileJSONBody) (UploadGroupFileData, error) {
	return Call[UploadGroupFileJSONBody, UploadGroupFileData](c, "UploadGroupFile", body)
}

// UploadPrivateFile 上传私聊文件
func (c Caller) UploadPrivateFile(body UploadPrivateFileJSONBody) (UploadPrivateFileData, error) {
	return Call[UploadPrivateFileJSONBody, UploadPrivateFileData](c, "UploadPrivateFile", body)
}

// GetFriendMsgHistory 获取好友历史消息
func (c Caller) GetFriendMsgHistory(body GetFriendMsgHistoryJSONBody) (GetFriendMsgHistoryData, error) {
	return Call[GetFriendMsgHistoryJSONBody, GetFriendMsgHistoryData](c, "GetFriendMsgHistory", body)
}

// GetGroupMsgHistory 获取群历史消息
func (c Caller) GetGroupMsgHistory(body GetGroupMsgHistoryJSONBody) (GetGroupMsgHistoryData, error) {
	return Call[GetGroupMsgHistoryJSONBody, GetGroupMsgHistoryData](c, "GetGroupMsgHistory", body)
}

// SetRestart 重启服务
func (c Caller) SetRestart(body SetRestartJSONBody) (SetRestartData, error) {
	return Call[SetRestartJSONBody, SetRestartData](c, "SetRestart", body)
}

// GetCookies 获取 Cookies
func (c Caller) GetCookies(body GetCookiesJSONBody) (GetCookiesData, error) {
	return Call[GetCookiesJSONBody, GetCookiesData](c, "GetCookies", body)
}

// GetCsrfToken 获取 CSRF Token
func (c Caller) GetCsrfToken(body GetCsrfTokenJSONBody) (GetCsrfTokenData, error) {
	return Call[GetCsrfTokenJSONBody, GetCsrfTokenData](c, "GetCsrfToken", body)
}

// GetCredentials 获取登录凭证
func (c Caller) GetCredentials(body GetCredentialsJSONBody) (GetCredentialsData, error) {
	return Call[GetCredentialsJSONBody, GetCredentialsData](c, "GetCredentials", body)
}

// ---- 请求处理 / 荣誉 ----

// SetFriendAddRequest 处理加好友请求
func (c Caller) SetFriendAddRequest(body SetFriendAddRequestJSONBody) (SetFriendAddRequestData, error) {
	return Call[SetFriendAddRequestJSONBody, SetFriendAddRequestData](c, "SetFriendAddRequest", body)
}

// SetGroupAddRequest 处理加群请求 / 邀请
func (c Caller) SetGroupAddRequest(body SetGroupAddRequestJSONBody) (SetGroupAddRequestData, error) {
	return Call[SetGroupAddRequestJSONBody, SetGroupAddRequestData](c, "SetGroupAddRequest", body)
}

// GetGroupHonorInfo 获取群荣誉信息
func (c Caller) GetGroupHonorInfo(body GetGroupHonorInfoJSONBody) (GetGroupHonorInfoData, error) {
	return Call[GetGroupHonorInfoJSONBody, GetGroupHonorInfoData](c, "GetGroupHonorInfo", body)
}
