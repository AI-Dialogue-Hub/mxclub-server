package wxwork

import (
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"mxclub/pkg/utils"
	"sync"
	"time"
)

type WxworkService struct {
	config         *WxWorkConfig
	accessTokenURL string
	accessToken    string
	expiresIn      int // 凭证的有效时间（秒）
	LastModifyTime time.Time
	logger         *xlog.Logger
	rwMutex        *sync.RWMutex
}

func NewWxworkService(config *WxWorkConfig) *WxworkService {
	svc := &WxworkService{config: config}
	svc.accessTokenURL = fmt.Sprintf(
		"https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%v&corpsecret=%v",
		svc.config.Corpid, svc.config.Corpsecret,
	)
	svc.logger = xlog.NewWith("wxworkService")
	svc.rwMutex = new(sync.RWMutex)
	return svc
}

func (svc *WxworkService) GetAccessToken() (string, error) {
	defer utils.RecoverByPrefixNoCtx("WxworkService")

	// 第一次检查缓存
	if token, ok := svc.getTokenByCache(); ok {
		return token, nil
	}

	// 加写锁
	svc.rwMutex.Lock()
	defer svc.rwMutex.Unlock()

	// 第二次检查缓存
	if token, ok := svc.getTokenByCache(); ok {
		return token, nil
	}

	// 请求新 Token
	got, err := utils.Get[AccessTokenResponse](svc.accessTokenURL)
	if err != nil || got == nil {
		svc.logger.Errorf("请求Token失败: err=%v", err)
		return "", fmt.Errorf("请求Token失败: %v", err)
	}
	if got.ErrCode != 0 {
		svc.logger.Errorf("API错误: errcode=%d, errmsg=%s", got.ErrCode, got.ErrMsg)
		return "", fmt.Errorf("API错误: %d-%s", got.ErrCode, got.ErrMsg)
	}

	// 更新缓存
	svc.accessToken = got.AccessToken
	svc.expiresIn = got.ExpiresIn
	svc.LastModifyTime = time.Now()
	svc.logger.Infof("Token刷新成功, expires_in=%d", got.ExpiresIn)
	return svc.accessToken, nil
}

func (svc *WxworkService) getTokenByCache() (token string, ok bool) {
	svc.rwMutex.RLock()
	defer svc.rwMutex.RUnlock()

	if svc.accessToken != "" && time.Since(svc.LastModifyTime) < time.Duration(svc.expiresIn-100)*time.Second {
		return svc.accessToken, true
	}
	return "", false
}
