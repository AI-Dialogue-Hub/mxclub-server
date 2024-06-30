package controller

import (
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/valyala/fasthttp"
	"io"
	"math/rand"
	"mime"
	"mxclub/apps/mxclub-admin/config"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (*DemoController) PostV1Upload(ctx jet.Ctx) (*api.Response, error) {
	var (
		logger = ctx.Logger()
	)
	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		logger.Error("Error retrieving the file", fasthttp.StatusBadRequest)
		return xjet.WrapperResult(ctx, nil, err)
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		logger.Error("Error opening the file", fasthttp.StatusInternalServerError)
		return xjet.WrapperResult(ctx, nil, err)
	}
	defer src.Close()

	// 读取文件内容
	fileData, err := io.ReadAll(src)
	if err != nil {
		logger.Error("Error reading the file", fasthttp.StatusInternalServerError)
		return xjet.WrapperResult(ctx, nil, err)
	}

	// 生成文件名
	fileName := generateFileName(file.Filename)

	// 拼接文件保存路径
	savePath := filepath.Join(config.GetConfig().File.FilePath, fileName)

	// 创建保存目录
	err = os.MkdirAll(filepath.Dir(savePath), 0755)
	if err != nil {
		logger.Error("Error creating the directory", fasthttp.StatusInternalServerError)
		return xjet.WrapperResult(ctx, nil, err)
	}

	// 创建一个新的文件来保存上传的文件
	dst, err := os.Create(savePath)
	if err != nil {
		logger.Error("Error creating the destination file", fasthttp.StatusInternalServerError)
		return xjet.WrapperResult(ctx, nil, err)
	}
	defer dst.Close()

	// 将文件数据写入目标文件
	_, err = dst.Write(fileData)
	if err != nil {
		logger.Error("Error writing the file", fasthttp.StatusInternalServerError)
		return xjet.WrapperResult(ctx, nil, err)
	}

	// 构建文件URI
	fileURI := fmt.Sprintf("%s%s", config.GetConfig().File.Domain, fileName)

	return xjet.WrapperResult(ctx, fileURI, nil)
}

func generateFileName(originalName string) string {
	// rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(100000)
	currentTime := time.Now().Format("20060102150405")

	fileExt := filepath.Ext(originalName)
	fileName := fmt.Sprintf("%s%d%s", currentTime, randomNumber, fileExt)

	return fileName
}

func (*DemoController) GetV1File0(ctx jet.Ctx, params *jet.Args) {
	var (
		logger = ctx.Logger()
		path   = params.CmdArgs[0]
	)
	if path == "" {
		logger.Errorf("path is empty, url:%v", path)
		return
	}
	fileName := filepath.Base(path)
	filePath := filepath.Join(config.GetConfig().File.FilePath, fileName)
	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if err != nil {
		logger.Error("File not found", fasthttp.StatusNotFound)
		xjet.Error(ctx, "File not found", fasthttp.StatusNotFound)
		return
	}

	// 读取文件内容
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error("Error reading the file", fasthttp.StatusInternalServerError)
		xjet.Error(ctx, "Error reading the file", fasthttp.StatusInternalServerError)
		return
	}

	var mineType = mime.TypeByExtension(filepath.Ext(filePath))

	ctx.Response().Header.Set("Content-Length", strconv.Itoa(len(fileData)))
	// 根据文件类型设置响应头
	if strings.HasPrefix(mime.TypeByExtension(filepath.Ext(filePath)), "image/") {
		ctx.Response().Header.Set("Content-Disposition", "inline")
		ctx.Response().Header.Set("Content-Type", mineType)
	} else {
		ctx.Response().Header.Set("Content-Disposition", "attachment; filename="+fileName)
		ctx.Response().Header.Set("Content-Type", "application/octet-stream")
	}
	ctx.Response().SetBody(fileData)
}
