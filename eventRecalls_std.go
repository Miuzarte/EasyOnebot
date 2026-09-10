package EasyOnebot

import (
	"runtime/debug"

	"github.com/Miuzarte/EasyOnebot/event"
	"github.com/Miuzarte/EasyOnebot/internal/utils"
)

// eventRecalls 遍历调用时不是异步的
type eventRecalls struct {
	messagePrivate []func(*event.MessagePrivate) // 私聊消息
	messageGroup   []func(*event.MessageGroup)   // 群消息

	noticeGroupUpload     []func(*event.NoticeGroupUpload)     // 群文件上传
	noticeGroupAdmin      []func(*event.NoticeGroupAdmin)      // 群管理员变动
	noticeGroupDecrease   []func(*event.NoticeGroupDecrease)   // 群成员减少
	noticeGroupIncrease   []func(*event.NoticeGroupIncrease)   // 群成员增加
	noticeGroupBan        []func(*event.NoticeGroupBan)        // 群成员禁言
	noticeFriendAdd       []func(*event.NoticeFriendAdd)       // 好友添加
	noticeGroupRecall     []func(*event.NoticeGroupRecall)     // 群消息撤回
	noticeFriendRecall    []func(*event.NoticeFriendRecall)    // 好友消息撤回
	noticeNotifyPoke      []func(*event.NoticeNotifyPoke)      // 群内戳一戳
	noticeNotifyLuckyKing []func(*event.NoticeNotifyLuckyKing) // 群红包运气王
	noticeNotifyHonor     []func(*event.NoticeNotifyHonor)     // 群成员荣誉变更

	requestFriend []func(*event.RequestFriend) // 好友请求
	requestGroup  []func(*event.RequestGroup)  // 加群请求／邀请

	metaEventLifecycle []func(*event.MetaEventLifecycle) // 生命周期
	metaEventHeartbeat []func(*event.MetaEventHeartbeat) // 心跳包

	eventRecalls_Nc
}

// OnMessagePrivate 处理私聊消息
func (b *Bot) OnMessagePrivate(handler func(*event.MessagePrivate)) *Bot {
	b.eventRecalls.messagePrivate = utils.OnDemandAppend(b.eventRecalls.messagePrivate, handler)
	return b
}

// OnMessageGroup 处理群消息
func (b *Bot) OnMessageGroup(handler func(*event.MessageGroup)) *Bot {
	b.eventRecalls.messageGroup = utils.OnDemandAppend(b.eventRecalls.messageGroup, handler)
	return b
}

// OnNoticeGroupUpload 处理群文件上传
func (b *Bot) OnNoticeGroupUpload(handler func(*event.NoticeGroupUpload)) *Bot {
	b.eventRecalls.noticeGroupUpload = utils.OnDemandAppend(b.eventRecalls.noticeGroupUpload, handler)
	return b
}

// OnNoticeGroupAdmin 处理群管理员变动
func (b *Bot) OnNoticeGroupAdmin(handler func(*event.NoticeGroupAdmin)) *Bot {
	b.eventRecalls.noticeGroupAdmin = utils.OnDemandAppend(b.eventRecalls.noticeGroupAdmin, handler)
	return b
}

// OnNoticeGroupDecrease 处理群成员减少
func (b *Bot) OnNoticeGroupDecrease(handler func(*event.NoticeGroupDecrease)) *Bot {
	b.eventRecalls.noticeGroupDecrease = utils.OnDemandAppend(b.eventRecalls.noticeGroupDecrease, handler)
	return b
}

// OnNoticeGroupIncrease 处理群成员增加
func (b *Bot) OnNoticeGroupIncrease(handler func(*event.NoticeGroupIncrease)) *Bot {
	b.eventRecalls.noticeGroupIncrease = utils.OnDemandAppend(b.eventRecalls.noticeGroupIncrease, handler)
	return b
}

// OnNoticeGroupBan 处理群成员禁言
func (b *Bot) OnNoticeGroupBan(handler func(*event.NoticeGroupBan)) *Bot {
	b.eventRecalls.noticeGroupBan = utils.OnDemandAppend(b.eventRecalls.noticeGroupBan, handler)
	return b
}

// OnNoticeFriendAdd 处理好友添加
func (b *Bot) OnNoticeFriendAdd(handler func(*event.NoticeFriendAdd)) *Bot {
	b.eventRecalls.noticeFriendAdd = utils.OnDemandAppend(b.eventRecalls.noticeFriendAdd, handler)
	return b
}

// OnNoticeGroupRecall 处理群消息撤回
func (b *Bot) OnNoticeGroupRecall(handler func(*event.NoticeGroupRecall)) *Bot {
	b.eventRecalls.noticeGroupRecall = utils.OnDemandAppend(b.eventRecalls.noticeGroupRecall, handler)
	return b
}

// OnNoticeFriendRecall 处理好友消息撤回
func (b *Bot) OnNoticeFriendRecall(handler func(*event.NoticeFriendRecall)) *Bot {
	b.eventRecalls.noticeFriendRecall = utils.OnDemandAppend(b.eventRecalls.noticeFriendRecall, handler)
	return b
}

// OnNoticeNotifyPoke 处理群内戳一戳
func (b *Bot) OnNoticeNotifyPoke(handler func(*event.NoticeNotifyPoke)) *Bot {
	b.eventRecalls.noticeNotifyPoke = utils.OnDemandAppend(b.eventRecalls.noticeNotifyPoke, handler)
	return b
}

