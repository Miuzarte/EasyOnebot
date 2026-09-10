package EasyOnebot

import (
	"regexp"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/Miuzarte/EasyOnebot/event"
	"github.com/Miuzarte/EasyOnebot/internal/utils"
	"github.com/Miuzarte/EasyOnebot/message"
)

type (
	rule     func(*Ctx) (pass bool)
	rules    []rule
	handler  func(*Ctx) (continue_ bool) // 返回 false 时终止剩余操作
	handlers []handler
	Matcher  struct {
		name           string
		sequential     bool // ture 时在匹配过程中持锁
		skipOnLockFail bool // true 时跳过锁失败的匹配
		mu             sync.Mutex
		rules          rules
		handlers       handlers
	}
	matchers struct {
		m  map[string]*Matcher // name: Matcher
		mu sync.Mutex
	}
)

func (rs *rules) append(rule ...rule) {
	*rs = utils.OnDemandAppend(*rs, rule...)
}

func (hs *handlers) append(handler ...handler) {
	*hs = utils.OnDemandAppend(*hs, handler...)
}

// rule returns the index of the first rule that fails, or len(rs) if all pass.
func (rs *rules) rule(ctx *Ctx) (pass bool, i int) {
	var r rule
	for i, r = range *rs {
		if !r(ctx) {
			return false, i
		}
	}
	return true, len(*rs)
}

// handle 返回最后执行的 handler 的索引
func (hs *handlers) handle(ctx *Ctx) (i int) {
	// 由调用方捕获 panic 方便打印信息
	var handler handler
	for i, handler = range *hs {
		if !handler(ctx) {
			break
		}
	}
	return i
}

// AddMatcher 添加匹配器,
// 已存在时会直接覆盖已存在的匹配器
func (b *Bot) AddMatcher(name string, m *Matcher) *Bot {
	b.matchers.mu.Lock()
	defer b.matchers.mu.Unlock()
	if len(m.rules) == 0 { // all pass
		b.log.Warn().Msg("Matcher.rules is empty")
	}
	if len(m.handlers) == 0 {
		b.log.Warn().Msg("Matcher.handlers is empty")
	}
	m.name = name
	b.matchers.m[name] = m
	return b
}

func (b *Bot) DelMatcher(name string) *Bot {
	b.matchers.mu.Lock()
	defer b.matchers.mu.Unlock()
	delete(b.matchers.m, name)
	return b
}

// NewMatcherWithoutFilter 不遵循名单过滤
func NewMatcherWithoutFilter() *Matcher {
	return &Matcher{}
}

// NewMatcher 默认遵循名单过滤
func NewMatcher() *Matcher {
	return (&Matcher{}).FollowFilter()
}

// WithMutex 使 matcher 在匹配与执行过程中持锁 避免并发,
// skipOnLockFail 为 true 时跳过锁失败的匹配
func (m *Matcher) WithMutex(skipOnLockFail bool) *Matcher {
	m.sequential = true
	m.skipOnLockFail = skipOnLockFail
	return m
}

func (m *Matcher) OnFunc(ruleFunc func(*Ctx) bool) *Matcher {
	m.rules.append(rule(ruleFunc))
	return m
}

func (m *Matcher) FollowFilter() *Matcher {
	return m.OnFunc(func(ctx *Ctx) bool {
		return !ctx.FilteredStranger && !ctx.FilteredGroup
	})
}

// OnTypes 匹配 [event.Event.TypeL1], [event.Event.TypeL2], [event.Event.TypeL3]
//
// check [event.Level1Type], [event.Level2Type], [event.Level3Type] for list
func (m *Matcher) OnTypes(typs ...[]string) *Matcher {
	compoundRule := make(rules, len(typs))
	switch len(typs) {
	case 3:
		compoundRule[2] = func(ctx *Ctx) bool { return slices.Contains(typs[2], ctx.Event.TypeL3) }
		fallthrough
	case 2:
		compoundRule[1] = func(ctx *Ctx) bool { return slices.Contains(typs[1], ctx.Event.TypeL2) }
		fallthrough
	case 1:
		compoundRule[0] = func(ctx *Ctx) bool { return slices.Contains(typs[0], ctx.Event.TypeL1) }
	default:
		panic("Matcher.OnTypes: invalied count of types")
	}
	return m.OnFunc(func(c *Ctx) bool {
		pass, _ := compoundRule.rule(c)
		return pass
	})
}

