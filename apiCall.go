package EasyOnebot

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/Miuzarte/EasyOnebot/api"
	"github.com/Miuzarte/EasyOnebot/api/napcat"
)

// MixCaller combines different implementations
type MixCaller struct {
	Std api.StdCaller
	Nc  api.NcCaller
}

// Call returns a [MixCaller] to call OneBot APIs
func (b *Bot) Call() MixCaller {
	return MixCaller{api.StdCaller{Callable: b}, api.NcCaller{Callable: b}}
}

// NapCat 返回基于本连接的 NapCat 类型化调用器。
//
// 与 [MixCaller] 的区别: 这里的参数与响应都是 OpenAPI spec 生成的类型,
// 响应 data 严格解码 (字段类型与 spec 不符时直接报错), 端点名来自 spec 查表。
// 新代码优先用它; Call() 保留给 OneBot 11 标准端点的旧写法。
func (b *Bot) NapCat() napcat.Caller {
	return napcat.Caller{Poster: b}
}

// PostReq implements [api.Callable]
func (b *Bot) PostReq(req *api.Request) (*api.Response, error) {
	if b.log.Trace().Enabled() {
		debug.PrintStack()
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = b.post(data)
	if err != nil {
		return nil, err
	}

	if b.apiPool.timeout == 0 {
		b.log.Error().Msg("apiPool.timeout is 0")
		b.log.Error().Msg(string(debug.Stack()))
		b.apiPool.timeout = time.Minute * 5
	}

	ctx, cancel := context.WithTimeout(context.Background(), b.apiPool.timeout)
	ch := make(chan *api.Response, 1)
	b.apiPool.put(req.Echo, ch)
	defer b.apiPool.del(req.Echo) // 清理
	defer close(ch)
	defer cancel()

	select {
	case resp := <-ch:
		switch {
		case resp.RetCode == 0 && resp.Status == "ok":
		case resp.RetCode == 1 && resp.Status == "async":
		default:
			// 理论上 resp 不会为 nil
			return resp, fmt.Errorf("failed to call api [%s]: %s", req.Action, resp)
		}
		return resp, nil
	case <-ctx.Done(): // 超时
		return nil, fmt.Errorf("api [%s] call timeout", req.Action)
	}
}

// PostReqNoResp implements [api.Callable]
func (b *Bot) PostReqNoResp(req *api.Request) (err error) {
	_, err = b.PostReq(req)
	return err
}
