package event

// https://github.com/botuniverse/onebot-11/blob/master/event/message.md

// MessageBase 消息事件
type MessageBase struct {
	Base // [TYPE_L1_MESSAGE] "message"

	// 消息类型
	MessageType string `json:"message_type" mapstructure:"message_type"`
	// "friend", "group", "other"      // 消息子类型，如果是好友则是 `friend`，如果是群临时会话则是 `group`
	// "normal", "anonymous", "notice" // 消息子类型，正常消息是 `normal`，匿名消息是 `anonymous`，系统提示（如「管理员已禁止群内匿名聊天」）是 `notice`
	SubType string `json:"sub_type" mapstructure:"sub_type"`
	// 消息 ID
	MessageId int `json:"message_id" mapstructure:"message_id"`
	// 发送者 QQ 号
	UserId int `json:"user_id" mapstructure:"user_id"`
	// [string], [message.SegmentArray] // 消息内容
	Message any `json:"message" mapstructure:"message"`
	// 原始消息内容
	RawMessage string `json:"raw_message" mapstructure:"raw_message"`
	// 字体
	Font int `json:"font" mapstructure:"font"`

	Message_Lgr
}

// MessagePrivate 私聊消息
type MessagePrivate struct {
	MessageBase // [TYPE_L2_MESSAGE_PRIVATE] "private"

	// 发送人信息
	Sender Sender `json:"sender" mapstructure:"sender"`
}

// MessageGroup 群消息
type MessageGroup struct {
	MessageBase // [TYPE_L2_MESSAGE_GROUP] "group"

	// 群号
	GroupId int `json:"group_id" mapstructure:"group_id"`
	// 发送人信息
	Sender GroupSender `json:"sender" mapstructure:"sender"`
	// 匿名信息，如果不是匿名消息则为 null
	Anonymous *Anonymous `json:"anonymous" mapstructure:"anonymous"`
}

// 需要注意的是，`sender` 中的各字段是尽最大努力提供的，也就是说，不保证每个字段都一定存在，也不保证存在的字段都是完全正确的（缓存可能过期）。
type Sender struct {
	// 发送者 QQ 号
	UserId int `json:"user_id" mapstructure:"user_id"`
	// 昵称
	Nickname string `json:"nickname" mapstructure:"nickname"`
	// 性别，`male` 或 `female` 或 `unknown`
	Sex string `json:"sex" mapstructure:"sex"`
	// 年龄
	Age int `json:"age" mapstructure:"age"`
}

// 需要注意的是，`sender` 中的各字段是尽最大努力提供的，也就是说，不保证每个字段都一定存在，也不保证存在的字段都是完全正确的（缓存可能过期）。尤其对于匿名消息，此字段不具有参考价值。
type GroupSender struct {
	Sender `mapstructure:"sender,squash"` // squash 让 mapstructure 不要将其展开为子结构

	// 群名片／备注
	Card string `json:"card" mapstructure:"card"`
	// 地区
	Area string `json:"area" mapstructure:"area"`
	// 成员等级
	Level string `json:"level" mapstructure:"level"`
	// 角色，`owner` 或 `admin` 或 `member`
	Role string `json:"role" mapstructure:"role"`
	// 专属头衔
	Title string `json:"title" mapstructure:"title"`
}

func (gs *GroupSender) GetCardOrNickname() string {
	if gs.Card != "" {
		return gs.Card
	}
	return gs.Nickname
}

type Anonymous struct {
	// 匿名用户 ID
	Id int `json:"id" mapstructure:"id"`
	// 匿名用户名称
	Name string `json:"name" mapstructure:"name"`
	// 匿名用户 flag，在调用禁言 API 时需要传入
	Flag string `json:"flag" mapstructure:"flag"`
}

type (
	Message_Lgr struct {
		MessageStyle *MessageStyle `json:"message_style" mapstructure:"message_style"`
	}
	MessageStyle struct {
		BubbleId              int  `json:"bubble_id" mapstructure:"bubble_id"`
		PendantId             int  `json:"pendant_id" mapstructure:"pendant_id"`
		FontId                int  `json:"font_id" mapstructure:"font_id"`
		FontEffectId          int  `json:"font_effect_id" mapstructure:"font_effect_id"`
		IsCsFontEffectEnabled bool `json:"is_cs_font_effect_enabled" mapstructure:"is_cs_font_effect_enabled"`
		BubbleDiyTextId       int  `json:"bubble_diy_text_id" mapstructure:"bubble_diy_text_id"`
	}
)
