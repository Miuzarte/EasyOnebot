package EasyOnebot // TODO: change to "EasyOneBot"

import (
	"context"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Miuzarte/EasyOnebot/api"
	"github.com/Miuzarte/EasyOnebot/internal/utils"

	"github.com/redis/rueidis"
	"github.com/rs/zerolog"
)

type Bot struct {
	ctx    context.Context
	cancel context.CancelFunc

	alias       []string         // 机器人别名 (默认包含 QQ 号, 昵称)
	sus         []int            // 超级用户
	selfId      int              // 机器人QQ号
	loginInfo   *api.LoginInfo   // 登录信息
	versionInfo *api.VersionInfo // onebot 实现信息

	settings settings // 设置
	filter   filter   // 过滤

	conn    connection // 连接
	hb      heartbeat  // 心跳
	apiPool apiRespChanPool
	redis   redisClient

	matchers     matchers     // 匹配器
	eventRecalls eventRecalls // 回调

	msgProcessEnabled atomic.Bool // 消息处理开关, 默认 true

	Log2Sus    log2Sus // 报告日志至超级用户
	log        *zerolog.Logger
	statistics statistics // 统计
}

// filter 过滤
type filter struct {
	strangers []int
	groups    []int
}

// settings 设置
type settings struct {
	msgLogOut       bool          // 日志输出收到的信息
	onlineNotify    bool          // 上线通知
	offlineNotify   bool          // 下线通知
	redisExpireMsg  time.Duration // 聊天记录过期时间
	redisExpireInfo time.Duration // 群/陌生人信息过期时间
}

func New() *Bot {
	// 默认日志, 外部可用 SetLogger 覆盖
	defaultLog := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Timestamp().
		Str("scope", "EasyOneBot").
		Logger()

	b := &Bot{
		ctx:    nil,
		cancel: nil,

		selfId:    0,
		alias:     []string{},
		sus:       []int{},
		loginInfo: nil,

		settings: settings{
			msgLogOut:       true,
			onlineNotify:    true,
			offlineNotify:   true,
			redisExpireMsg:  time.Hour * 24 * 7,
			redisExpireInfo: time.Hour,
		},
		filter: filter{},

		conn: connection{
			ws:           nil,
			url:          "",
			mu:           sync.Mutex{},
			dialCount:    0,
			lastDialTime: time.Time{},
		},
		hb: heartbeat{
			signal:   make(chan int, 1),
			running:  false,
			interval: 0,
			count:    0,
			lost:     0,
		},
		apiPool: apiRespChanPool{
			// timeout: time.Minute * 2, // 可撤回的截止时间
			timeout: time.Minute * 5,
			m:       map[string]chan *api.Response{},
			mu:      sync.Mutex{},
		},
		redis: redisClient{},

		matchers:     matchers{m: map[string]*Matcher{}},
		eventRecalls: eventRecalls{},

		log:        &defaultLog,
		Log2Sus:    log2Sus{},
		statistics: statistics{},
	}
	b.Log2Sus = log2Sus{b.wrapLogFunc()}
	b.msgProcessEnabled.Store(true)

	return b
}

// GetLoginInfoCache 获取登录信息, 如果已经缓存则直接返回
func (b *Bot) GetLoginInfoCache() (*api.LoginInfo, error) {
	if b.loginInfo != nil {
		return b.loginInfo, nil
	}
	return b.Call().Std.GetLoginInfo()
}

// GetVersionInfoCache 获取 OneBot 实现信息, 如果已经缓存则直接返回
func (b *Bot) GetVersionInfoCache() (*api.VersionInfo, error) {
	if b.versionInfo != nil {
		return b.versionInfo, nil
	}
	return b.Call().Std.GetVersionInfo()
}

// SetWsUrl 设置 websocket url, 无默认值
func (b *Bot) SetWsUrl(url string) *Bot {
	if !strings.HasPrefix(url, "ws://") {
		url = "ws://" + url
	}
	b.conn.url = strings.TrimSuffix(url, "/")
	return b
}

// SetToken 设置 OneBot token, 无默认值
func (b *Bot) SetToken(token string) *Bot {
	b.conn.token = token
	return b
}

// SetLogger 覆盖内部 logger, 传 nil 时忽略
//
// 不覆盖时使用内部默认 logger (ConsoleWriter, stderr)
func (b *Bot) SetLogger(l *zerolog.Logger) *Bot {
	if l != nil {
		b.log = l
	}
	return b
}

