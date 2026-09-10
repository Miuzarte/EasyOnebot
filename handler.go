package EasyOnebot

import (
	"fmt"
	"os"
	"runtime/debug"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Miuzarte/EasyOnebot/api"
	"github.com/Miuzarte/EasyOnebot/event"
	"github.com/Miuzarte/EasyOnebot/internal/utils"
	"github.com/Miuzarte/EasyOnebot/message"
)

func (b *Bot) handleApiResp(r *api.Response) {
	b.statistics.apiResp++
	ch, ok := b.apiPool.get(r.Echo)
	if !ok {
		b.log.Warn().Msgf("收到未知/过期响应: %v", r.Echo)
		return
	}
	ch <- r
}

func (b *Bot) handleEvent(e *event.Event) {
	if e.UserId == e.SelfId { // 自己发的消息
		return
	}
	b.statistics.event++

	e.InitFields()
	go b.classifyEvent(e) // 辨别
	ctx := b.toCtx(e)
	go b.matchEvent(ctx) // 匹配
	b.log.Trace().Msgf("ctx: %+v", ctx)

	if e.TypeL2 == event.TYPE_L2_MESSAGE_GROUP && e.GroupId != 0 && len(ctx.Ats) != 0 {
		go func() {
			var err error
			for _, atId := range ctx.Ats {
				err = b.redisPushAt(e.GroupId, e.UserId, e.MessageId)
				if err != nil {
					break
				}
				if atId != 0 {
					err = b.redisPushAtBy(e.GroupId, atId, e.MessageId)
				}
				if err != nil {
					break
				}
			}
			if err != nil {
				b.log.Warn().Err(err).Msg("缓存 at 关系失败")
			}
		}()
	}

	if e.TypeL2 == event.TYPE_L2_NOTICE_GROUP_RECALL {
		err := b.redisPushRecall(e.GroupId, e.UserId, e.MessageId)
		if err != nil {
			b.log.Warn().Err(err).Msg("缓存群撤回失败")
		}
	}
}

