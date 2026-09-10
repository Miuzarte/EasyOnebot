package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/sync/singleflight"
)

type Callable interface {
	PostReq(*Request) (*Response, error) // CallApi 调用 API 返回响应
	PostReqNoResp(*Request) error        // CallApiNoResp 调用 API 屏蔽响应
}

// Caller 只负责序列化请求和解析响应, 响应中 retcode 不为 0 之类的错误由 [Callable] 处理
type (
	// [StdCaller] implements OneBot standard API
	StdCaller struct{ Callable }
	// [NcCaller] implements NapCat extension API
	//
	// 包含 NapCat 兼容的 Lagrange 扩展端点 (合并转发、戳一戳、群文件等)。
	NcCaller struct{ Callable }
)

// 参数验证
var validate = customValidator{validator.New(validator.WithRequiredStructEnabled())}

type customValidator struct {
	*validator.Validate
}

// userId, groupId ...
func (cv *customValidator) idGt0(id int) error {
	return cv.Var(id, "gt=0")
}

// messageId
func (cv *customValidator) msgIdNe0(messageId int) error {
	return cv.Var(messageId, "ne=0")
}

func (cv *customValidator) require(v any) error {
	return cv.Var(v, "required")
}

// 为类似 [GetMsg] 的 api 实现 singleflight
var sfg singleflight.Group

// DecodeResponseAssertion 断言基础类型
func DecodeResponseAssertion[T any](r *Response, err error) (T, error) {
	if err != nil {
		return *new(T), err
	}
	resp, ok := r.Data.(T)
	if !ok {
		pT := new(T)
		return *pT, fmt.Errorf("can not assert data(type of %T) to %T", r.Data, *pT)
	}
	return resp, nil
}

// DecodeResponse 解析结构体
func DecodeResponse[T any](r *Response, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	resp := new(T)
	err = mapstructure.Decode(r.Data, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DecodeResponseSlice 解析切片
func DecodeResponseSlice[ST ~[]T, T any](r *Response, err error) (ST, error) {
	if err != nil {
		return nil, err
	}
	resp := make(ST, 0)
	err = mapstructure.Decode(r.Data, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
