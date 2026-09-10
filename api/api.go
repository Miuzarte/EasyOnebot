package api

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
	"time"
)

var echoSeq uintptr

// NewReq 创建一个 API 调用
func NewReq(action string, params map[string]any) *Request {
	ac := &Request{Action: action}
	if len(params) > 0 {
		ac.Params = params
	}

	ac.Echo = ac.Action +
		"_" + strconv.FormatInt(time.Now().UnixNano(), 10) +
		"_" + strconv.FormatUint(uint64(atomic.AddUintptr(&echoSeq, 1)), 16)

	return ac
}

/*
Request 是 API 调用的原始请求数据结构

https://github.com/botuniverse/onebot-11/blob/master/communication/ws.md#api-%E6%8E%A5%E5%8F%A3

这里的 `action` 参数用于指定要调用的 API，具体参考 [API](https://github.com/botuniverse/onebot-11/blob/master/api)；
`params` 用于传入参数，如果要调用的 API 不需要参数，则可以不加；
`echo` 字段是可选的，类似于 `JSON RPC` 的 `id` 字段，用于唯一标识一次请求，可以是任何类型的数据，OneBot 将会在调用结果中原样返回。
*/
type Request struct {
	Action string         `json:"action"`
	Params map[string]any `json:"params"`
	Echo   string         `json:"echo"`
} // 仅编码, 不需要 mapstructure

func (r *Request) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

/*
Response 是 API 调用的原始响应数据结构

https://github.com/botuniverse/onebot-11/blob/master/communication/ws.md#api-%E6%8E%A5%E5%8F%A3

客户端向 OneBot 发送 JSON 之后，OneBot 会往回发送一个调用结果，结构和 `HTTP 的响应` 相似，（除了包含请求中传入的 `echo` 字段外）
唯一的区别在于，通过 HTTP 调用 API 时，HTTP 状态码反应的错误情况被移动到响应 JSON 的 `retcode` 字段，
例如，HTTP 返回 404 的情况，对应到 WebSocket 的回复，是：

	{
		"status": "failed",
		"retcode": 1404,
		"data": null,
		"echo": "123"
	}
*/
type Response struct {
	Status  string `json:"status" mapstructure:"status"`
	RetCode int    `json:"retcode" mapstructure:"retcode"`
	Data    any    `json:"data" mapstructure:"data"` // map[string]any / []any
	Echo    string `json:"echo" mapstructure:"echo"`
} // 从 map[string]any / []any 解码, 需要 mapstructure

func (r *Response) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