func (b *Bot) classifyEvent(e *event.Event) {
	var typedEvent any

	switch e.TypeL1 {
	case event.TYPE_L1_MESSAGE:
		switch e.TypeL2 {
		case event.TYPE_L2_MESSAGE_PRIVATE:
			b.statistics.messagePrivate++
			mp := utils.AnyCopy[event.MessagePrivate](e)
			typedEvent = mp

		case event.TYPE_L2_MESSAGE_GROUP:
			b.statistics.messageGroup++
			mg := utils.AnyCopy[event.MessageGroup](e)
			typedEvent = mg

		}

	case event.TYPE_L1_REQUEST:
		switch e.TypeL2 {
		case event.TYPE_L2_REQUEST_FRIEND:
			b.statistics.requestFriend++
			rf := utils.AnyCopy[event.RequestFriend](e)
			typedEvent = rf

		case event.TYPE_L2_REQUEST_GROUP:
			b.statistics.requestGroup++
			rg := utils.AnyCopy[event.RequestGroup](e)
			typedEvent = rg

		}

	case event.TYPE_L1_NOTICE:
		switch e.TypeL2 {
		case event.TYPE_L2_NOTICE_GROUP_UPLOAD:
			b.statistics.noticeGroupUpload++
			ngu := utils.AnyCopy[event.NoticeGroupUpload](e)
			typedEvent = ngu

		case event.TYPE_L2_NOTICE_GROUP_ADMIN:
			b.statistics.noticeGroupAdmin++
			nga := utils.AnyCopy[event.NoticeGroupAdmin](e)
			typedEvent = nga

		case event.TYPE_L2_NOTICE_GROUP_DECREASE:
			b.statistics.noticeGroupDecrease++
			ngd := utils.AnyCopy[event.NoticeGroupDecrease](e)
			typedEvent = ngd

		case event.TYPE_L2_NOTICE_GROUP_INCREASE:
			b.statistics.noticeGroupIncrease++
			ngi := utils.AnyCopy[event.NoticeGroupIncrease](e)
			typedEvent = ngi

		case event.TYPE_L2_NOTICE_GROUP_BAN:
			b.statistics.noticeGroupBan++
			ngb := utils.AnyCopy[event.NoticeGroupBan](e)
			typedEvent = ngb

		case event.TYPE_L2_NOTICE_FRIEND_ADD:
			b.statistics.noticeFriendAdd++
			nfa := utils.AnyCopy[event.NoticeFriendAdd](e)
			typedEvent = nfa

		case event.TYPE_L2_NOTICE_GROUP_RECALL:
			b.statistics.noticeGroupRecall++
			ngr := utils.AnyCopy[event.NoticeGroupRecall](e)
			typedEvent = ngr

		case event.TYPE_L2_NOTICE_FRIEND_RECALL:
			b.statistics.noticeFriendRecall++
			nfr := utils.AnyCopy[event.NoticeFriendRecall](e)
			typedEvent = nfr

		case event.TYPE_L2_NOTICE_NOTIFY:
			switch e.TypeL3 {
			case event.TYPE_L3_NOTICE_NOTIFY_POKE:
				b.statistics.noticeNotifyPoke++
				nnp := utils.AnyCopy[event.NoticeNotifyPoke](e)
				normalizePoke(nnp)
				typedEvent = nnp

			case event.TYPE_L3_NOTICE_NOTIFY_LUCKY_KING:
				b.statistics.noticeNotifyLuckyKing++
				nnlk := utils.AnyCopy[event.NoticeNotifyLuckyKing](e)
				typedEvent = nnlk

			case event.TYPE_L3_NOTICE_NOTIFY_HONOR:
				b.statistics.noticeNotifyHonor++
				nnh := utils.AnyCopy[event.NoticeNotifyHonor](e)
				typedEvent = nnh
			}

		case event.TYPE_L2_NOTICE_BOT_OFFLINE: // NapCat 扩展
			b.statistics.noticeBotOffline++
			noff := utils.AnyCopy[event.NoticeBotOffline](e)
			typedEvent = noff

		case event.TYPE_L2_NOTICE_BOT_ONLINE: // NapCat 扩展
			b.statistics.noticeBotOnline++
			non := utils.AnyCopy[event.NoticeBotOnline](e)
			typedEvent = non

		case event.TYPE_L2_NOTICE_ESSENCE: // NapCat 扩展
			b.statistics.noticeEssence++
			ne := utils.AnyCopy[event.NoticeEssence](e)
			typedEvent = ne

		case event.TYPE_L2_NOTICE_GROUP_NAME: // NapCat 扩展
			b.statistics.noticeGroupName++
			ngn := utils.AnyCopy[event.NoticeGroupName](e)
			typedEvent = ngn

		case event.TYPE_L2_NOTICE_REACTION: // NapCat 扩展
			b.statistics.noticeReaction++
			nr := utils.AnyCopy[event.NoticeReaction](e)
			typedEvent = nr

		case event.TYPE_L2_NOTICE_OFFLINE_FILE: // NapCat 扩展
			b.statistics.noticeOfflineFile++
			nofff := utils.AnyCopy[event.NoticeOfflineFile](e)
			typedEvent = nofff

		}

	case event.TYPE_L1_METAEVENT:
		switch e.TypeL2 {
		case event.TYPE_L2_META_LIFECYCLE:
			b.statistics.metaEventLifecycle++
			mel := utils.AnyCopy[event.MetaEventLifecycle](e)
			typedEvent = mel

		case event.TYPE_L2_META_HEARTBEAT:
			b.statistics.metaEventHeartbeat++
			meh := utils.AnyCopy[event.MetaEventHeartbeat](e)
			typedEvent = meh

		}
	}

	if typedEvent == nil {
		b.log.Warn().Msgf("nil event, raw: %v", e.RawEvent)
	}

	go b.doEventRecallInternal(typedEvent)
	go b.doEventRecall(typedEvent)
	if b.settings.msgLogOut {
		switch e.TypeL2 {
		case event.TYPE_L2_META_HEARTBEAT:
			b.log.Trace().Msg(b.FormatEvent(typedEvent))
		case event.TYPE_L2_META_LIFECYCLE:
			b.log.Debug().Msg(b.FormatEvent(typedEvent))
		default:
			b.log.Info().Msg(b.FormatEvent(typedEvent))
		}
	}
}

