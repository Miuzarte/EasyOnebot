package api

/*
GetFriendList_Lgr 获取好友列表

std: [GetFriendList]

https://lagrange-onebot.apifox.cn/236981445e0
*/
func GetFriendList_Lgr() *Request {
	return NewReq("get_friend_list", nil)
}

type GetFriendListResp_Lgr []struct {
	UserId   int    `json:"user_id" mapstructure:"user_id"`   // Uin
	Qid      string `json:"q_id" mapstructure:"q_id"`         // QID (可选)
	Nickname string `json:"nickname" mapstructure:"nickname"` // 昵称
	Group    struct {
		GroupId   int    `json:"group_id" mapstructure:"group_id"`     // 分组 ID
		GroupName string `json:"group_name" mapstructure:"group_name"` // 分组名称
	} `json:"group" mapstructure:"group"` // 分组信息
	Remark string `json:"remark" mapstructure:"json"` // 备注
}

// func GetGroupInfo(groupId int, noCache ...bool) *Request
// std: [GetGroupInfo] https://lagrange-onebot.apifox.cn/236981449e0
// func GetGroupList(noCache ...bool) *Request
// std: [GetGroupList] https://lagrange-onebot.apifox.cn/236981452e0
// func GetGroupMemberInfo(groupId, userId int, noCache ...bool) *Request
// std: [GetGroupMemberInfo] https://lagrange-onebot.apifox.cn/236981473e0
// func GetGroupMemberList(groupId int) *Request
// std: [GetGroupMemberList] https://lagrange-onebot.apifox.cn/236981479e0
// func GetLoginInfo() *Request
// std: [GetLoginInfo] https://lagrange-onebot.apifox.cn/236981487e0
// func GetStatus_Lgr() *Request
// std: [GetStatus] https://lagrange-onebot.apifox.cn/236981584e0

/*
GetStrangerInfo_Lgr 获取陌生人信息

std: [GetStrangerInfo]

https://lagrange-onebot.apifox.cn/236981590e0

参数:

	userId: 用户 Uin
	noCache: <默认值: false>
*/
func GetStrangerInfo_Lgr(userId int, noCache ...bool) *Request {
	if len(noCache) == 0 {
		return NewReq("get_stranger_info", map[string]any{
			"user_id": userId,
		})
	} else {
		return NewReq("get_stranger_info", map[string]any{
			"user_id":  userId,
			"no_cache": noCache[0],
		})
	}
}

type validateGetStrangerInfo_Lgr struct {
	UserId  int    `validate:"gt=0"`
	NoCache []bool `validate:"omitempty"`
}

type GetStrangerInfoResp_Lgr struct {
	UserId       int         `json:"user_id" mapstructure:"user_id"`           // 用户 Uin
	Qid          string      `json:"q_id" mapstructure:"q_id"`                 // QID (可选)
	Nickname     string      `json:"nickname" mapstructure:"nickname"`         // 昵称
	Avatar       string      `json:"avatar" mapstructure:"avatar"`             // 头像
	Level        int         `json:"level" mapstructure:"level"`               // 等级
	Age          int         `json:"age" mapstructure:"age"`                   // 年龄
	Sex          string      `json:"sex" mapstructure:"sex"`                   // 性别
	RegisterTime string      `json:"RegisterTime" mapstructure:"RegisterTime"` // 注册时间
	Sign         string      `json:"sign" mapstructure:"sign"`                 // 个性签名
	Status       StatusClass `json:"status" mapstructure:"status"`             // 当前状态信息
	Business     []Business  `json:"Business" mapstructure:"Business"`
}

type StatusClass struct {
	FaceID   int    `json:"face_id" mapstructure:"face_id"`     // 表情 ID (可选)
	Message  string `json:"message" mapstructure:"message"`     // 信息 (可选)
	StatusID int    `json:"status_id" mapstructure:"status_id"` // 状态 ID
}

type Business struct {
	Icon   *string `json:"icon" mapstructure:"icon"`
	Ispro  int     `json:"ispro" mapstructure:"ispro"`
	Isyear int     `json:"isyear" mapstructure:"isyear"`
	Level  int     `json:"level" mapstructure:"level"`
	Name   *string `json:"name" mapstructure:"name"`
	Type   int     `json:"type" mapstructure:"type"`
}

// GetStrangerInfo_Lgr 获取陌生人信息
func (c LgrCaller) GetStrangerInfo_Lgr(userId int, noCache ...bool) (*GetStrangerInfoResp_Lgr, error) {
	err := validate.Struct(&validateGetStrangerInfo_Lgr{userId, noCache})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[GetStrangerInfoResp_Lgr](c.PostReq(GetStrangerInfo_Lgr(userId, noCache...)))
}

// func GetVersionInfo() *Request
// std: [GetVersionInfo] https://lagrange-onebot.apifox.cn/236981622e0
