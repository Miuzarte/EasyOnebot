package EasyOnebot

import (
	"context"
	"fmt"

	"github.com/Miuzarte/EasyOnebot/api"
	"github.com/Miuzarte/EasyOnebot/api/napcat"
	"github.com/Miuzarte/EasyOnebot/event"
	"github.com/Miuzarte/EasyOnebot/message"
)

// 本文件承载 ctx 的发送/撤回便捷方法。
//
// 内部统一走 [Bot.NapCat] (spec 生成类型 + 严格解码), 对外返回 [SendMsgResult] ——
// 群与私聊的响应字段相同, 这里统一成一套, 省掉调用方的类型分支。

// SendMsgResult 是发送消息/合并转发的结果, 字段与 spec 的 <OperationID>Data 对齐。
type SendMsgResult struct {
	MessageID int64   // 消息 ID
	ForwardID *string // 合并转发的 forward_id
	ResID     *string // 合并转发的 res_id
}

func resultFromGroup(d napcat.SendGroupMsgData) SendMsgResult {
	return SendMsgResult{MessageID: d.MessageID, ForwardID: d.ForwardID, ResID: d.ResID}
}

func resultFromPrivate(d napcat.SendPrivateMsgData) SendMsgResult {
	return SendMsgResult{MessageID: d.MessageID, ForwardID: d.ForwardID, ResID: d.ResID}
}

func resultFromGroupForward(d napcat.SendGroupForwardMsgData) SendMsgResult {
	return SendMsgResult{MessageID: d.MessageID, ForwardID: d.ForwardID, ResID: d.ResID}
}

func resultFromPrivateForward(d napcat.SendPrivateForwardMsgData) SendMsgResult {
	return SendMsgResult{MessageID: d.MessageID, ForwardID: d.ForwardID, ResID: d.ResID}
}

// SendPrivateMsgSegs 发送私聊消息至 Event.UserId
func (c *Ctx) SendPrivateMsgSegs(msg ...message.Segment) (*SendMsgResult, error) {
	return c.SendPrivateMsg(msg)
}

// SendPrivateMsgf 发送私聊消息至 Event.UserId
func (c *Ctx) SendPrivateMsgf(format string, a ...any) (*SendMsgResult, error) {
	return c.SendPrivateMsg(fmt.Sprintf(format, a...))
}

// SendPrivateMsg 发送私聊消息至 Event.UserId
func (c *Ctx) SendPrivateMsg(msg any, autoEscape ...bool) (*SendMsgResult, error) {
	segs := message.UnmarshalMessage(msg)
	body := napcat.SendPrivateMsgJSONBody{
		UserID:     new(fmt.Sprint(c.Event.UserId)),
		AutoEscape: boxAutoEscapePrivate(autoEscapePtr(autoEscape)),
	}
	var err error
	body.Message, err = napcat.BoxMessage(segs)
	if err != nil {
		return nil, err
	}
	data, err := c.Bot.NapCat().SendPrivateMsg(body)
	if err != nil {
		return nil, err
	}
	return new(resultFromPrivate(data)), nil
}

// SendGroupMsgSegs 发送群聊消息至 Event.GroupId
func (c *Ctx) SendGroupMsgSegs(msg ...message.Segment) (*SendMsgResult, error) {
	return c.SendGroupMsg(msg)
}

// SendGroupMsgf 发送群聊消息至 Event.GroupId
func (c *Ctx) SendGroupMsgf(format string, a ...any) (*SendMsgResult, error) {
	return c.SendGroupMsg(fmt.Sprintf(format, a...))
}

// SendGroupMsg 发送群聊消息至 Event.GroupId
func (c *Ctx) SendGroupMsg(msg any, autoEscape ...bool) (*SendMsgResult, error) {
	segs := message.UnmarshalMessage(msg)
	body := napcat.SendGroupMsgJSONBody{
		GroupID:    new(fmt.Sprint(c.Event.GroupId)),
		AutoEscape: boxAutoEscapeGroup(autoEscapePtr(autoEscape)),
	}
	var err error
	body.Message, err = napcat.BoxMessage(segs)
	if err != nil {
		return nil, err
	}
	data, err := c.Bot.NapCat().SendGroupMsg(body)
	if err != nil {
		return nil, err
	}
	return new(resultFromGroup(data)), nil
}

// SendMsgSegs 发送消息至 Event.MessageType, Event.UserId, Event.GroupId
func (c *Ctx) SendMsgSegs(msg ...message.Segment) (*SendMsgResult, error) {
	return c.SendMsg(msg)
}

// SendMsgf 发送消息至 Event.MessageType, Event.UserId, Event.GroupId
func (c *Ctx) SendMsgf(format string, a ...any) (*SendMsgResult, error) {
	return c.SendMsg(fmt.Sprintf(format, a...))
}

// SendMsg 发送消息至 Event.MessageType, Event.UserId, Event.GroupId
//
// msg 可以是 string (解析 CQ 码), 也可以是消息段 / 消息段数组。
func (c *Ctx) SendMsg(msg any, autoEscape ...bool) (*SendMsgResult, error) {
	if c.Event.TypeL2 == event.TYPE_L2_MESSAGE_PRIVATE {
		return c.SendPrivateMsg(msg, autoEscape...)
	}
	return c.SendGroupMsg(msg, autoEscape...)
}

