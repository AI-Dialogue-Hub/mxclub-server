package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/config"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/common/xupload"
)

func (*DemoController) PostV1Upload(ctx jet.Ctx) (*api.Response, error) {
	strategy, _ := xupload.FetchStrategy(config.GetConfig().UploadConfig.StorageType)
	fileURL, err := strategy.PutFromWeb(ctx)
	if err != nil {
		ctx.Logger().Errorf("PutFromWeb ERROR:%v", err)
	}
	return xjet.WrapperResult(ctx, fileURL, nil)
}

func (*DemoController) GetV1File0(ctx jet.Ctx, params *api.PathParam) {
	strategy, _ := xupload.FetchStrategy(config.GetConfig().UploadConfig.StorageType)
	if config.GetConfig().UploadConfig.StorageType == "oss" {
		filePath, _ := params.GetString(0)
		// 从本地上传到oss中
		fileURL, _ := strategy.PutFromLocalFile(ctx, config.GetConfig().File.FilePath+"/"+filePath)
		ctx.Logger().Infof("parse success:%v", fileURL)
	}
}
