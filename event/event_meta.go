package event

// https://github.com/botuniverse/onebot-11/blob/master/event/meta.md

// MetaEventBase 元事件
//
// 消息、通知、请求三大类事件是与聊天软件直接相关的、机器人真实接收到的事件，除了这些，OneBot 自己还会产生一类事件，这里称之为「元事件」，例如生命周期事件、心跳事件等，这类事件与 OneBot 本身的运行状态有关，而与聊天软件无关。元事件的上报方式和普通事件完全一样。
type MetaEventBase struct {
	Base // [TYPE_L1_METAEVENT] "meta_event"

	MetaEventType string `json:"meta_event_type" mapstructure:"meta_event_type"`
}

// MetaEventLifecycle 生命周期
//
// 注意，目前生命周期元事件中，只有 HTTP POST 的情况下可以收到 `enable` 和 `disable`，只有正向 WebSocket 和反向 WebSocket 可以收到 `connect`。
type MetaEventLifecycle struct {
	MetaEventBase // [TYPE_L2_META_LIFECYCLE] "lifecycle"

	// "enable", "disable", "connect" // 事件子类型，分别表示 OneBot 启用、停用、WebSocket 连接成功
	SubType string `json:"sub_type" mapstructure:"sub_type"`
}

// MetaEventHeartbeat 心跳
type MetaEventHeartbeat struct {
	MetaEventBase // [TYPE_L2_META_HEARTBEAT] "heartbeat"

	// 到下次心跳的间隔，单位毫秒
	Interval int `json:"interval" mapstructure:"interval"`
	// 状态信息 // check [Status] / [Status_Lgr] for more details
	Status map[string]any `json:"status" mapstructure:"status"`
}

// Status 状态信息
type Status struct {
	// 当前 QQ 在线，null 表示无法查询到在线状态
	Online bool `json:"online" mapstructure:"online"`
	// 状态符合预期，意味着各模块正常运行、功能正常，且 QQ 在线
	Good bool `json:"good" mapstructure:"good"`
}
