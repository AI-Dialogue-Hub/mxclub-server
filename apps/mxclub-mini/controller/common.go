package controller

import (
	"github.com/fengyuan-liang/GoKit/utils"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/config"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/common/xupload"
)

func init() {
	xjet.NewCommonJetController[DemoController]()
}

type DemoController struct {
	jet.BaseJetController
}

type Param struct {
	Name string `form:"name" validate:"required" reg_err_info:"name字段不能为空"`
	Age  int    `form:"age" validate:"lte=120,gte=1" reg_err_info:"年龄不合法"`
}

func (*DemoController) GetV1Hello(ctx jet.Ctx, p *Param) (*api.Response, error) {
	ctx.Logger().Infof("request uri:%v", string(ctx.Request().RequestURI()))
	ctx.Logger().Infof("request param:%v", utils.ObjToJsonStr(p))
	return api.Success(ctx.Logger().ReqId, "hello world"), nil
}

func (*DemoController) PostV1Upload(ctx jet.Ctx) (*api.Response, error) {
	strategy, _ := xupload.FetchStrategy(config.GetConfig().UploadConfig.StorageType)
	fileURL, err := strategy.PutFromWeb(ctx)
	if err != nil {
		ctx.Logger().Errorf("PutFromWeb ERROR:%v", err)
	}
	return xjet.WrapperResult(ctx, fileURL, nil)
}

func (*DemoController) GetV1File0(ctx jet.Ctx, params *api.PathParam) {
	strategy, err := xupload.FetchStrategy(config.GetConfig().UploadConfig.StorageType)
	if err != nil {
		ctx.Logger().Errorf("FetchStrategy ERROR:%v", err)
		return
	}
	path, _ := params.GetString(0)

	strategy.GetFile(ctx, config.GetConfig().File.FilePath, path)
}
