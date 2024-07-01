package xjet

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/pkg/api"
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

func IsNil(val any) bool {
	if val == nil {
		return true
	}
	return false
}
