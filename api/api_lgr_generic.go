package api

/*
FetchCustomFace 获取自定义Face

https://lagrange-onebot.apifox.cn/236974550e0
*/
func FetchCustomFace() *Request {
	return NewReq("fetch_custom_face", nil)
}

type FetchCustomFaceResp = []string // 表情 Url 列表

// FetchCustomFace 获取自定义Face
func (c LgrCaller) FetchCustomFace() (FetchCustomFaceResp, error) {
	return DecodeResponseSlice[FetchCustomFaceResp](c.PostReq(FetchCustomFace()))
}

/*
FetchMFaceKey 获取mface key

https://lagrange-onebot.apifox.cn/236974658e0

参数:

	emojiIds: 表情 Id 列表
*/
func FetchMFaceKey(emojiIds []string) *Request {
	return NewReq("fetch_mface_key", map[string]any{
		"emoji_ids": emojiIds,
	})
}

type FetchMFaceKeyResp = []string // 商城表情 Key 列表

// FetchMFaceKey 获取mface key
func (c LgrCaller) FetchMFaceKey(emojiIds []string) (FetchCustomFaceResp, error) {
	err := validate.require(emojiIds)
	if err != nil {
		return nil, err
	}
	return DecodeResponseSlice[FetchCustomFaceResp](c.PostReq(FetchMFaceKey(emojiIds)))
}

/*
JoinFriendEmojiChain 加入好友表情接龙

https://lagrange-onebot.apifox.cn/236974739e0

参数:

	userId: 用户 Uin
	messageId: 期望加入表情接龙的消息id
	emojiId: 表情id
*/
func JoinFriendEmojiChain(userId, messageId, emojiId int) *Request {
	return NewReq(".join_friend_emoji_chain", map[string]any{
		"user_id":    userId,
		"message_id": messageId,
		"emoji_id":   emojiId,
	})
}

type validateJoinFriendEmojiChain struct {
	UserId    int `validate:"gt=0"`
	MessageId int `validate:"ne=0"`
	EmojiId   int `validate:"required"`
}

// JoinFriendEmojiChain 加入好友表情接龙
func (c LgrCaller) JoinFriendEmojiChain(userId, messageId, emojiId int) error {
	err := validate.Struct(&validateJoinFriendEmojiChain{userId, messageId, emojiId})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(JoinFriendEmojiChain(userId, messageId, emojiId))
}

/*
GetAiCharacters 获取群 Ai 语音可用声色列表

https://lagrange-onebot.apifox.cn/236974974e0

参数:

	groupId: 群 Uin
	chatType(可选): 语音类型 <enum: 1, 2>
*/
func GetAiCharacters(groupId int, chatType ...int) *Request {
	if len(chatType) == 0 {
		return NewReq("get_ai_characters", map[string]any{
			"group_id": groupId,
		})
	} else {
		return NewReq("get_ai_characters", map[string]any{
			"group_id":  groupId,
			"chat_type": chatType[0],
		})
	}
}

type validateGetAiCharacters struct {
	GroupId  int   `validate:"gt=0"`
	ChatType []int `validate:"omitempty,div,oneof=1 2"`
}

type GetAiCharactersResp []struct {
	Type       string      `json:"type" mapstructure:"type"`             // Ai 声色分类
	Characters []Character `json:"characters" mapstructure:"characters"` // 分类下 Ai 声色列表
}

type Character struct {
	CharacterId   string `json:"character_id" mapstructure:"character_id"`     // Ai 声色 ID
	CharacterName string `json:"character_name" mapstructure:"character_name"` // Ai 声色名称
	PreviewUrl    string `json:"preview_url" mapstructure:"preview_url"`       // Ai 声色预览语音 Url
}

// GetAiCharacters 获取群 Ai 语音可用声色列表
func (c LgrCaller) GetAiCharacters(groupId int, chatType ...int) (GetAiCharactersResp, error) {
	err := validate.Struct(&validateGetAiCharacters{groupId, chatType})
	if err != nil {
		return nil, err
	}
	return DecodeResponseSlice[GetAiCharactersResp](c.PostReq(GetAiCharacters(groupId, chatType...)))
}

// func GetCookies(domain string) *Request
// std: [GetCookies] https://lagrange-onebot.apifox.cn/236975000e0
// func GetCredentials(domain string) *Request
// std: [GetCredentials] https://lagrange-onebot.apifox.cn/236975179e0
// func GetCsrfToken() *Request
// std: [GetCsrfToken] https://lagrange-onebot.apifox.cn/236975210e0

