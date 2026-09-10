package event

type Level1Type = string

const ( // level 1
	TYPE_L1_MESSAGE   Level1Type = "message"    // 消息 [Message]
	TYPE_L1_NOTICE    Level1Type = "notice"     // 通知 [Notice]
	TYPE_L1_REQUEST   Level1Type = "request"    // 请求 [Request]
	TYPE_L1_METAEVENT Level1Type = "meta_event" // 元事件 [MetaEvent]

	TYPE_L1_MESSAGE_SENT Level1Type = "message_sent" // 自身发送的消息 (go-cqhttp)
)

type Level2Type = string

const ( // level 2
	TYPE_L2_MESSAGE_PRIVATE Level2Type = "private" // 私聊消息 [MessagePrivate]
	TYPE_L2_MESSAGE_GROUP   Level2Type = "group"   // 群消息 [MessageGroup]

	TYPE_L2_NOTICE_GROUP_UPLOAD   Level2Type = "group_upload"   // 群文件上传 [NoticeGroupUpload]
	TYPE_L2_NOTICE_GROUP_ADMIN    Level2Type = "group_admin"    // 群管理员变动 [NoticeGroupAdmin]
	TYPE_L2_NOTICE_GROUP_DECREASE Level2Type = "group_decrease" // 群成员减少 [NoticeGroupDecrease]
	TYPE_L2_NOTICE_GROUP_INCREASE Level2Type = "group_increase" // 群成员增加 [NoticeGroupIncrease]
	TYPE_L2_NOTICE_GROUP_BAN      Level2Type = "group_ban"      // 群禁言 [NoticeGroupBan]
	TYPE_L2_NOTICE_FRIEND_ADD     Level2Type = "friend_add"     // 好友添加 [NoticeFriendAdd]
	TYPE_L2_NOTICE_GROUP_RECALL   Level2Type = "group_recall"   // 群消息撤回 [NoticeGroupRecall]
	TYPE_L2_NOTICE_FRIEND_RECALL  Level2Type = "friend_recall"  // 好友消息撤回 [NoticeFriendRecall]
	TYPE_L2_NOTICE_NOTIFY         Level2Type = "notify"         // 系统提示 [NoticeNotify]...

	TYPE_L2_REQUEST_FRIEND Level2Type = "friend" // 加好友请求 [RequestFriend]
	TYPE_L2_REQUEST_GROUP  Level2Type = "group"  // 加群请求／邀请 [RequestGroup]

	TYPE_L2_META_LIFECYCLE Level2Type = "lifecycle" // 生命周期 [MetaEventLifecycle]
	TYPE_L2_META_HEARTBEAT Level2Type = "heartbeat" // 心跳 [MetaEventHeartbeat]
)

const ( // level 2 (NapCat 扩展)
	TYPE_L2_NOTICE_BOT_OFFLINE  Level2Type = "bot_offline"  // [NoticeBotOffline]
	TYPE_L2_NOTICE_BOT_ONLINE   Level2Type = "bot_online"   // [NoticeBotOnline]
	TYPE_L2_NOTICE_ESSENCE      Level2Type = "essence"      // [NoticeEssence]
	TYPE_L2_NOTICE_GROUP_NAME   Level2Type = "group_name"   // [NoticeGroupName]
	TYPE_L2_NOTICE_REACTION     Level2Type = "reaction"     // [NoticeReaction]
	TYPE_L2_NOTICE_OFFLINE_FILE Level2Type = "offline_file" // [NoticeOfflineFile]
)

type Level3Type = string

const ( // level 3
	TYPE_L3_MESSAGE_PRIVATE_FRIEND Level3Type = "friend"
	TYPE_L3_MESSAGE_PRIVATE_GROUP  Level3Type = "group"
	TYPE_L3_MESSAGE_PRIVATE_OTHER  Level3Type = "other"

	TYPE_L3_MESSAGE_GROUP_NORMAL    Level3Type = "normal"
	TYPE_L3_MESSAGE_GROUP_ANONYMOUS Level3Type = "anonymous"
	TYPE_L3_MESSAGE_GROUP_NOTICE    Level3Type = "notice"

	TYPE_L3_NOTICE_GROUP_ADMIN_SET   Level3Type = "set"
	TYPE_L3_NOTICE_GROUP_ADMIN_UNSET Level3Type = "unset"

	TYPE_L3_NOTICE_GROUP_DECREASE_LEAVE   Level3Type = "leave"
	TYPE_L3_NOTICE_GROUP_DECREASE_KICK    Level3Type = "kick"
	TYPE_L3_NOTICE_GROUP_DECREASE_KICK_ME Level3Type = "kick_me"

	TYPE_L3_NOTICE_GROUP_INCREASE_APPROVE Level3Type = "approve"
	TYPE_L3_NOTICE_GROUP_INCREASE_INVITE  Level3Type = "invite"

	TYPE_L3_NOTICE_GROUP_BAN_BAN      Level3Type = "ban"
	TYPE_L3_NOTICE_GROUP_BAN_LIFT_BAN Level3Type = "lift_ban"

	TYPE_L3_NOTICE_NOTIFY_POKE       Level3Type = "poke"       // [NoticeNotifyPoke]
	TYPE_L3_NOTICE_NOTIFY_LUCKY_KING Level3Type = "lucky_king" // [NoticeNotifyLuckyKing]
	TYPE_L3_NOTICE_NOTIFY_HONOR      Level3Type = "honor"      // [NoticeNotifyHonor]

	TYPE_L3_REQUEST_GROUP_ADD    Level3Type = "add"
	TYPE_L3_REQUEST_GROUP_INVITE Level3Type = "invite"

	TYPE_L3_META_LIFECYCLE_ENABLE  Level3Type = "enable"
	TYPE_L3_META_LIFECYCLE_DISABLE Level3Type = "disable"
	TYPE_L3_META_LIFECYCLE_CONNECT Level3Type = "connect"
)

const ( // level 3 (NapCat 扩展)
	TYPE_L3_NOTICE_REACTION_ADD    Level3Type = "add"
	TYPE_L3_NOTICE_REACTION_REMOVE Level3Type = "remove"
)