// SendMsgReply 以回复形式 [Ctx.SendMsg]
func (c *Ctx) SendMsgReply(msg any, autoEscape ...bool) (*SendMsgResult, error) {
	segChain := message.SegmentArray{
		message.Reply(c.Event.MessageId),
	}
	switch msg := msg.(type) {
	case string: // 跳过字符串的解析
		segChain.Append(message.Text(msg))
	default:
		segChain.Append(message.UnmarshalMessage(msg)...)
	}
	return c.SendMsg(segChain, autoEscape...)
}

// SendMsgReplySegs 以回复形式 [Ctx.SendMsg]
func (c *Ctx) SendMsgReplySegs(msg ...message.Segment) (*SendMsgResult, error) {
	segChain := message.SegmentArray{
		message.Reply(c.Event.MessageId),
	}
	segChain.Append(msg...)
	return c.SendMsg(segChain)
}

// SendMsgReplyf 以回复形式 [Ctx.SendMsg]
func (c *Ctx) SendMsgReplyf(format string, a ...any) (*SendMsgResult, error) {
	return c.SendMsgReply(fmt.Sprintf(format, a...))
}

// DeleteMsg 撤回指定消息
func (c *Ctx) DeleteMsg(messageID int64) error {
	_, err := c.Bot.NapCat().DeleteMsg(napcat.DeleteMsgJSONBody{MessageID: boxDeleteMsgID(messageID)})
	return err
}

// DeleteMsgSf 撤回指定消息 (与 [Ctx.DeleteMsg] 等价, 保留旧名字)
func (c *Ctx) DeleteMsgSf(messageID int64) (error, bool) {
	return c.DeleteMsg(messageID), false
}

// SendForwardMsgAuto 发送合并转发
//
// 传入的必须是 node 消息段数组 (由 message.Node* 构造)。
func (c *Ctx) SendForwardMsgAuto(nodes message.SegmentArray, _ ...any) (*SendMsgResult, error) {
	for i, seg := range nodes {
		if seg.Type != "node" {
			panic(fmt.Sprintf("Ctx.SendForwardMsgAuto: unexpected segment type: %s at index %d", seg.Type, i))
		}
	}
	boxed, err := napcat.BoxMessage(nodes)
	if err != nil {
		return nil, err
	}

	switch c.Event.TypeL2 {
	case event.TYPE_L2_MESSAGE_GROUP:
		data, err := c.Bot.NapCat().SendGroupForwardMsg(napcat.SendGroupForwardMsgJSONBody{
			GroupID: new(fmt.Sprint(c.Event.GroupId)),
			Message: boxed,
		})
		if err != nil {
			return nil, err
		}
		return new(resultFromGroupForward(data)), nil
	case event.TYPE_L2_MESSAGE_PRIVATE:
		data, err := c.Bot.NapCat().SendPrivateForwardMsg(napcat.SendPrivateForwardMsgJSONBody{
			UserID:  new(fmt.Sprint(c.Event.UserId)),
			Message: boxed,
		})
		if err != nil {
			return nil, err
		}
		return new(resultFromPrivateForward(data)), nil
	default:
		return nil, fmt.Errorf("unsupported event type: %s", c.Event.TypeL2)
	}
}

// GetReplyMsg 根据回复的消息 ID 获取消息内容
func (c *Ctx) GetReplyMsg() (*api.GetMsgResp, error) {
	resp, err, _ := c.Std.GetMsgSf(c.ReplyId)
	// resp.Message.TryAtoi()
	return resp, err
}

// ---- 内部辅助 ----

func autoEscapePtr(autoEscape []bool) *bool {
	if len(autoEscape) == 0 {
		return nil
	}
	return &autoEscape[0]
}

// boxAutoEscapeGroup / boxAutoEscapePrivate 把 bool 包成生成类型里的 auto_escape 联合类型
func boxAutoEscapeGroup(v *bool) *napcat.SendGroupMsgJSONBody_AutoEscape {
	if v == nil {
		return nil
	}
	out := napcat.BoxTo[bool, napcat.SendGroupMsgJSONBody_AutoEscape](*v)
	return &out
}

func boxAutoEscapePrivate(v *bool) *napcat.SendPrivateMsgJSONBody_AutoEscape {
	if v == nil {
		return nil
	}
	out := napcat.BoxTo[bool, napcat.SendPrivateMsgJSONBody_AutoEscape](*v)
	return &out
}

// boxDeleteMsgID 把 int64 消息 ID 包成 delete_msg 要求的联合类型
func boxDeleteMsgID(id int64) napcat.DeleteMsgJSONBody_MessageID {
	return napcat.BoxTo[int64, napcat.DeleteMsgJSONBody_MessageID](id)
}

// RegisterTempMatcher 注册临时匹配器, 默认加锁执行, 且取消后会自动删除匹配器
func (c *Ctx) RegisterTempMatcher(ctx context.Context, name string, handler func(*Ctx), afterJob func()) {
	bot := c.Bot

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