// SetLogLevel 设置此 Bot 的日志等级, 默认跟随 zerolog 全局等级
func (b *Bot) SetLogLevel(level zerolog.Level) *Bot {
	l := b.log.Level(level)
	b.log = &l
	return b
}

// setScope 替换日志的 scope, 用于登录后带上账号昵称
//
// 外部 logger 若已带 scope 字段, JSON 中会出现两个同名 key,
// 控制台渲染取最后一个, 即此处设置的值
func (b *Bot) setScope(scope string) {
	l := b.log.With().Str("scope", scope).Logger()
	b.log = &l
}

// SetLogOutMsg 设置是否打印收到的信息, 默认 true
func (b *Bot) SetLogOutMsg(enable bool) *Bot {
	b.settings.msgLogOut = enable
	return b
}

// SetOnlineNotify 设置上线通知, 默认 true
func (b *Bot) SetOnlineNotify(enable bool) *Bot {
	b.settings.onlineNotify = enable
	return b
}

// SetOfflineNotify 设置下线通知, 默认 true
func (b *Bot) SetOfflineNotify(enable bool) *Bot {
	b.settings.offlineNotify = enable
	return b
}

// SetRedisExpire 设置 Redis 过期时间, 默认无
func (b *Bot) SetRedisExpire(msg, info time.Duration) *Bot {
	b.settings.redisExpireMsg = msg
	b.settings.redisExpireInfo = info
	return b
}

// SetRedisClient 设置 [rueidis.Client] 客户端, nil 为不使用
func (b *Bot) SetRedisClient(client rueidis.Client) *Bot {
	b.redis.Client = client
	return b
}

// SetNicknames 设置机器人昵称, 不区分大小写
func (b *Bot) SetNicknames(nicknames []string) *Bot {
	b.alias = utils.StringsToLowerMulti(nicknames)
	return b
}

// AddNicknames 添加机器人昵称, 不区分大小写
func (b *Bot) AddNicknames(nicknames ...string) *Bot {
	b.alias = utils.NoReduplicateAppend(b.alias, utils.StringsToLowerMulti(nicknames)...)
	return b
}

// DelNickname 删除机器人昵称
func (b *Bot) DelNickname(nicknames ...string) *Bot {
	b.alias = utils.DelElements(b.alias, utils.StringsToLowerMulti(nicknames)...)
	return b
}

// SetSuperusers 设置超级用户
func (b *Bot) SetSuperusers(userIds []int) *Bot {
	b.sus = utils.DelElements(userIds, 0)
	return b
}

// AddSuperusers 添加超级用户
func (b *Bot) AddSuperusers(userIds ...int) *Bot {
	b.sus = utils.NoReduplicateAppend(b.sus, utils.DelElements(userIds, 0)...)
	return b
}

// DelSuperusers 删除超级用户
func (b *Bot) DelSuperusers(userIds ...int) *Bot {
	b.sus = utils.DelElements(b.sus, userIds...)
	return b
}

// SetFilterStranger 设置陌生人过滤
func (b *Bot) SetFilterStranger(userIds []int) *Bot {
	b.filter.strangers = userIds
	return b
}

// AddFilterStranger 添加陌生人过滤
func (b *Bot) AddFilterStranger(userIds ...int) *Bot {
	utils.NoReduplicateAppend(b.filter.strangers, utils.DelElements(userIds, 0)...)
	return b
}

// DelFilterStranger 删除陌生人过滤
func (b *Bot) DelFilterStranger(userIds ...int) *Bot {
	b.filter.strangers = utils.DelElements(b.filter.strangers, userIds...)
	return b
}

// SetFilterGroup 设置群过滤
func (b *Bot) SetFilterGroup(groupIds []int) *Bot {
	b.filter.groups = groupIds
	return b
}

// SetMsgProcessEnabled 设置是否处理消息(触发 matcher)
// 默认为 true, 设为 false 时 WS 连接保持, 但所有 matcher 不会执行,
// 适用于初始化阶段避免处理消息导致竞态
func (b *Bot) SetMsgProcessEnabled(enabled bool) *Bot {
	b.msgProcessEnabled.Store(enabled)
	return b
}

// AddFilterGroup 添加群过滤
func (b *Bot) AddFilterGroup(groupIds ...int) *Bot {
	utils.NoReduplicateAppend(b.filter.groups, utils.DelElements(groupIds, 0)...)
	return b
}

// DelFilterGroup 删除群过滤
func (b *Bot) DelFilterGroup(groupIds ...int) *Bot {
	b.filter.groups = utils.DelElements(b.filter.groups, groupIds...)
	return b
}
