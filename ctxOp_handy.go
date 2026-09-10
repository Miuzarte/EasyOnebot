package EasyOnebot

import (
	"context"
	"fmt"

	"github.com/Miuzarte/EasyOnebot/api"
	"github.com/Miuzarte/EasyOnebot/event"
	"github.com/Miuzarte/EasyOnebot/message"
)

// SendPrivateMsgSegs 发送私聊消息至 Event.UserId
func (c *Ctx) SendPrivateMsgSegs(msg ...message.Segment) (*api.SendPrivateMsgResp, error) {
	return c.Std.SendPrivateMsg(c.Event.UserId, msg)
}

// SendPrivateMsgf 发送私聊消息至 Event.UserId
func (c *Ctx) SendPrivateMsgf(format string, a ...any) (*api.SendPrivateMsgResp, error) {
	return c.Std.SendPrivateMsg(c.Event.UserId, fmt.Sprintf(format, a...))
}

// SendGroupMsgSegs 发送群聊消息至 Event.GroupId
func (c *Ctx) SendGroupMsgSegs(msg ...message.Segment) (*api.SendGroupMsgResp, error) {
	return c.Std.SendGroupMsg(c.Event.GroupId, msg)
}

// SendGroupMsgf 发送群聊消息至 Event.GroupId
func (c *Ctx) SendGroupMsgf(format string, a ...any) (*api.SendGroupMsgResp, error) {
	return c.Std.SendGroupMsg(c.Event.GroupId, fmt.Sprintf(format, a...))
}

// SendMsgSegs 发送消息至 Event.MessageType, Event.UserId, Event.GroupId
func (c *Ctx) SendMsgSegs(msg ...message.Segment) (*api.SendAnyMsgResp, error) {
	return c.Std.SendMsg(c.Event.MessageType, c.Event.UserId, c.Event.GroupId, msg)
}

// SendMsgf 发送消息至 Event.MessageType, Event.UserId, Event.GroupId
func (c *Ctx) SendMsgf(format string, a ...any) (*api.SendAnyMsgResp, error) {
	return c.Std.SendMsg(c.Event.MessageType, c.Event.UserId, c.Event.GroupId, fmt.Sprintf(format, a...))
}

// SendMsgReply 以回复形式 [Ctx.SendMsg]
func (c *Ctx) SendMsgReply(msg any, autoEscape ...bool) (*api.SendAnyMsgResp, error) {
	segChain := message.SegmentArray{
		message.Reply(c.Event.MessageId),
	}
	switch msg := msg.(type) {
	case string: // 跳过字符串的解析
		segChain.Append(message.Text(msg))
	default:
		segChain.Append(message.UnmarshalMessage(msg)...)
	}
	return c.Std.SendMsg(c.Event.MessageType, c.Event.UserId, c.Event.GroupId, segChain, autoEscape...)
}

// SendMsgReplySegs 以回复形式 [Ctx.SendMsg]
func (c *Ctx) SendMsgReplySegs(msg ...message.Segment) (*api.SendAnyMsgResp, error) {
	segChain := message.SegmentArray{
		message.Reply(c.Event.MessageId),
	}
	segChain.Append(msg...)
	return c.Std.SendMsg(c.Event.MessageType, c.Event.UserId, c.Event.GroupId, segChain)
}

// SendMsgReplyf 以回复形式 [Ctx.SendMsg]
func (c *Ctx) SendMsgReplyf(format string, a ...any) (*api.SendAnyMsgResp, error) {
	return c.SendMsgReply(fmt.Sprintf(format, a...))
}

// GetReplyMsg 根据回复的消息 ID 获取消息内容
func (c *Ctx) GetReplyMsg() (*api.GetMsgResp, error) {
	resp, err, _ := c.Std.GetMsgSf(c.ReplyId)
	// resp.Message.TryAtoi()
	return resp, err
}

// SendForwardMsgAuto 发送合并转发
func (c *Ctx) SendForwardMsgAuto(message message.SegmentArray) (*api.SendAnyForwardMsgResp, error) {
	for i, seg := range message {
		if seg.Type != "node" {
			panic(fmt.Sprintf("Ctx.SendForwardMsgAuto: unexpected segment type: %s at index %d", seg.Type, i))
		}
	}
	switch c.Event.TypeL2 {
	case "group":
		return c.Lgr.SendGroupForwardMsg(c.Event.GroupId, message)
		// return c.Nc.SendForwardMsg(c.Event.GroupId, 0, message, nil, "外显", "底下文本", "内容")
	case "private":
		return c.Lgr.SendPrivateForwardMsg(c.Event.UserId, message)
		// return c.Nc.SendForwardMsg(0, c.Event.UserId, message, nil, "外显", "底下文本", "内容")
	default:
		return nil, fmt.Errorf("unsupported event type: %s", c.Event.TypeL2)
	}
}

// RegisterTempMatcher 注册临时匹配器, 默认加锁执行, 且取消后会自动删除匹配器
func (c *Ctx) RegisterTempMatcher(ctx context.Context, name string, handler func(*Ctx), afterJob func()) {
	bot := c.Std.Callable.(*Bot)

	matcher := NewMatcher().WithMutex(false).
		OnTypes(
			[]string{c.Event.TypeL1},
			[]string{c.Event.TypeL2},
			[]string{c.Event.TypeL3},
		).Do(handler)
	switch c.Event.TypeL2 {
	case event.TYPE_L2_MESSAGE_GROUP:
		matcher.IsGroup(c.Event.GroupId)
		// fallthrough // 不限制用户, 其他群成员也可触发
	case event.TYPE_L2_MESSAGE_PRIVATE:
		matcher.IsUser(c.Event.UserId)
	}

	bot.AddMatcher(name, matcher)
	go func() {
		<-ctx.Done()
		bot.DelMatcher(name)
		if afterJob != nil {
			afterJob()
		}
	}()
}
