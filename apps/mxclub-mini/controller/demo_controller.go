package controller

import (
	"github.com/fengyuan-liang/GoKit/utils"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/pkg/api"
)

func init() {
	jet.Provide(func() jet.ControllerResult { return jet.NewJetController(new(DemoController)) })
}

type DemoController struct {
	jet.BaseJetController
}

type Param struct {
	Name string `form:"name" validate:"required" reg_err_info:"name字段不能为空"`
	Age  int    `form:"age" validate:"lte=120,gte=0" reg_err_info:"年龄不合法"`
}

func (*DemoController) GetV1Hello(ctx jet.Ctx, p *Param) *api.Response {
	ctx.Logger().Infof("request uri:%v", string(ctx.Request().RequestURI()))
	ctx.Logger().Infof("request param:%v", utils.ObjToJsonStr(p))
	return api.Success(ctx.Logger().ReqId, "hello world")
}