func (b *Bot) toCtx(e *event.Event) (ctx *Ctx) {
	ctx = &Ctx{Event: e, MixCaller: b.Call(), Bot: b}

	ctx.FilteredStranger = slices.Contains(b.filter.strangers, e.UserId)
	ctx.FilteredGroup = slices.Contains(b.filter.groups, e.GroupId)
	ctx.IsSuperuser = slices.Contains(b.sus, e.UserId)

	if l := len(e.RawMessage); l > 0 {
		// json 转义影响正则匹配
		ctx.UnescapedMessage = strings.ReplaceAll(utils.UnescapeCqCode(e.RawMessage), `\/`, `/`)
		ctx.ParsedSegments = message.ParseCqCodes(e.RawMessage) // 里面解析了再 unescape

		ctx.IsForwardMsg = ctx.ParsedSegments.WithType(message.TYPE_FORWARD)
		ctx.IsJsonMsg = ctx.ParsedSegments.WithType(message.TYPE_JSON)
		ctx.IsXmlMsg = ctx.ParsedSegments.WithType(message.TYPE_XML)
		ctx.IsCardMsg = ctx.IsXmlMsg || ctx.IsJsonMsg

		for _, replySeg := range ctx.ParsedSegments.GetType(message.TYPE_REPLY) {
			replyId, ok := replySeg.Data["id"].(int)
			if !ok {
				b.log.Warn().Msgf("reply segment without int \"id\": %v", replySeg)
			}
			if replyId != 0 {
				ctx.ReplyId = replyId
			}
		}

		if e.GroupId != 0 {
			for _, atSeg := range ctx.ParsedSegments.GetType(message.TYPE_AT) {
				atId, ok := atSeg.Data["qq"].(int)
				if !ok {
					b.log.Warn().Msgf("at segment without int \"qq\": %v", atSeg)
				}
				if atId != 0 {
					ctx.Ats = utils.NoReduplicateAppend(ctx.Ats, atId)
				}
			}
		}
	}

	ctx.IsToMe = e.TypeL2 == "private" ||
		e.TargetId == b.selfId ||
		slices.Contains(ctx.Ats, b.selfId) ||
		utils.StringsContainsMulti(e.RawMessage, b.alias)

	return ctx
}

func (b *Bot) matchEvent(ctx *Ctx) {
	if !b.msgProcessEnabled.Load() {
		return
	}
	b.log.Trace().Msgf("len(b.matchers.m): %d", len(b.matchers.m))
	for mName, matcher := range b.matchers.m {
		go func(mName string, matcher *Matcher) {
			if matcher.sequential {
				if matcher.skipOnLockFail {
					if !matcher.mu.TryLock() {
						b.log.Debug().Msgf("matcher %s lock failed, skip", mName)
						return
					}
				} else {
					matcher.mu.Lock()
				}
				b.log.Trace().Msgf("matcher %s lock success", mName)
				defer func() {
					matcher.mu.Unlock()
					b.log.Trace().Msgf("matcher %s unlock success", mName)
				}()
			}

			pass, i := matcher.rules.rule(ctx)

			traceLog := fmt.Sprintf("matching: %s (%d/%d)", mName, i, len(matcher.rules))
			b.log.Trace().Msg(traceLog)

			if pass {
				i := new(int)
				defer func() {
					if err := recover(); err != nil {
						panicErr := fmt.Sprintf("matcher %s[%d] panic: %v", mName, *i, err)
						stack := debug.Stack()
						ctx := ctx.String()

						b.log.Error().Msg(panicErr)
						os.Stderr.Write(stack)
						b.log.Error().Msgf("ctx: %s", ctx)
						b.log.Error().Msgf("match trace: %s", traceLog)

						b.Log2Sus.Error(panicErr)
						b.Log2Sus.Error(string(stack))
						b.Log2Sus.Error("ctx: ", ctx)
						b.Log2Sus.Error("match trace: ", traceLog)
					}
				}()
				*i = matcher.handlers.handle(ctx)
			}
		}(mName, matcher)
	}
}

