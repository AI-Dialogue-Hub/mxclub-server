package wxnotify

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"mxclub/pkg/utils"
	"time"
)

type INotifyTokenService interface {
	// FetchToken 获取发送推送消息的token
	FetchToken() (string, error)
}

type notifyTokenService struct {
	grantType string // 消息推送使用默认值client_credential
	appid     string
	secret    string
	// ================
	fetchTokenURI string
	accessToken   string     // 缓存的token
	expiresIn     int        // 过期时间
	lastFetchTime *time.Time // 最后一次请求时间
}

func NewNotifyTokenService(appid, secret string) INotifyTokenService {
	return &notifyTokenService{
		grantType:     "client_credential",
		appid:         appid,
		secret:        secret,
		fetchTokenURI: fmt.Sprintf(TOKEN_URI, appid, secret),
	}
}

var (
	// TOKEN_URI token获取地址
	TOKEN_URI = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v"
	logger    = xlog.NewWith("notifyTokenService")
)

func (svc *notifyTokenService) SetExpiresIn(expiresIn int) {
	svc.expiresIn = expiresIn
	svc.lastFetchTime = utils.CaseToPoint(time.Now())
}

func (svc *notifyTokenService) FetchToken() (string, error) {
	// 1. 判断是否过期，优先使用没过期的token
	if svc.lastFetchTime != nil && time.Now().Sub(*svc.lastFetchTime).Seconds() < float64(svc.expiresIn-100) {
		return svc.accessToken, nil
	}
	type TokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	// 2. token过期 重新请求
	tokenResponse, err := utils.Get[TokenResp](svc.fetchTokenURI)
	if err != nil || tokenResponse == nil {
		logger.Errorf("fetch token error, %v", err)
		return "", errors.New("fetch token error")
	}
	logger.Infof("fetch token response => %v", utils.ObjToJsonStr(tokenResponse))
	if tokenResponse.AccessToken != "" {
		if tokenResponse.ExpiresIn > 0 {
			svc.SetExpiresIn(tokenResponse.ExpiresIn)
		}
		svc.accessToken = tokenResponse.AccessToken
		return svc.accessToken, nil
	}
	return "", nil
}
