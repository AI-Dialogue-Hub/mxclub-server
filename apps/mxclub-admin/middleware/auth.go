package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/router"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/constant"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"mxclub/apps/mxclub-admin/config"
	"mxclub/domain/user/po"
	"mxclub/pkg/api"
	"time"
)

const AuthHeaderKey = "Authorization"

func AuthMiddleware(next router.IJetRouter) (router.IJetRouter, error) {
	return jet.JetHandlerFunc(func(ctx *fasthttp.RequestCtx) {
		if err := handleJwtAuth(ctx); err == nil {
			next.ServeHTTP(ctx)
		} else {
			ctx.Response.Header.SetServer("JetServer")
			ctx.Response.Header.Set("Content-Type", constant.MIMETextPlain)
			ctx.SetBodyString(err.Error())
		}
	}), nil
}

var logger = xlog.NewWith("auth_middleware")

func handleJwtAuth(ctx *fasthttp.RequestCtx) (err error) {
	path := string(ctx.Path())
	if config.IsOpenApi(path) {
		return nil
	}
	jwtToken := string(ctx.Request.Header.Peek(AuthHeaderKey))
	if jwtToken == "" {
		logger.Error("empty Authorization")
		logger.Errorf("path:%v", path)
		err = api.ErrorUnauthorized(logger.ReqId)
	}
	tokenInfo, err := ParseAuthToken(jwtToken)
	if err != nil {
		return err
	}
	ctx.SetUserValue("tokenInfo", tokenInfo)
	return
}

type AuthToken struct {
	jwt.StandardClaims
	UserName string
	UserPO   *po.User
}

func MustGenAuthToken(ctx jet.Ctx, userPO *po.User) string {
	token, err := GenAuthTokenByUserName(userPO)
	if err != nil {
		ctx.Logger().Infof("GenAuthTokenByUserName error:%v", err.Error())
		return ""
	}
	return token
}

func GenAuthTokenByUserName(userPO *po.User) (string, error) {
	authToken := &AuthToken{
		UserName: userPO.Name,
		UserPO:   userPO,
	}
	if authToken.ExpiresAt == 0 {
		authToken.ExpiresAt = time.Now().Unix() + 7*86400
	}

	signKey := []byte(config.GetConfig().Server.JwtKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authToken)
	ss, err := token.SignedString(signKey)
	return ss, err
}

func ParseAuthToken(token string) (*AuthToken, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &AuthToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().Server.JwtKey), nil
	})
	if err != nil {
		if isExpiredTokenError(err) {
			return nil, api.Error("", api.ErrorCodeTokenExpired, "token 已经过期")
		} else {
			return nil, api.ErrorUnauthorized("")
		}
	}

	if claims, ok := jwtToken.Claims.(*AuthToken); ok && jwtToken.Valid {
		return claims, nil
	}

	return nil, api.ErrorUnauthorized("")
}

func isExpiredTokenError(err error) bool {
	var jwtErr *jwt.ValidationError
	if errors.As(err, &jwtErr) {
		return jwtErr.Errors&jwt.ValidationErrorExpired != 0
	}
	return false
}

func ParseAuthTokenByCtx(ctx jet.Ctx) (*AuthToken, error) {
	jwtToken := string(ctx.Request().Header.Peek(AuthHeaderKey))
	return ParseAuthToken(jwtToken)
}

func MustGetUserInfo(ctx jet.Ctx) *po.User {
	authInfo, _ := ParseAuthTokenByCtx(ctx)
	return authInfo.UserPO
}
