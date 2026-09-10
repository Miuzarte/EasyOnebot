package api

/*
SetFriendAddRequest 设置好友添加请求

std: [SetFriendAddRequest]

https://lagrange-onebot.apifox.cn/236975604e0

参数:

	flag: 要设置的请求的 Flag
	approve: 是否同意 <默认值: true>
	reason(可选): 拒绝原因
*/
func SetFriendAddRequest_Lgr(flag string, approve bool, reason ...string) *Request {
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

type validateSetFriendAddRequest_Lgr struct {
	Flag   string   `validate:"required"`
	Reason []string `validate:"omitempty"`
}

func (c LgrCaller) SetFriendAddRequest_Lgr(flag string, approve bool, reason ...string) error {
	err := validate.Struct(&validateSetFriendAddRequest_Lgr{flag, reason})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetFriendAddRequest_Lgr(flag, approve, reason...))
}

/*
SetGroupAddRequest 处理加群请求／邀请

std: [SetGroupAddRequest]

https://lagrange-onebot.apifox.cn/236975617e0

参数:

	flag: 要设置的请求的 Flag
	approve: 是否同意 <默认值: true>
	reason(可选): 拒绝原因
*/
func SetGroupAddRequest_Lgr(flag string, approve bool, reason ...string) *Request {
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

type validateSetGroupAddRequest_Lgr struct {
	Flag   string   `validate:"required"`
	Reason []string `validate:"omitempty"`
}

func (c LgrCaller) SetGroupAddRequest_Lgr(flag string, approve bool, reason ...string) error {
	err := validate.Struct(&validateSetGroupAddRequest_Lgr{flag, reason})
	if err != nil {
		return err
	}
	return c.PostReqNoResp(SetGroupAddRequest_Lgr(flag, approve, reason...))
}
