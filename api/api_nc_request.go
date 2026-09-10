package api

/*
SetFriendAddRequest 设置好友添加请求

std: [SetFriendAddRequest]

见 NapCat 端点 [set_friend_add_request] (docs/onebot-napcat-endpoints.md)

参数:

	flag: 要设置的请求的 Flag
	approve: 是否同意 <默认值: true>
	reason(可选): 拒绝原因
*/
func SetFriendAddRequest_Nc(flag string, approve bool, reason ...string) *Request {
	if len(reason) == 0 {
		return NewReq("set_friend_add_request", map[string]any{
			"flag":    flag,
			"approve": approve,
		})
	} else {
		return NewReq("set_friend_add_request", map[string]any{
			"flag":    flag,
			"approve": approve,
			"reason":  reason[0],
		})
	}
}

type validateSetFriendAddRequest struct {
	Flag   string   `validate:"required"`
	Reason []string `validate:"omitempty"`
}

func (c NcCaller) SetFriendAddRequest(flag string, approve bool, reason ...string) error {
	err := validate.Struct(&validateSetFriendAddRequest{flag, reason})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetFriendAddRequest_Nc(flag, approve, reason...))
}

/*
SetGroupAddRequest 处理加群请求／邀请

std: [SetGroupAddRequest]

见 NapCat 端点 [set_group_add_request] (docs/onebot-napcat-endpoints.md)

参数:

	flag: 要设置的请求的 Flag
	approve: 是否同意 <默认值: true>
	reason(可选): 拒绝原因
*/
func SetGroupAddRequest_Nc(flag string, approve bool, reason ...string) *Request {
	if len(reason) == 0 {
		return NewReq("set_group_add_request", map[string]any{
			"flag":    flag,
			"approve": approve,
		})
	} else {
		return NewReq("set_group_add_request", map[string]any{
			"flag":    flag,
			"approve": approve,
			"reason":  reason[0],
		})
	}
}

type validateSetGroupAddRequest_Nc struct {
	Flag   string   `validate:"required"`
	Reason []string `validate:"omitempty"`
}

func (c NcCaller) SetGroupAddRequest(flag string, approve bool, reason ...string) error {
	err := validate.Struct(&validateSetGroupAddRequest_Nc{flag, reason})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupAddRequest_Nc(flag, approve, reason...))
}
