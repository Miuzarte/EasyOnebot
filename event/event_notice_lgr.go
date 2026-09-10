package event

type NoticeBotOffline struct {
	NoticeBase // [TYPE_L2_NOTICE_BOT_OFFLINE] "bot_offline"

	Tag     string `json:"tag" mapstructure:"tag"`
	Message string `json:"message" mapstructure:"message"`
}

type NoticeBotOnline struct {
	NoticeBase // [TYPE_L2_NOTICE_BOT_ONLINE] "bot_online"

	// 0: Login, 1: Reconnect
	Reason int `json:"reason" mapstructure:"reason"`
}

type NoticeEssence struct { // GroupEssence
	NoticeBase // [TYPE_L2_NOTICE_ESSENCE] "essence"

	// "add", "delete"
	SubType    string `json:"sub_type" mapstructure:"sub_type"`
	GroupId    int    `json:"group_id" mapstructure:"group_id"`
	SenderId   int    `json:"sender_id" mapstructure:"sender_id"`
	OperatorId int    `json:"operator_id" mapstructure:"operator_id"`
	MessageId  int    `json:"message_id" mapstructure:"message_id"`
}

type NoticeGroupName struct { // GroupNameChange
	NoticeBase // [TYPE_L2_NOTICE_GROUP_NAME] "group_name"

	GroupId int    `json:"group_id" mapstructure:"group_id"`
	Name    string `json:"name" mapstructure:"name"`
}

type NoticeReaction struct {
	NoticeBase // [TYPE_L2_NOTICE_REACTION] "reaction"

	GroupId    int `json:"group_id" mapstructure:"group_id"`
	MessageId  int `json:"message_id" mapstructure:"message_id"`
	OperatorId int `json:"operator_id" mapstructure:"operator_id"`
	// "add" / "remove"
	SubType string `json:"sub_type" mapstructure:"sub_type"`
	Code    string `json:"code" mapstructure:"code"`
	Count   int    `json:"count" mapstructure:"count"`
}

type NoticeOfflineFile struct {
	NoticeBase // [TYPE_L2_NOTICE_OFFLINE_FILE] "offline_file"

	File offlineFileInfo `json:"file" mapstructure:"file"`
}

type offlineFileInfo struct {
	Id   string `json:"id" mapstructure:"id"`
	Name string `json:"name" mapstructure:"name"`
	Size int    `json:"size" mapstructure:"size"`
	Url  string `json:"url" mapstructure:"url"`
	Hash string `json:"hash" mapstructure:"hash"`
}
