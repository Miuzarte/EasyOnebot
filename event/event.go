package event

import (
	"encoding/json"
)

type Base struct {
	// 事件发生的时间戳
	Time int `json:"time" mapstructure:"time"`
	// 收到事件的机器人 QQ 号
	SelfId int `json:"self_id" mapstructure:"self_id"`
	// 上报类型
	PostType string `json:"post_type" mapstructure:"post_type"`
}

// Event has everything,
// can convert to any type of event directly using [copier.Copy],
// by checking the TypeL1, TypeL2, TypeL3 fields.
type Event struct {
	// 同 PostType
	TypeL1 Level1Type `json:"-" mapstructure:"-"`
	// 详细类型
	TypeL2 Level2Type `json:"-" mapstructure:"-"`
	// 同 SubType
	TypeL3 Level3Type `json:"-" mapstructure:"-"`

	// [Base]
	Time     int    `json:"time" mapstructure:"time"`
	SelfId   int    `json:"self_id" mapstructure:"self_id"`
	PostType string `json:"post_type" mapstructure:"post_type"` // [L1]

	// [MessageBase]
	MessageType string `json:"message_type" mapstructure:"message_type"` // [L2] 消息类型
	SubType     string `json:"sub_type" mapstructure:"sub_type"`         // [L3]
	MessageId   int    `json:"message_id" mapstructure:"message_id"`
	UserId      int    `json:"user_id" mapstructure:"user_id"`
	Message     any    `json:"message" mapstructure:"message"`
	RawMessage  string `json:"raw_message" mapstructure:"raw_message"`
	Font        int    `json:"font" mapstructure:"font"`

	// [MessagePrivate]
	// [MessageGroup]
	GroupId   int         `json:"group_id" mapstructure:"group_id"`
	Sender    GroupSender `json:"sender" mapstructure:"sender"`
	Anonymous *Anonymous  `json:"anonymous" mapstructure:"anonymous"`

	// [NoticeBase]
	NoticeType string `json:"notice_type" mapstructure:"notice_type"` // [L2] 通知类型

	// [NoticeGroupUpload]
	File fileInfo `json:"file" mapstructure:"file"`

	// [NoticeGroupAdmin]

	// [NoticeGroupDecrease]
	OperatorId int `json:"operator_id" mapstructure:"operator_id"`

	// [NoticeGroupIncrease]

	// [NoticeGroupBan]
	Duration int `json:"duration" mapstructure:"duration"`

	// [NoticeFriendAdd]

	// [NoticeGroupRecall]
	// [NoticeFriendRecall]
	NoticeRecall_Lgr

	// [NoticeNotify]
	TargetId int `json:"target_id" mapstructure:"target_id"`

	// [NoticeNotifyPoke]
	NoticeNotifyPoke_Lgr

	// [NoticeNotifyLuckyKing]

	// [NoticeNotifyHonor]
	HonorType string `json:"honor_type" mapstructure:"honor_type"`

	// [RequestBase]
	RequestType string `json:"request_type" mapstructure:"request_type"` // [L2] 请求类型
	Comment     string `json:"comment" mapstructure:"comment"`
	Flag        string `json:"flag" mapstructure:"flag"`

	// [RequestFriend]

	// [RequestGroup]
	RequestGroup_Lgr

	// [MetaEventBase]
	MetaEventType string `json:"meta_event_type" mapstructure:"meta_event_type"` // [L2] 元事件类型

	// [MetaEventLifecycle]

	// [MetaEventHeartbeat]
	Interval int            `json:"interval" mapstructure:"interval"`
	Status   map[string]any `json:"status" mapstructure:"status"`

	// 为某些非标准实现预留的字段
	RawEvent map[string]any `json:"-" mapstructure:"-"`
}

func (e *Event) InitFields() {
	e.TypeL1 = e.PostType
	e.TypeL3 = e.SubType
	switch e.PostType {
	case TYPE_L1_MESSAGE, TYPE_L1_MESSAGE_SENT:
		e.TypeL2 = e.MessageType
	case TYPE_L1_NOTICE:
		e.TypeL2 = e.NoticeType
	case TYPE_L1_REQUEST:
		e.TypeL2 = e.RequestType
	case TYPE_L1_METAEVENT:
		e.TypeL2 = e.MetaEventType
	}
}

func (e *Event) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func ParseHonorType(honorType string) string {
	if v, ok := HonorTypes[honorType]; ok {
		return v
	}
	return honorType
}

// HonorTypes 荣誉类型中文
var HonorTypes = map[string]string{
	"talkative": "龙王",
	"performer": "群聊之火",
	"emotion":   "快乐源泉",
}
