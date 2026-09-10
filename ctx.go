package EasyOnebot

import (
	"encoding/json"
	"sync"

	"github.com/Miuzarte/EasyOnebot/event"
	"github.com/Miuzarte/EasyOnebot/message"
)

type Ctx struct {
	Event            *event.Event
	ParsedSegments   message.SegmentArray // 根据 RawMessage 解析出的消息段, 相比于 Message 字段保证了类型
	UnescapedMessage string               // 反转义后的 RawMessage
	ReplyId          int                  // 回复的消息 ID
	Ats              []int                // 消息中 at 的 QQ 号

	FilteredStranger bool
	FilteredGroup    bool
	IsSuperuser      bool
	IsForwardMsg     bool
	IsXmlMsg         bool
	IsJsonMsg        bool
	IsCardMsg        bool
	IsToMe           bool

	Submatches   regSubmatches // ONLY exists when using `OnRegexpFindAllStringSubmatch`
	submatchesMu sync.Mutex
	MixCaller
	Bot *Bot `json:"-"` // ctx 便捷方法内部走 Bot.NapCat() 类型化调用; 含 func 字段, 不能进 JSON
}

func (c *Ctx) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

type regSubmatches map[string][][]string // name: submatches

func (rs regSubmatches) Get(moduleName string) [][]string {
	if rs == nil {
		return nil
	}
	// 处理完后只读, 不需要锁
	return rs[moduleName]
}