func (b *Bot) formatGroupInfo(groupId int) string {
	gi, _ := b.GetGroupInfoTryCache(groupId)
	if gi != nil {
		return "(" + utils.RuneCut(gi.GroupName) + " " + itoa(groupId) + ")"
	} else {
		return itoa(groupId)
	}
}

func (b *Bot) formatStrangerInfo(userId int) string {
	si, _ := b.GetStrangerInfoTryCache(userId)
	if si != nil {
		return "(" + utils.RuneCut(si.Nickname) + " " + itoa(userId) + ")"
	} else {
		return itoa(userId)
	}
}

func (b *Bot) FormatEvent(e any) string {
	switch e := e.(type) {
	case *event.MessagePrivate:
		switch e.SubType {
		case event.TYPE_L3_MESSAGE_PRIVATE_FRIEND:
			return fmt.Sprintf(
				"收到好友 (%s %d) 的私聊消息(%d): %s",
				utils.RuneCut(e.Sender.Nickname), e.Sender.UserId, e.MessageId, e.RawMessage,
			)
		case event.TYPE_L3_MESSAGE_PRIVATE_GROUP:
			return fmt.Sprintf(
				"收到群成员 (%s %d) 的临时会话消息(%d): %s",
				utils.RuneCut(e.Sender.Nickname), e.Sender.UserId, e.MessageId, e.RawMessage,
			)
		case event.TYPE_L3_MESSAGE_PRIVATE_OTHER:
			return fmt.Sprintf(
				"收到 (%s %d) 的其他消息(%d): %s",
				utils.RuneCut(e.Sender.Nickname), e.Sender.UserId, e.MessageId, e.RawMessage,
			)
		default:
			return fmt.Sprintf("unknown MessagePrivate subtype(%s): %+v", e.SubType, e)
		}

	case *event.MessageGroup:
		// 获取群名
		group := b.formatGroupInfo(e.GroupId)

		switch e.SubType {
		case event.TYPE_L3_MESSAGE_GROUP_NORMAL:
			return fmt.Sprintf(
				"在 %s 收到 (%s %d) 的群聊消息(%d): %s",
				group, utils.RuneCut(e.Sender.Nickname), e.Sender.UserId, e.MessageId, e.RawMessage,
			)
		case event.TYPE_L3_MESSAGE_GROUP_ANONYMOUS:
			return fmt.Sprintf(
				"在 %s 收到 (%s %d) 的匿名消息(%d): %s",
				group, utils.RuneCut(e.Sender.Nickname), e.Sender.UserId, e.MessageId, e.RawMessage,
			)
		case event.TYPE_L3_MESSAGE_GROUP_NOTICE:
			return fmt.Sprintf(
				"在 %s 收到 (%s %d) 的系统提示(%d): %s",
				group, utils.RuneCut(e.Sender.Nickname), e.Sender.UserId, e.MessageId, e.RawMessage,
			)
		default:
			return fmt.Sprintf("unknown MessageGroup subtype(%s): %+v", e.SubType, e)
		}

	case *event.RequestFriend:
		// 获取陌生人昵称
		stranger := b.formatStrangerInfo(e.UserId)

		return fmt.Sprintf(
			"收到 %s 的好友请求: %s (%s)",
			stranger, e.Comment, e.Flag,
		)

	case *event.RequestGroup:
		// 获取群名、陌生人昵称
		var group, stranger string
		wg := sync.WaitGroup{}
		wg.Go(func() { group = b.formatGroupInfo(e.GroupId) })
		wg.Go(func() { stranger = b.formatStrangerInfo(e.UserId) })
		wg.Wait()

		switch e.SubType {
		case event.TYPE_L3_REQUEST_GROUP_ADD:
			return fmt.Sprintf(
				"在 %s 收到 %s 的入群请求: %s (%s)",
				group, stranger, e.Comment, e.Flag,
			)
		case event.TYPE_L3_REQUEST_GROUP_INVITE:
			return fmt.Sprintf(
				"收到 %s 邀请加群 %s: (%s)",
				stranger, group, e.Flag,
			)
		default:
			return fmt.Sprintf("unknown RequestGroup subtype: %s", e.SubType)
		}

	case *event.NoticeGroupUpload:
		// 获取群名、成员昵称
		var group, stranger string
		wg := sync.WaitGroup{}
		wg.Go(func() { group = b.formatGroupInfo(e.GroupId) })
		wg.Go(func() { stranger = b.formatStrangerInfo(e.UserId) })
		wg.Wait()

		return fmt.Sprintf(
			"群 %s 成员 %s 上传了文件: %s(%s)",
			group, stranger, e.File.Name, utils.FormatBytes(uint64(e.File.Size)),
		)

	case *event.NoticeGroupAdmin:
		// 获取群名、成员昵称
		var group, stranger string
		wg := sync.WaitGroup{}
		wg.Go(func() { group = b.formatGroupInfo(e.GroupId) })
		wg.Go(func() { stranger = b.formatStrangerInfo(e.UserId) })
		wg.Wait()

		switch e.SubType {
		case event.TYPE_L3_NOTICE_GROUP_ADMIN_SET:
			return fmt.Sprintf(
				"群 %s 设置了新管理员 %s",
				group, stranger,
			)
		case event.TYPE_L3_NOTICE_GROUP_ADMIN_UNSET:
			return fmt.Sprintf(
				"群 %s 取消了管理员 %s",
				group, stranger,
			)
		default:
			return fmt.Sprintf("unknown NoticeGroupAdmin subtype: %s", e.SubType)
		}

	case *event.NoticeGroupDecrease:
		// 获取群名、管理员昵称、成员昵称
		var group, operator, stranger string
		wg := sync.WaitGroup{}
		wg.Go(func() { group = b.formatGroupInfo(e.GroupId) })
		wg.Go(func() { operator = b.formatStrangerInfo(e.OperatorId) })
		wg.Go(func() { stranger = b.formatStrangerInfo(e.UserId) })
		wg.Wait()

		switch e.SubType {
		case event.TYPE_L3_NOTICE_GROUP_DECREASE_LEAVE:
			return fmt.Sprintf(
				"群 %s 成员 %s 离开了群",
				group, stranger,
			)
		case event.TYPE_L3_NOTICE_GROUP_DECREASE_KICK:
			return fmt.Sprintf(
				"群 %s 管理员 %s 踢出了成员 %s",
				group, operator, stranger,
			)
		case event.TYPE_L3_NOTICE_GROUP_DECREASE_KICK_ME:
			b.log.Warn().Msgf(
				"群 %s 管理员 %s 踢出了本账号",
				group, operator,
			)
		default:
			return fmt.Sprintf("unknown NoticeGroupDecrease subtype: %s", e.SubType)
		}

	case *event.NoticeGroupIncrease:
		// 获取群名、管理员昵称、陌生人昵称
		var group, operator, stranger string
		wg := sync.WaitGroup{}
		wg.Go(func() { group = b.formatGroupInfo(e.GroupId) })
		wg.Go(func() { operator = b.formatStrangerInfo(e.OperatorId) })
		wg.Go(func() { stranger = b.formatStrangerInfo(e.UserId) })
		wg.Wait()

		switch e.SubType {
		case event.TYPE_L3_NOTICE_GROUP_INCREASE_APPROVE:
			return fmt.Sprintf(
				"群 %s 管理员 %s 同意了 %s 的入群请求",
				group, operator, stranger,
			)
		case event.TYPE_L3_NOTICE_GROUP_INCREASE_INVITE:
			return fmt.Sprintf(
				"群 %s 管理员 %s 邀请了 %s 加入群",
				group, operator, stranger,
			)
		default:
			return fmt.Sprintf("unknown NoticeGroupIncrease subtype: %s", e.SubType)
		}

	case *event.NoticeGroupBan:
		// 获取群名、管理员昵称、成员昵称
		var group, operator, stranger string
		wg := sync.WaitGroup{}
		wg.Go(func() { group = b.formatGroupInfo(e.GroupId) })
		wg.Go(func() { operator = b.formatStrangerInfo(e.OperatorId) })
		wg.Go(func() { stranger = b.formatStrangerInfo(e.UserId) })
		wg.Wait()

		switch e.SubType {
		case event.TYPE_L3_NOTICE_GROUP_BAN_BAN:
			return fmt.Sprintf(
				"群 %s 管理员 %s 禁言了 %s (%s)",
				group, operator, stranger, utils.FormatDuration(time.Duration(e.Duration)*time.Second),
			)
		case event.TYPE_L3_NOTICE_GROUP_BAN_LIFT_BAN:
			return fmt.Sprintf(
				"群 %s 管理员 %s 解除了 %s 的禁言",
				group, operator, stranger,
			)
		default:
			return fmt.Sprintf("unknown NoticeGroupBan subtype: %s", e.SubType)
		}

	case *event.NoticeFriendAdd:
		// 获取陌生人昵称
		stranger := b.formatStrangerInfo(e.UserId)

		return fmt.Sprintf("添加了新好友 %s", stranger)

	case *event.NoticeGroupRecall:
		selfRecall := e.OperatorId == e.UserId // 用户自己撤回自己
		// 获取群名、成员昵称、消息内容、操作者昵称
		var group, stranger, operator string
		msg := "<没收到>"
		wg := sync.WaitGroup{}
		wg.Go(func() { group = b.formatGroupInfo(e.GroupId) })
		wg.Go(func() { stranger = b.formatStrangerInfo(e.UserId) })
		if !selfRecall {
			wg.Go(func() { operator = b.formatStrangerInfo(e.OperatorId) })
		}
		if e.UserId != b.selfId { // bot 本体撤回的消息拿不到
			wg.Go(func() {
				mg, _ := b.GetMessageGroupTryCache(e.GroupId, e.MessageId)
				if mg != nil {
					msg = mg.RawMessage
				}
			})
		}
		wg.Wait()

		if !selfRecall {
			return fmt.Sprintf(
				"群 %s 管理员 %s 撤回了成员 %s 的消息(%d): %s",
				group, operator, stranger, e.MessageId, msg,
			)
		} else {
			return fmt.Sprintf(
				"群 %s 成员 %s 撤回了消息(%d): %s",
				group, stranger, e.MessageId, msg,
			)
		}

	case *event.NoticeFriendRecall:
		if e.UserId != b.selfId { // 获取陌生人昵称
			stranger := b.formatStrangerInfo(e.UserId)

			return fmt.Sprintf(
				"好友 %s 撤回了消息: %d",
				stranger, e.MessageId,
			)
		} else {
			return fmt.Sprintf(
				"撤回了自己的消息: %d", e.MessageId,
			)
		}

	case *event.NoticeNotifyPoke:
		selfPoke := e.TargetId == e.UserId
		// 获取群名、成员昵称、被戳成员昵称
		var group, stranger, target string
		wg := sync.WaitGroup{}
		wg.Go(func() { stranger = b.formatStrangerInfo(e.UserId) })
		if e.GroupId != 0 {
			wg.Go(func() { group = b.formatGroupInfo(e.GroupId) })
		}
		if !selfPoke && e.TargetId != b.selfId {
			wg.Go(func() { target = b.formatStrangerInfo(e.TargetId) })
		}
		wg.Wait()

		if e.GroupId != 0 {
			if !selfPoke {
				if e.TargetId != b.selfId {
					return fmt.Sprintf(
						"群 %s 成员 %s 戳了戳 %s",
						group, stranger, target,
					)
				} else {
					return fmt.Sprintf(
						"群 %s 成员 %s 戳了戳 Bot",
						group, stranger,
					)
				}
			} else {
				return fmt.Sprintf(
					"群 %s 成员 %s 戳了戳自己",
					group, stranger,
				)
			}
		} else {
			if !selfPoke {
				return fmt.Sprintf(
					"好友 %s 戳了戳 Bot",
					stranger,
				)
			} else {
				return fmt.Sprintf(
					"好友 %s 戳了戳自己",
					stranger,
				)
			}
		}

	case *event.NoticeNotifyLuckyKing:
		// 获取群名、成员昵称、运气王昵称
		var group, stranger, target string
		wg := sync.WaitGroup{}
		wg.Go(func() { group = b.formatGroupInfo(e.GroupId) })
		wg.Go(func() { stranger = b.formatStrangerInfo(e.UserId) })
		if e.UserId != e.TargetId {
			wg.Go(func() { target = b.formatStrangerInfo(e.TargetId) })
		} else {
			target = stranger
		}
		wg.Wait()

		return fmt.Sprintf(
			"群 %s 成员 %s 发送的红包由 %s 夺得运气王",
			group, stranger, target,
		)

	case *event.NoticeNotifyHonor:
		// 获取群名、成员昵称
		var group, stranger string
		wg := sync.WaitGroup{}
		wg.Go(func() { group = b.formatGroupInfo(e.GroupId) })
		wg.Go(func() { stranger = b.formatStrangerInfo(e.UserId) })
		wg.Wait()

		return fmt.Sprintf(
			"群 %s 成员 %s 荣誉变更: %s",
			group, stranger, event.ParseHonorType(e.HonorType),
		)

	case *event.MetaEventLifecycle:
		return fmt.Sprintf("lifecycle: %v", *e)

	case *event.MetaEventHeartbeat:
		return fmt.Sprintf("heartbeat: %v", *e)

	case nil:
		return "<nil event>"
	}
	return fmt.Sprintf("<unknown event(%T): %v>", e, e)
}

