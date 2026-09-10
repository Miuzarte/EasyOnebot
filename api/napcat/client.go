package napcat

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Miuzarte/EasyOnebot/api"
)

// Poster 是调用方的传输能力, 由 EasyOnebot 的 Bot 实现。
//
// 只依赖这两个方法, 便于测试替身与其它传输实现。
type Poster interface {
	PostReq(*api.Request) (*api.Response, error)
}

// Caller 是 NapCat 端点的类型化调用器。
//
// 端点名来自 spec (见 [EndpointToAction]), 参数与响应都是生成类型,
// 响应 data 用严格 json 解码 —— NapCat 返回的结构与 spec 不符时直接报错,
// 而不是像 mapstructure 那样静默凑合。
type Caller struct {
	Poster
}

// endpoints 是 operationId -> 端点名, 进程启动时从 embed 的 spec 读一次。
//
// spec 已编进二进制, 这里的失败只可能是构建期问题, 所以直接 panic 而不是把错误
// 传染给每个调用点。
var endpoints = func() map[string]string {
	m, err := EndpointToAction()
	if err != nil {
		panic("napcat: 读取 spec 失败: " + err.Error())
	}
	return m
}()

// Call 按 operationId 调用端点, 并把响应 data 严格解码为 D。
//
// 返回的 error 包含两类失败: 传输/retcode 层失败, 以及 data 与 D 类型不匹配。
func Call[P any, D any](p Poster, operationID string, params P) (D, error) {
	var zero D

	action, ok := endpoints[operationID]
	if !ok {
		return zero, fmt.Errorf("napcat: 未知 operationId %q (spec 里没有该端点)", operationID)
	}

	raw, err := json.Marshal(params)
	if err != nil {
		return zero, fmt.Errorf("napcat: 编码 %s 参数失败: %w", action, err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return zero, fmt.Errorf("napcat: %s 参数不是对象: %w", action, err)
	}

	resp, err := p.PostReq(api.NewReq(action, m))
	if err != nil {
		return zero, err
	}
	if resp == nil || resp.Data == nil {
		// 无 data 的端点 (例如 delete_msg) 允许调用方用 D = any 接收
		return zero, nil
	}

	// 统一经 json.RawMessage 走一遍, 保证严格解码不受传输层解码方式影响
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return zero, fmt.Errorf("napcat: %s 响应无法重新编码: %w", action, err)
	}
	if err := json.Unmarshal(data, &zero); err != nil {
		if ute, ok := errors.AsType[*json.UnmarshalTypeError](err); ok {
			return zero, &DecodeError{Action: action, Field: ute.Field, Want: ute.Type.String(), Got: ute.Value, Err: err}
		}
		return zero, fmt.Errorf("napcat: 解码 %s 响应失败: %w", action, err)
	}
	return zero, nil
}

// DecodeError 表示响应结构与生成类型不匹配, 通常意味着 NapCat 改了返回或 spec 落后。
type DecodeError struct {
	Action string
	Field  string
	Want   string
	Got    string
	Err    error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("napcat: %s 响应字段 %s 类型不匹配: 期望 %s, 实际 %s (spec 可能落后于 NapCat)",
		e.Action, e.Field, e.Want, e.Got)
}

func (e *DecodeError) Unwrap() error { return e.Err }
