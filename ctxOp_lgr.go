package EasyOnebot

import (
	"github.com/Miuzarte/EasyOnebot/api"
	"github.com/Miuzarte/EasyOnebot/message"
)

/*
	# TODO: complete
*/

// GetFriendMsgHistory 获取好友历史消息记录
func (c *Ctx) GetFriendMsgHistory(count int) (*api.GetFriendMsgHistoryResp, error) {
	return c.Lgr.GetFriendMsgHistory(c.Event.UserId, c.Event.MessageId, count)
}

// GetGroupMsgHistory 获取群组历史消息记录
func (c *Ctx) GetGroupMsgHistory(count int) (*api.GetGroupMsgHistoryResp, error) {
	return c.Lgr.GetGroupMsgHistory(c.Event.GroupId, c.Event.MessageId, count)
}

// SendForwardMsg 构造合并转发消息
func (c *Ctx) SendForwardMsg(messages message.SegmentArray) (api.SendForwardMsg_LgrResp, error) {
	return c.Lgr.SendForwardMsg(messages)
}

// SendGroupForwardMsg 发送合并转发 (群聊)
func (c *Ctx) SendGroupForwardMsg(messages message.SegmentArray) (*api.SendPrivateForwardMsgResp, error) {
	return c.Lgr.SendGroupForwardMsg(c.Event.GroupId, messages)
}

// SendPrivateForwardMsg 发送合并转发 (好友)
func (c *Ctx) SendPrivateForwardMsg(messages message.SegmentArray) (*api.SendGroupForwardMsgResp, error) {
	return c.Lgr.SendPrivateForwardMsg(c.Event.UserId, messages)
}

// UploadGroupFile 上传群文件
func (c *Ctx) UploadGroupFile(file, name, folder string) error {
	return c.Lgr.UploadGroupFile(c.Event.GroupId, file, name, folder)
}

// UploadPrivateFile 私聊发送文件
func (c *Ctx) UploadPrivateFile(file, name string) error {
	return c.Lgr.UploadPrivateFile(c.Event.UserId, file, name)
}

// GetGroupRootFiles 获取群根目录文件列表
func (c *Ctx) GetGroupRootFiles() (*api.GetGroupRootFilesResp, error) {
	return c.Lgr.GetGroupRootFiles(c.Event.GroupId)
}

// GetGroupFilesByFolder 获取群子目录文件列表
func (c *Ctx) GetGroupFilesByFolder(folderId string) (*api.GetGroupFilesByFolderResp, error) {
	return c.Lgr.GetGroupFilesByFolder(c.Event.GroupId, folderId)
}

// GetGroupFileUrl 获取群文件资源链接
func (c *Ctx) GetGroupFileUrl(fileId, busid string) (*api.GetGroupFileUrlResp, error) {
	return c.Lgr.GetGroupFileUrl(c.Event.GroupId, fileId, busid)
}

// FriendPoke 好友戳一戳
func (c *Ctx) FriendPoke() error {
	return c.Lgr.FriendPoke(c.Event.UserId)
}

// GroupPoke 群组戳一戳
func (c *Ctx) GroupPoke() error {
	return c.Lgr.GroupPoke(c.Event.GroupId, c.Event.UserId)
}