// OnTypeL1 匹配 [event.Event.TypeL1]
//
// check [event.Level1Type] for list
func (m *Matcher) OnTypeL1(typ ...event.Level1Type) *Matcher {
	if len(typ) == 0 {
		panic("Matcher.OnTypeL1: empty types")
	}
	return m.OnFunc(func(ctx *Ctx) bool {
		return slices.Contains(typ, ctx.Event.TypeL1)
	})
}

// OnTypeL2 匹配 [event.Event.TypeL2]
//
// check [event.Level2Type] for list
func (m *Matcher) OnTypeL2(typ ...event.Level2Type) *Matcher {
	if len(typ) == 0 {
		panic("Matcher.OnTypeL2: empty types")
	}
	return m.OnFunc(func(ctx *Ctx) bool {
		return slices.Contains(typ, ctx.Event.TypeL2)
	})
}

// OnTypeL3 匹配 [event.Event.TypeL3]
//
// check [event.Level3Type] for list
func (m *Matcher) OnTypeL3(typ ...event.Level3Type) *Matcher {
	if len(typ) == 0 {
		panic("Matcher.OnTypeL3: empty types")
	}
	return m.OnFunc(func(ctx *Ctx) bool {
		return slices.Contains(typ, ctx.Event.TypeL3)
	})
}

// WithType 匹配消息存在指定消息段类型
func (m *Matcher) WithType(typ ...message.SegType) *Matcher {
	if len(typ) == 0 {
		panic("Matcher.WithType: empty types")
	}
	return m.OnFunc(func(ctx *Ctx) bool {
		return ctx.ParsedSegments.WithType(typ...)
	})
}

func (m *Matcher) WithoutType(typ ...message.SegType) *Matcher {
	if len(typ) == 0 {
		panic("Matcher.WithoutType: empty types")
	}
	return m.OnFunc(func(ctx *Ctx) bool {
		return ctx.ParsedSegments.WithoutType(typ...)
	})
}

func (m *Matcher) OnlyType(typ ...message.SegType) *Matcher {
	if len(typ) == 0 {
		panic("Matcher.OnlyType: empty types")
	}
	return m.OnFunc(func(ctx *Ctx) bool {
		return ctx.ParsedSegments.OnlyType(typ...)
	})
}

// IsGroup 匹配消息来自指定群
func (m *Matcher) IsGroup(gid ...int) *Matcher {
	if len(gid) == 0 {
		panic("Matcher.IsGroup: empty gid")
	}
	return m.OnFunc(func(ctx *Ctx) bool {
		return slices.Contains(gid, ctx.Event.GroupId)
	})
}

func (m *Matcher) IsNotGroup(gid ...int) *Matcher {
	if len(gid) == 0 {
		panic("Matcher.IsNotGroup: empty gid")
	}
	return m.OnFunc(func(ctx *Ctx) bool {
		return !slices.Contains(gid, ctx.Event.GroupId)
	})
}

// IsUser 匹配消息来自指定用户
func (m *Matcher) IsUser(uid ...int) *Matcher {
	if len(uid) == 0 {
		panic("Matcher.IsUser: empty uid")
	}
	return m.OnFunc(func(ctx *Ctx) bool {
		return slices.Contains(uid, ctx.Event.UserId)
	})
}

func (m *Matcher) IsNotUser(uid ...int) *Matcher {
	if len(uid) == 0 {
		panic("Matcher.IsNotUser: empty uid")
	}
	return m.OnFunc(func(ctx *Ctx) bool {
		return !slices.Contains(uid, ctx.Event.UserId)
	})
}

// IsSuperuser 匹配消息来自超级用户 [Bot.sus]
func (m *Matcher) IsSuperuser() *Matcher {
	return m.OnFunc(func(ctx *Ctx) bool { return ctx.IsSuperuser })
}

func (m *Matcher) IsNotForwardMsg() *Matcher {
	return m.OnFunc(func(ctx *Ctx) bool { return !ctx.IsForwardMsg })
}

func (m *Matcher) IsNotXmlMsg() *Matcher {
	return m.OnFunc(func(ctx *Ctx) bool { return !ctx.IsXmlMsg })
}

