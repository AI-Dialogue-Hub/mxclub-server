package local

import (
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/valyala/fasthttp"
	"io"
	"math/rand"
	"mime"
	"mime/multipart"
	"mxclub/pkg/common/xjet"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type LocalUploadStrategy struct{}

func (*LocalUploadStrategy) PutFromLocalFile(ctx jet.Ctx, filePath string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (*LocalUploadStrategy) PutFromWeb(ctx jet.Ctx) (string, error) {
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

	// 读取文件内容
	fileData, err := io.ReadAll(src)
	if err != nil {
		logger.Error("Error reading the file", fasthttp.StatusInternalServerError)
		return "", err
	}

	// 生成文件名
	fileName := generateFileName(file.Filename)

	// 拼接文件保存路径
	savePath := filepath.Join(config.FilePath, fileName)

	// 创建保存目录
	err = os.MkdirAll(filepath.Dir(savePath), 0755)
	if err != nil {
		logger.Error("Error creating the directory", fasthttp.StatusInternalServerError)
		return "", err
	}

	// 创建一个新的文件来保存上传的文件
	dst, err := os.Create(savePath)
	if err != nil {
		logger.Error("Error creating the destination file", fasthttp.StatusInternalServerError)
		return "", err
	}
	defer func(dst *os.File) {
		err = dst.Close()
		if err != nil {
			logger.Errorf("dst close error:%v", err.Error())
		}
	}(dst)

	// 将文件数据写入目标文件
	_, err = dst.Write(fileData)
	if err != nil {
		logger.Error("Error writing the file", fasthttp.StatusInternalServerError)
		return "", err
	}

	// 构建文件URI
	fileURI := fmt.Sprintf("%s%s", config.Domain, fileName)

	return fileURI, nil
}

func generateFileName(originalName string) string {
	randomNumber := rand.Intn(100000)
	currentTime := time.Now().Format("20060102150405")

	fileExt := filepath.Ext(originalName)
	fileName := fmt.Sprintf("%s%d%s", currentTime, randomNumber, fileExt)

	return fileName
}

func (*LocalUploadStrategy) GetFile(ctx jet.Ctx, path string, filePath string) {
	var (
		logger = ctx.Logger()
	)
	if path == "" {
		logger.Errorf("path is empty, url:%v", path)
		return
	}
	fileName := filepath.Base(path)
	fullFilePath := filepath.Join(filePath, fileName)
	// 检查文件是否存在
	_, err := os.Stat(fullFilePath)
	if err != nil {
		logger.Errorf("File not found, %v, %v", fileName, fullFilePath)
		xjet.Error(ctx, "File not found", fasthttp.StatusNotFound)
		return
	}

	// 读取文件内容
	fileData, err := os.ReadFile(fullFilePath)
	if err != nil {
		logger.Error("Error reading the file", fasthttp.StatusInternalServerError)
		xjet.Error(ctx, "Error reading the file", fasthttp.StatusInternalServerError)
		return
	}

	var mineType = mime.TypeByExtension(filepath.Ext(fullFilePath))

	ctx.Response().Header.Set("Content-Length", strconv.Itoa(len(fileData)))
	// 根据文件类型设置响应头
	if strings.HasPrefix(mime.TypeByExtension(filepath.Ext(fullFilePath)), "image/") {
		ctx.Response().Header.Set("Content-Disposition", "inline")
		ctx.Response().Header.Set("Content-Type", mineType)
	} else {
		ctx.Response().Header.Set("Content-Disposition", "attachment; filename="+fileName)
		ctx.Response().Header.Set("Content-Type", "application/octet-stream")
	}
	ctx.Response().SetBody(fileData)
}
