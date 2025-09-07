package xoss

import (
	"bytes"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"io"
	"mime/multipart"
	"mxclub/pkg/utils"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type OssUploadStrategy struct{}

// PutFromLocalFile filePath 文件保存路径
func (*OssUploadStrategy) PutFromLocalFile(ctx jet.Ctx, filePath string) (string, error) {
	logger := ctx.Logger()
	fileUrlPath := fmt.Sprintf("%v/%v", ossCfg.StoragePath, getFileNameFromPath(filePath))
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

func getFileNameFromPath(filePath string) string {
	_, fileName := path.Split(filePath)
	return fileName
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
	defer src.Close()

	var finalFile io.Reader

	if ossCfg.Image2WebpURI != "" {
		// 启用压缩功能
		// 先压缩图片到WebP格式
		compressedImage, err := compressImageToWebP(src, file.Filename)
		if err != nil {
			logger.Errorf("Failed to compress image: %v", err)
			// 如果压缩失败，回退到使用原始文件
			if seeker, ok := src.(io.Seeker); ok {
				seeker.Seek(0, io.SeekStart)
			}
			compressedImage = src
			file.Filename = strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename)) + ".webp"
		} else {
			defer compressedImage.Close()
			file.Filename = strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename)) + ".webp"
		}

		finalFile = compressedImage
	} else {
		finalFile = src
	}

	filePath := fmt.Sprintf("%s/%s", ossCfg.StoragePath, utils.GenerateFileName(file.Filename))

	request := &oss.AppendObjectRequest{
		Bucket:   oss.Ptr(ossCfg.Bucket),
		Key:      oss.Ptr(filePath),
		Position: oss.Ptr(int64(0)),
		Body:     finalFile,
	}

	result, err := client.AppendObject(todoCtx, request)
	if err != nil {
		logger.Errorf("AppendObject ERROR:%v", err)
		return "", err
	}

	logger.Infof("append object result:%#v\n", result)

	return BuildURL(filePath), nil
}

var (
	compressImageLogger = xlog.NewWith("compressImageLogger")
)

// compressImageToWebP 调用WebP转换服务压缩图片
func compressImageToWebP(src multipart.File, filename string) (io.ReadCloser, error) {
	defer utils.TraceElapsedWithPrefix(compressImageLogger, "compressImageLogger")()
	// 重置文件读取位置
	if seeker, ok := src.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	// 读取文件内容到内存
	fileData, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 创建multipart表单数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建文件字段
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	// 写入文件数据
	if _, err := part.Write(fileData); err != nil {
		return nil, fmt.Errorf("failed to write file data: %w", err)
	}

	// 关闭writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", ossCfg.Image2WebpURI, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")

	// 设置超时
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to WebP service: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("WebP service returned non-200 status: %d", resp.StatusCode)
	}

	// 返回响应体（WebP图片数据）
	return resp.Body, nil
}

func (o *OssUploadStrategy) GetFile(ctx jet.Ctx, path string, filePath string) {
	// noting to do
}
