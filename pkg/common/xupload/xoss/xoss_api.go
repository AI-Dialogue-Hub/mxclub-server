package xoss

import (
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/valyala/fasthttp"
	"mime/multipart"
	"mxclub/pkg/utils"
	"os"
)

type OssUploadStrategy struct{}

func (*OssUploadStrategy) PutFromLocalFile(ctx jet.Ctx, filePath string) (string, error) {
	logger := ctx.Logger()
	fileUrlPath := fmt.Sprintf("%v/%v", ossCfg.StoragePath, filePath)
	appendFile, err := client.AppendFile(todoCtx, ossCfg.Bucket, fileUrlPath)
	if err != nil {
		logger.Errorf("failed to append file %v", err)
		return "", err
	}
	readFile, _ := os.ReadFile(filePath)
	logger.Infof("append file af:%#v\n", appendFile)
	n, err := appendFile.Write(readFile)
	if err != nil {
		logger.Errorf("failed to af write %v", err)
		return "", err
	}
	defer utils.HandleClose(appendFile)
	logger.Infof("af write n:%#v\n", n)
	return BuildURL(fileUrlPath), nil
}

func (*OssUploadStrategy) PutFromWeb(ctx jet.Ctx) (string, error) {
	var (
		logger = ctx.Logger()
	)
	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		logger.Error("Error retrieving the file", fasthttp.StatusBadRequest)
		return "", err
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		logger.Error("Error opening the file", fasthttp.StatusInternalServerError)
		return "", err
	}

	defer func(src multipart.File) {
		err = src.Close()
		if err != nil {
			logger.Errorf("src close error:%v", err.Error())
		}
	}(src)

	filePath := fmt.Sprintf("%s/%s", ossCfg.StoragePath, utils.GenerateFileName(file.Filename))

	request := &oss.AppendObjectRequest{
		Bucket:   oss.Ptr(ossCfg.Bucket),
		Key:      oss.Ptr(filePath),
		Position: oss.Ptr(int64(0)),
		Body:     src,
	}

	result, err := client.AppendObject(todoCtx, request)

	if err != nil {
		logger.Errorf("AppendObject ERROR:%v", err)
	}

	logger.Infof("append object result:%#v\n", result)

	return BuildURL(filePath), nil
}

func (*OssUploadStrategy) GetFile(ctx jet.Ctx, path string, filePath string) {
	// noting to do
}