// OnNoticeNotifyLuckyKing 处理群红包运气王
func (b *Bot) OnNoticeNotifyLuckyKing(handler func(*event.NoticeNotifyLuckyKing)) *Bot {
	b.eventRecalls.noticeNotifyLuckyKing = utils.OnDemandAppend(b.eventRecalls.noticeNotifyLuckyKing, handler)
	return b
}

// OnNoticeNotifyHonor 处理群成员荣誉变更
func (b *Bot) OnNoticeNotifyHonor(handler func(*event.NoticeNotifyHonor)) *Bot {
	b.eventRecalls.noticeNotifyHonor = utils.OnDemandAppend(b.eventRecalls.noticeNotifyHonor, handler)
	return b
}

// OnRequestFriend 处理好友请求
func (b *Bot) OnRequestFriend(handler func(*event.RequestFriend)) *Bot {
	b.eventRecalls.requestFriend = utils.OnDemandAppend(b.eventRecalls.requestFriend, handler)
	return b
}

// OnRequestGroup 处理群组请求
func (b *Bot) OnRequestGroup(handler func(*event.RequestGroup)) *Bot {
	b.eventRecalls.requestGroup = utils.OnDemandAppend(b.eventRecalls.requestGroup, handler)
	return b
}

// OnMetaEventLifecycle 处理生命周期事件
func (b *Bot) OnMetaEventLifecycle(handler func(*event.MetaEventLifecycle)) *Bot {
	b.eventRecalls.metaEventLifecycle = utils.OnDemandAppend(b.eventRecalls.metaEventLifecycle, handler)
	return b
}

// OnMetaEventHeartbeat 处理心跳包
func (b *Bot) OnMetaEventHeartbeat(handler func(*event.MetaEventHeartbeat)) *Bot {
	b.eventRecalls.metaEventHeartbeat = utils.OnDemandAppend(b.eventRecalls.metaEventHeartbeat, handler)
	return b
}

func (b *Bot) doEventRecallInternal(typedEvent any) {
	switch e := typedEvent.(type) {
	case *event.MessagePrivate:
		b.redisSetMessagePrivate(e)
	case *event.MessageGroup:
		b.redisSetMessageGroup(e)
		b.redisSetGroupCardName(e.GroupId, e.Sender.UserId, e.Sender.Card)
	case *event.MetaEventLifecycle:
		b.selfId = e.SelfId
	case *event.MetaEventHeartbeat:
		if !b.hb.running {
			b.hb.running = true
			go b.heartbeatLoop()
		}
		b.hb.signal <- e.Interval
	}
}

// maybeTODO: 反射实现, 性能降低
func (b *Bot) doEventRecall(typedEvent any) {
	i := new(int)
	defer func() {
		if err := recover(); err != nil {
			b.log.Error().Msgf("panic while handling %T[%d]: %v", typedEvent, *i, err)
			debug.PrintStack()
		}
	}()
	switch e := typedEvent.(type) {
	case *event.MessagePrivate:
		for j, f := range b.eventRecalls.messagePrivate {
			*i = j
			f(e)
		}
	case *event.MessageGroup:
		for j, f := range b.eventRecalls.messageGroup {
			*i = j
			f(e)
		}
	case *event.NoticeGroupUpload:
		for j, f := range b.eventRecalls.noticeGroupUpload {
			*i = j
			f(e)
		}
	case *event.NoticeGroupAdmin:
		for j, f := range b.eventRecalls.noticeGroupAdmin {
			*i = j
			f(e)
		}
	case *event.NoticeGroupDecrease:
		for j, f := range b.eventRecalls.noticeGroupDecrease {
			*i = j
			f(e)
		}
	case *event.NoticeGroupIncrease:
		for j, f := range b.eventRecalls.noticeGroupIncrease {
			*i = j
			f(e)
		}
	case *event.NoticeGroupBan:
		for j, f := range b.eventRecalls.noticeGroupBan {
			*i = j
			f(e)
		}
	case *event.NoticeFriendAdd:
		for j, f := range b.eventRecalls.noticeFriendAdd {
			*i = j
			f(e)
		}
	case *event.NoticeGroupRecall:
		for j, f := range b.eventRecalls.noticeGroupRecall {
			*i = j
			f(e)
		}
	case *event.NoticeFriendRecall:
		for j, f := range b.eventRecalls.noticeFriendRecall {
			*i = j
			f(e)
		}
	case *event.NoticeNotifyPoke:
		for j, f := range b.eventRecalls.noticeNotifyPoke {
			*i = j
			f(e)
		}
	case *event.NoticeNotifyLuckyKing:
		for j, f := range b.eventRecalls.noticeNotifyLuckyKing {
			*i = j
			f(e)
		}
	case *event.NoticeNotifyHonor:
		for j, f := range b.eventRecalls.noticeNotifyHonor {
			*i = j
			f(e)
		}
	case *event.RequestFriend:
		for j, f := range b.eventRecalls.requestFriend {
			*i = j
			f(e)
		}
	case *event.RequestGroup:
		for j, f := range b.eventRecalls.requestGroup {
			*i = j
			f(e)
		}
	case *event.MetaEventLifecycle:
		for j, f := range b.eventRecalls.metaEventLifecycle {
			*i = j
			f(e)
		}
	case *event.MetaEventHeartbeat:
		for j, f := range b.eventRecalls.metaEventHeartbeat {
			*i = j
			f(e)
		}
	default:
		b.doEventRecall_Nc(typedEvent)
	}
}
