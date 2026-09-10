package api

import "github.com/Miuzarte/EasyOnebot/message"

/*
/send_forward_msg 发送合并转发消息

https://napcat.apifox.cn/226659136e0

参数:

	groupId
	userId
	messages: 合并转发消息
	news: []string
	prompt: 外显
	summary: 底下文本
	source: 内容
*/
func SendForwardMsg_Nc(groupId, userId int, messages message.SegmentArray, news []string, prompt, summary, source string) *Request {
	type NewsItem struct {
		Text string `json:"text"`
	}

	params := map[string]any{}
	if groupId > 0 {
		params["group_id"] = groupId
	}
	if userId > 0 {
		params["user_id"] = userId
	}
	params["messages"] = messages

	if len(news) > 0 {
		newsItems := make([]NewsItem, 0, len(news))
		for _, n := range news {
			newsItems = append(newsItems, NewsItem{Text: n})
		}
		params["news"] = newsItems
	}
	if prompt != "" {
		params["prompt"] = prompt
	}
	if summary != "" {
		params["summary"] = summary
	}
	if source != "" {
		params["source"] = source
	}

	return NewReq("send_forward_msg", params)
}

type validateSendForwardMsg_Nc struct {
	// GroupId  int                  `validate:"required_without=UserId,gt=0"`
	// UserId   int                  `validate:"required_without=GroupId,gt=0"`
	Messages message.SegmentArray `validate:"required,min=1"`
}

type SendForwardMsg_NcResp struct {
	MessageId int    `json:"message_id" mapstructure:"message_id"`
	ForwardId string `json:"forward_id" mapstructure:"forward_id"`
}

func (c NcCaller) SendForwardMsg(groupId, userId int, messages message.SegmentArray, news []string, prompt, summary, source string) (*SendForwardMsg_NcResp, error) {
	err := validate.Struct(&validateSendForwardMsg_Nc{messages})
	if err != nil {
		return nil, err
	}
	return DecodeResponse[SendForwardMsg_NcResp](c.PostReq(SendForwardMsg_Nc(groupId, userId, messages, news, prompt, summary, source)))
}