func (m *Matcher) IsNotJsonMsg() *Matcher {
	return m.OnFunc(func(ctx *Ctx) bool { return !ctx.IsJsonMsg })
}

func (m *Matcher) IsNotCardMsg() *Matcher {
	return m.OnFunc(func(ctx *Ctx) bool { return !ctx.IsCardMsg })
}

// IsToMe 匹配是私聊消息、群聊@、群聊提到昵称/别名(处理假@)
func (m *Matcher) IsToMe() *Matcher {
	return m.OnFunc(func(ctx *Ctx) bool { return ctx.IsToMe })
}

// OnReply 匹配回复消息段
func (m *Matcher) OnReply() *Matcher {
	return m.OnFunc(func(ctx *Ctx) bool { return ctx.ReplyId != 0 })
}

// OnStringsContains 不区分大小写
func (m *Matcher) OnStringsContains(substr ...string) *Matcher {
	if len(substr) == 0 {
		panic("Matcher.OnStringsContains: empty substrings")
	}
	return m.OnFunc(func(ctx *Ctx) bool {
		return utils.StringsContainsMulti(ctx.Event.RawMessage, substr)
	})
}

// OnRegexpMatchString 匹配正则,
// 匹配 RawMessage
func (m *Matcher) OnRegexpMatchString(reg *regexp.Regexp) *Matcher {
	return m.OnFunc(func(ctx *Ctx) bool {
		return reg.MatchString(ctx.Event.RawMessage)
	})
}

// OnRegexpFindAllStringSubmatch 匹配正则捕获组,
// 匹配反转义后的消息, 不希望匹配到卡片消息时可以配合 [Matcher.IsNotCardMsg]
func (m *Matcher) OnRegexpFindAllStringSubmatch(reg *regexp.Regexp) *Matcher {
	return m.OnFunc(func(ctx *Ctx) bool {
		submatches := reg.FindAllStringSubmatch(ctx.UnescapedMessage, -1)
		if len(submatches) == 0 {
			return false
		}

		ctx.submatchesMu.Lock()
		defer ctx.submatchesMu.Unlock()
		if ctx.Submatches == nil {
			ctx.Submatches = map[string][][]string{}
		}
		ctx.Submatches[m.name] = submatches

		return true
	})
}

// DoCondition 在 handler 返回 false 时终止后续操作
func (m *Matcher) DoCondition(handler handler) *Matcher {
	m.handlers.append(handler)
	return m
}

// Once 需要放在只执行一次的 handler 之前,
// 保证接下来的所有操作只执行一次, 但仍需要外部销毁 matcher
func (m *Matcher) Once() *Matcher {
	once := &atomic.Bool{}
	return m.DoCondition(func(c *Ctx) bool {
		if once.CompareAndSwap(false, true) {
			return true
		}
		return false // stop remains
	})
}

// Do 执行 handler
func (m *Matcher) Do(handler func(*Ctx)) *Matcher {
	return m.DoCondition(func(c *Ctx) bool {
		handler(c)
		return true // 始终继续
	})
}

// Send 发送消息到上下文, 返回 nil 时不发送
//
// [TODO] 返回切片时发送合并转发
func (m *Matcher) Send(msgFunc func(*Ctx) any) *Matcher {
	return m.Do(func(c *Ctx) {
		msg := msgFunc(c)
		if msg != nil {
			_, err := c.SendMsg(msg)
			if err != nil {
				c.SendMsg(err)
			}
		}
	})
}

// Reply 以回复形式发送消息到上下文, 返回 nil 时不发送
//
// [TODO] 返回切片时发送合并转发
func (m *Matcher) Reply(msgFunc func(*Ctx) any) *Matcher {
	return m.Do(func(c *Ctx) {
		msg := msgFunc(c)
		if msg != nil {
			c.SendMsgReply(msg)
		}
	})
}

// Send 发送合并转发消息到上下文, 返回 nil 时不发送
func (m *Matcher) SendForward(msgFunc func(*Ctx) message.SegmentArray) *Matcher {
	return m.Do(func(c *Ctx) {
		msg := msgFunc(c)
		if msg != nil {
			c.SendForwardMsgAuto(msg)
		}
	})
}
