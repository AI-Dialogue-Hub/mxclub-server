package xjet

import (
	"bytes"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/valyala/fasthttp"
	"mxclub/pkg/api"
	"net/http"
	"net/url"
)

func NewCommonJetController[T jet.IJetController]() {
	jet.Provide(func() jet.ControllerResult { return jet.NewJetController(new(T)) })
}

type JetContext struct {
	Ctx jet.Ctx
}

func WrapperResult(ctx jet.Ctx, result any, err error) (*api.Response, error) {
	if err != nil {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, err.Error())
	}
	return api.Success(ctx.Logger().ReqId, result), nil
}

var defaultContentType = []byte("text/plain; charset=utf-8")

func Error(ctx jet.Ctx, msg string, statusCode int) {
	ctx.Response().Reset()
	ctx.Response().SetStatusCode(statusCode)
	ctx.Response().Header.SetContentTypeBytes(defaultContentType)
	ctx.Response().SetBodyString(msg)
}

func IsAnyEmpty(strs ...string) bool {
	if strs == nil || len(strs) == 0 {
		return true
	}
	for _, str := range strs {
		if str == "" {
			return true
		}
	}
	return false
}

// ConvertFastHTTPRequestToStandard converts a *fasthttp.Request to a *http.Request
func ConvertFastHTTPRequestToStandard(fastReq *fasthttp.Request) (*http.Request, error) {
	// 获取请求方法
	method := string(fastReq.Header.Method())

	// 获取请求URL
	uri := fastReq.URI()
	url := &url.URL{
		Scheme:   string(uri.Scheme()),
		Host:     string(uri.Host()),
		Path:     string(uri.Path()),
		RawPath:  string(uri.PathOriginal()),
		RawQuery: string(uri.QueryString()),
	}

	// 获取请求头部
	header := make(http.Header)
	fastReq.Header.VisitAll(func(key, value []byte) {
		header.Add(string(key), string(value))
	})

	// 获取请求正文并复制
	body := fastReq.Body()
	bodyCopy := make([]byte, len(body))
	copy(bodyCopy, body)
	bodyReader := bytes.NewReader(bodyCopy)

	// 创建新的http.Request实例
	req, err := http.NewRequest(method, url.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create http.Request: %w", err)
	}

	// 复制请求头部
	req.Header = header

	return req, nil
}