/*
JoinGroupEmojiChain 加入群聊表情接龙

https://lagrange-onebot.apifox.cn/236975310e0

参数:

	groupId: 群号
	messageId: 期望加入表情接龙的消息id
	emojiId: 表情id
*/
func JoinGroupEmojiChain(groupId, messageId, emojiId int) *Request {
	return NewReq(".join_group_emoji_chain", map[string]any{
		"group_id":   groupId,
		"message_id": messageId,
		"emoji_id":   emojiId,
	})
}

type validateJoinGroupEmojiChain struct {
	GroupId   int `validate:"gt=0"`
	MessageId int `validate:"ne=0"`
	EmojiId   int `validate:"required"`
}

// JoinGroupEmojiChain 加入群聊表情接龙
func (c LgrCaller) JoinGroupEmojiChain(groupId, messageId, emojiId int) error {
	err := validate.Struct(&validateJoinGroupEmojiChain{groupId, messageId, emojiId})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(JoinGroupEmojiChain(groupId, messageId, emojiId))
}

/*
OcrImage OCR图像识别

https://lagrange-onebot.apifox.cn/236975354e0

参数:

	image: image 链接, 支持 http/https/file/base64
*/
func OcrImage(image any) *Request {
	return NewReq("ocr_image", map[string]any{
		"image": image,
	})
}

type OcrImageResp struct {
	Language string `json:"language" mapstructure:"language"` // 语言
	Texts    []Text `json:"texts" mapstructure:"texts"`       // 文本信息列表
}

type Text struct {
	Confidence  int          `json:"confidence" mapstructure:"confidence"`   // 匹配率
	Coordinates []Coordinate `json:"coordinates" mapstructure:"coordinates"` // 位置列表
	Text        string       `json:"text" mapstructure:"text"`
}

type Coordinate struct {
	X int `json:"x" mapstructure:"x"`
	Y int `json:"y" mapstructure:"y"`
}

// OcrImage OCR图像识别
func (c LgrCaller) OcrImage(image any) (*OcrImageResp, error) {
	err := validate.require(image)
	if err != nil {
		return nil, err
	}
	return DecodeResponse[OcrImageResp](c.PostReq(OcrImage(image)))
}

/*
SetQqAvatar 设置QQ头像

https://lagrange-onebot.apifox.cn/236975389e0

参数:

	file: file 链接, 支持 http/https/file/base64
*/
func SetQqAvatar(file any) *Request {
	return NewReq("set_qq_avatar", map[string]any{
		"file": file,
	})
}

// SetQqAvatar 设置QQ头像
func (c LgrCaller) SetQqAvatar(file any) error {
	err := validate.require(file)
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetQqAvatar(file))
}

// func SendLike(userId, times int) *Request
// std: [SendLike] https://lagrange-onebot.apifox.cn/236975389e0
// func SetRestart() *Request
// std: [SetRestart] https://lagrange-onebot.apifox.cn/236975407e0

/*
DeleteFriend 删除好友

https://lagrange-onebot.apifox.cn/238992623e0

参数:

	userId: 用户 Uin
	block: 是否加入黑名单
*/
func DeleteFriend(userId int, block bool) *Request {
	return NewReq("delete_friend", map[string]any{
		"user_id": userId,
		"block":   block,
	})
}

// DeleteFriend 删除好友
func (c LgrCaller) DeleteFriend(userId int, block bool) error {
	err := validate.idGt0(userId)
	if err != nil {
		return err
	}
	return c.PostReqNoResp(DeleteFriend(userId, block))
}

/*
GetRKey 获取rkey

https://lagrange-onebot.apifox.cn/256115438e0
*/
func GetRKey() *Request {
	return NewReq("get_rkey", nil)
}

type GetRKeyResp struct {
	Rkeys []Rkey `json:"rkeys" mapstructureL:"rkeys"`
}

type Rkey struct {
	CreatedAt int    `json:"created_at" mapstructureL:"created_at"` // rkey的创建时间
	Rkey      string `json:"rkey" mapstructureL:"rkey"`             // rkey
	TTL       int    `json:"ttl" mapstructureL:"ttl"`               // rkey的有效时间
	Type      string `json:"type" mapstructureL:"type"`             // 类型，group, private 二者之一
}

// GetRKey 获取rkey
func (c LgrCaller) GetRKey() (*GetRKeyResp, error) {
	return DecodeResponse[GetRKeyResp](c.PostReq(GetRKey()))
}
