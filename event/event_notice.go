package event

// https://github.com/botuniverse/onebot-11/blob/master/event/notice.md

// NoticeBase 通知事件
type NoticeBase struct {
	Base // [TYPE_L1_NOTICE] "notice"

	// 通知类型
	NoticeType string `json:"notice_type" mapstructure:"notice_type"`
}

// NoticeGroupUpload 群文件上传
type NoticeGroupUpload struct {
	NoticeBase // [TYPE_L2_NOTICE_GROUP_UPLOAD] "group_upload"

	// 群号
	GroupId int `json:"group_id" mapstructure:"group_id"`
	// 发送者 QQ 号
	UserId int `json:"user_id" mapstructure:"user_id"`
	// 文件信息
	File fileInfo `json:"file" mapstructure:"file"`
}

type fileInfo struct {
	// 文件 ID
	Id string `json:"id" mapstructure:"id"`
	// 文件名
	Name string `json:"name" mapstructure:"name"`
	// 文件大小（字节数）
	Size int `json:"size" mapstructure:"size"`
	// busid（目前不清楚有什么作用）
	Busid int `json:"busid" mapstructure:"busid"`

	fileInfo_Nc
}

type fileInfo_Nc struct {
	Url string `json:"url" mapstructure:"url"`
}

// NoticeGroupAdmin 群管理员变动
type NoticeGroupAdmin struct {
	NoticeBase // [TYPE_L2_NOTICE_GROUP_ADMIN] "group_admin"

	// "set", "unset" // 事件子类型，分别表示设置和取消管理员
	SubType string `json:"sub_type" mapstructure:"sub_type"`
	// 群号
	GroupId int `json:"group_id" mapstructure:"group_id"`
	// 管理员 QQ 号
	UserId int `json:"user_id" mapstructure:"user_id"`
}

// NoticeGroupDecrease 群成员减少
type NoticeGroupDecrease struct {
	NoticeBase // [TYPE_L2_NOTICE_GROUP_DECREASE] "group_decrease"

	// "leave", "kick", "kick_me" // 事件子类型，分别表示主动退群、成员被踢、登录号被踢
	SubType string `json:"sub_type" mapstructure:"sub_type"`
	// 群号
	GroupId int `json:"group_id" mapstructure:"group_id"`
	// 操作者 QQ 号（如果是主动退群，则和 `user_id` 相同）
	OperatorId int `json:"operator_id" mapstructure:"operator_id"`
	// 离开者 QQ 号
	UserId int `json:"user_id" mapstructure:"user_id"`
}

// NoticeGroupIncrease 群成员增加
type NoticeGroupIncrease struct {
	NoticeBase // [TYPE_L2_NOTICE_GROUP_INCREASE] "group_increase"

	// "approve", "invite", "invite_approve"(NapCat 扩展) // 事件子类型，分别表示管理员已同意入群、管理员邀请入群
	SubType string `json:"sub_type" mapstructure:"sub_type"`
	// 群号
	GroupId int `json:"group_id" mapstructure:"group_id"`
	// 操作者 QQ 号
	OperatorId int `json:"operator_id" mapstructure:"operator_id"`
	// 加入者 QQ 号
	UserId int `json:"user_id" mapstructure:"user_id"`
}

// NoticeGroupBan 群禁言
type NoticeGroupBan struct {
	NoticeBase // [TYPE_L2_NOTICE_GROUP_BAN] "group_ban"

	// "ban", "lift_ban" // 事件子类型，分别表示禁言、解除禁言
	SubType string `json:"sub_type" mapstructure:"sub_type"`
	// 群号
	GroupId int `json:"group_id" mapstructure:"group_id"`
	// 操作者 QQ 号
	OperatorId int `json:"operator_id" mapstructure:"operator_id"`
	// 被禁言 QQ 号
	UserId int `json:"user_id" mapstructure:"user_id"`
	// 禁言时长，单位秒
	Duration int `json:"duration" mapstructure:"duration"`
}

// NoticeFriendAdd 好友添加
type NoticeFriendAdd struct {
	NoticeBase // [TYPE_L2_NOTICE_FRIEND_ADD] "friend_add"

	// 新添加好友 QQ 号
	UserId int `json:"user_id" mapstructure:"user_id"`
}

// NoticeGroupRecall 群消息撤回
type NoticeGroupRecall struct {
	NoticeBase // [TYPE_L2_NOTICE_GROUP_RECALL] "group_recall"

	// 群号
	GroupId int `json:"group_id" mapstructure:"group_id"`
	// 消息发送者 QQ 号
	UserId int `json:"user_id" mapstructure:"user_id"`
	// 操作者 QQ 号
	OperatorId int `json:"operator_id" mapstructure:"operator_id"`
	// 被撤回的消息 ID
	MessageId int `json:"message_id" mapstructure:"message_id"`

	NoticeRecall_Nc
}

// NoticeFriendRecall 好友消息撤回
type NoticeFriendRecall struct {
	NoticeBase // [TYPE_L2_NOTICE_FRIEND_RECALL] "friend_recall"

	// 好友 QQ 号
	UserId int `json:"user_id" mapstructure:"user_id"`
	// 被撤回的消息 ID
	MessageId int `json:"message_id" mapstructure:"message_id"`

	NoticeRecall_Nc
}

type NoticeRecall_Nc struct {
	Tip string `json:"tip" mapstructure:"tip"`
}

type NoticeNotify struct {
	NoticeBase // [TYPE_L2_NOTICE_NOTIFY] "notify"

	// "poke", "lucky_king", "honor" // 提示类型
	SubType string `json:"sub_type" mapstructure:"sub_type"`
	// 群号
	GroupId int `json:"group_id" mapstructure:"group_id"`
	// 发送者 QQ 号
	// 红包发送者 QQ 号
	// 成员 QQ 号
	UserId int `json:"user_id" mapstructure:"user_id"`
	// 被戳者 QQ 号
	// 运气王 QQ 号
	TargetId int `json:"target_id" mapstructure:"target_id"`
}

// NoticeNotifyPoke 群内戳一戳
type NoticeNotifyPoke struct {
	NoticeNotify // [TYPE_L3_NOTICE_NOTIFY_POKE] "poke"

	NoticeNotifyPoke_Nc
}

type NoticeNotifyPoke_Nc struct {
	Action       string `json:"action" mapstructure:"action"`
	Suffix       string `json:"suffix" mapstructure:"suffix"`
	ActionImgUrl string `json:"action_img_url" mapstructure:"action_img_url"`

	// SenderId 是 poke 的发起者
	//
	// NapCat 的 poke 语义与 OneBot 11 标准不同: user_id 与 target_id 都是被戳方,
	// 发起方放在 sender_id (标准里没有这个字段)。非 NapCat 实现不带 sender_id,
	// Easyonebot 在分发前会把它补齐成 user_id, 消费方统一读 SenderId。
	// 详见 docs/napcat-protocol-differences.md
	SenderId int `json:"sender_id" mapstructure:"sender_id"`
}

// NoticeNotifyLuckyKing 群红包运气王
type NoticeNotifyLuckyKing struct {
	NoticeNotify // [TYPE_L3_NOTICE_NOTIFY_LUCKY_KING] "lucky_king"
}

// NoticeNotifyHonor 群成员荣誉变更
type NoticeNotifyHonor struct {
	NoticeNotify // [TYPE_L3_NOTICE_NOTIFY_HONOR] "honor"

	// "talkative", "erformer", "emotion" // 荣誉类型，分别表示龙王、群聊之火、快乐源泉
	HonorType string `json:"honor_type" mapstructure:"honor_type"`
}
