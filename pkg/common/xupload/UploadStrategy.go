package xupload

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/pkg/common/xupload/local"
	"mxclub/pkg/common/xupload/xoss"
)

type UploadStrategy interface {
	PutFromLocalFile(ctx jet.Ctx, filePath string) (string, error)
	PutFromWeb(ctx jet.Ctx) (string, error)
	// GetFile
	// @param path 文件名称
	// @param filePath 文件保存路径
	GetFile(ctx jet.Ctx, path string, filePath string)
}

var _ UploadStrategy = (*local.LocalUploadStrategy)(nil)
var _ UploadStrategy = (*xoss.OssUploadStrategy)(nil)

var UploadStrategyFactory = map[string]UploadStrategy{
	"local": new(local.LocalUploadStrategy),
	"oss":   new(xoss.OssUploadStrategy),
}

func FetchStrategy(storageType string) (UploadStrategy, error) {
	strategy, ok := UploadStrategyFactory[storageType]
	if !ok {
		return nil, errors.New("unsupported storage type")
	}
	if strategy == nil {
		return UploadStrategyFactory["local"], nil
	}
	return strategy, nil
}
