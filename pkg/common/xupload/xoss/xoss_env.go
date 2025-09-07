package xoss

import (
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"time"
)

var (
	client  *oss.Client
	ossCfg  *Config
	todoCtx = context.Background()
)

func NewClient(config *Config) *oss.Client {
	ossCfg = config
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.Ak, config.SK)).
		// 设置HTTP连接超时时间为20秒
		WithConnectTimeout(20 * time.Second).
		WithReadWriteTimeout(60 * time.Second).
		// 不校验SSL证书校验
		WithInsecureSkipVerify(true).
		// 设置日志
		WithLogLevel(oss.LogInfo).
		WithRegion(config.Region)

	client = oss.NewClient(cfg)
	if ossCfg.Image2WebpURI != "" {
		xlog.Infof("use compress function, image2webp uri is:%v", ossCfg.Image2WebpURI)
	}
	return client
}

func BuildURL(filePath string) string {
	return ossCfg.Domain + "/" + filePath
}
