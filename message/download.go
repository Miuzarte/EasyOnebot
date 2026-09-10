package message

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/sync/singleflight"
)

var sfg singleflight.Group

type finResponse struct {
	response *http.Response
	body     []byte
}

func sfGet(ctx context.Context, url string) (resp *finResponse, err error, shared bool) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err, false
	}
	return sfDo(req)
}

func sfDo(req *http.Request) (resp *finResponse, err error, shared bool) {
	url := req.URL.String() // 用 url 做标识, POST 不能这么搞

	v, err, shared := sfg.Do(url, func() (any, error) {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("bad status code: %s", resp.Status)
		}
		if err != nil {
			sfg.Forget(url)
		}
		return &finResponse{
			response: resp,
			body:     body,
		}, err
	})

	if err != nil {
		return nil, err, shared
	}
	return v.(*finResponse), err, shared
}

func (seg Segment) GetUrl() (url string, err error) {
	switch seg.Type {
	case TYPE_IMAGE,
		TYPE_VIDEO,
		TYPE_RECORD,
		TYPE_FILE,
		TYPE_SHARE:

		if u := seg.Data["url"].(string); u != "" && strings.HasPrefix(u, "http") {
			return u, nil
		}
		if f := seg.Data["file"].(string); f != "" && strings.HasPrefix(f, "http") {
			return f, nil
		}
		return "", fmt.Errorf("failed to get url from seg: %s", seg)
	}
	return "", fmt.Errorf("%s download not implemented", seg.Type)
}

func (seg Segment) Download(ctx context.Context) (body []byte, err error) {
	url, err := seg.GetUrl()
	if err != nil {
		return nil, err
	}
	resp, err, _ := sfGet(ctx, url)
	if resp != nil {
		body = resp.body
	}
	return body, err
}
