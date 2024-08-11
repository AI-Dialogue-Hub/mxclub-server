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
	xjet.NewCommonJetController[CommonController]()
}

type CommonController struct {
	jet.BaseJetController
}

type Param struct {
	Name string `form:"name" validate:"required" reg_err_info:"name字段不能为空"`
	Age  int    `form:"age" validate:"lte=120,gte=1" reg_err_info:"年龄不合法"`
}

func (*CommonController) GetV1Hello(ctx jet.Ctx, p *Param) (*api.Response, error) {
	ctx.Logger().Infof("request uri:%v", string(ctx.Request().RequestURI()))
	ctx.Logger().Infof("request param:%v", utils.ObjToJsonStr(p))
	return api.Success(ctx.Logger().ReqId, "hello world"), nil
}

func (*CommonController) PostV1Upload(ctx jet.Ctx) (*api.Response, error) {
	strategy, _ := xupload.FetchStrategy(config.GetConfig().UploadConfig.StorageType)
	fileURL, err := strategy.PutFromWeb(ctx)
	if err != nil {
		ctx.Logger().Errorf("PutFromWeb ERROR:%v", err)
	}
	return xjet.WrapperResult(ctx, fileURL, nil)
}

func (ctr ProductController) GetV1File0(ctx jet.Ctx, params *api.PathParam) {
	storageType := config.GetConfig().UploadConfig.StorageType
	strategy, err := xupload.FetchStrategy(storageType)
	if err != nil {
		ctx.Logger().Errorf("FetchStrategy ERROR:%v", err)
		return
	}
	path, _ := params.GetString(0)
	if storageType == "oss" {
		localStrategy, _ := xupload.FetchLocalStrategy()
		localStrategy.GetFile(ctx, path, config.GetConfig().File.FilePath)
	} else {
		strategy.GetFile(ctx, path, config.GetConfig().File.FilePath)
	}
}
