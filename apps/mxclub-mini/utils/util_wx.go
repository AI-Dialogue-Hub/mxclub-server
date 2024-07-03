package utils

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"mxclub/apps/mxclub-mini/config"
	httpUtil "mxclub/pkg/utils"
)

// WxResp GetWxOpenId response like this
type WxResp struct {
	SessionKey string `json:"session_key" form:"session_key"`
	OpenId     string `json:"open_id" form:"open_id"`
}

func GetWxOpenId(authCode string) (string, error) {
	reqUrl := "https://api.weixin.qq.com/sns/jscode2session"
	reqUrl += "?appid=" + config.GetConfig().WxConfig.Ak
	reqUrl += "&secret=" + config.GetConfig().WxConfig.Sk
	reqUrl += "&js_code=" + authCode
	reqUrl += "&grant_type=authorization_code"
	gotMap, err := httpUtil.Get[map[string]any](reqUrl)
	if err != nil {
		xlog.Error("GetWxOpenId err:", err)
		return "", err
	}
	m := *gotMap
	if openId, ok := m["openid"]; ok {
		return openId.(string), nil
	}
	return "", errors.New("openid not found")
}