// statistics 统计
type statistics struct {
	statistics_Std
	statistics_Nc
}

type statistics_Std struct {
	receive int // 下发的数据包数量

	apiResp int // api 响应数量
	event   int // 事件数量

	messagePrivate int // 私聊消息数量
	messageGroup   int // 群聊消息数量

	noticeGroupUpload     int // 群文件上传数量
	noticeGroupAdmin      int // 群管理员变动数量
	noticeGroupDecrease   int // 群成员减少数量
	noticeGroupIncrease   int // 群成员增加数量
	noticeGroupBan        int // 群成员禁言数量
	noticeFriendAdd       int // 好友添加数量
	noticeGroupRecall     int // 群消息撤回数量
	noticeFriendRecall    int // 好友消息撤回数量
	noticeNotifyPoke      int // 群内戳一戳数量
	noticeNotifyLuckyKing int // 群红包运气王数量
	noticeNotifyHonor     int // 群成员荣誉变更数量

	requestFriend int // 好友请求数量
	requestGroup  int // 加群请求／邀请数量

	metaEventLifecycle int // 生命周期数量
	metaEventHeartbeat int // 心跳包数量
}

type statistics_Nc struct {
	noticeBotOffline  int
	noticeBotOnline   int
	noticeEssence     int
	noticeGroupName   int
	noticeReaction    int
	noticeOfflineFile int
}

// normalizePoke 统一 poke 通知的发起者字段
//
// NapCat 的 poke 事件里 user_id 与 target_id 都是被戳方, 发起方是 sender_id;
// OneBot 11 标准的实现不带 sender_id, 那种实现里 user_id 才是发起方。
// 归一化后消费方 (OnNoticeNotifyPoke) 统一读 SenderId。
// 详见 docs/napcat-protocol-differences.md
func normalizePoke(p *event.NoticeNotifyPoke) {
	if p == nil {
		return
	}
	if p.SenderId == 0 {
		p.SenderId = p.UserId
	}
}
