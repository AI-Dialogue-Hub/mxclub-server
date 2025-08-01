package xjet

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/core/context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
	"net"
	"net/http"
	"strings"
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

func WraResult(ctx jet.Ctx, result any, err error) (*api.R[any], error) {
	if err != nil {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, err.Error())
	}
	return api.Ok(ctx.Logger().ReqId, result), nil
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
func ConvertFastHTTPRequestToStandard(ctx jet.Ctx) (*http.Request, error) {
	request := new(http.Request)
	err := fasthttpadaptor.ConvertRequest(ctx.FastHttpCtx(), request, true)
	if err != nil {
		ctx.Logger().Errorf("ConvertFastHTTPRequestToStandard ERROR:%v", err)
		return nil, err
	}
	return request, nil
}

var defaultCtx = context.NewContext(new(fasthttp.RequestCtx), xlog.NewWith("defaultCtx"))

func NewDefaultJetContext() jet.Ctx {
	return defaultCtx
}

func GetClientIP(ctx jet.Ctx) string {
	defer utils.RecoverAndLogError(ctx)
	fastHttpCtx := ctx.FastHttpCtx()
	// 尝试从X-Forwarded-For获取
	if ip := string(fastHttpCtx.Request.Header.Peek("X-Forwarded-For")); ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// 尝试从X-Real-IP获取
	if ip := string(fastHttpCtx.Request.Header.Peek("X-Real-IP")); ip != "" {
		return ip
	}

	// 从RemoteAddr获取
	ip, _, err := net.SplitHostPort(fastHttpCtx.RemoteAddr().String())
	if err != nil {
		return fastHttpCtx.RemoteAddr().String()
	}
	return ip
}
