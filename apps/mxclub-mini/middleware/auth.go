package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/fengyuan-liang/jet-web-fasthttp/core/router"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/constant"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/valyala/fasthttp"
	"mxclub/apps/mxclub-mini/config"
	"mxclub/pkg/api"
	"time"
)

const AuthHeaderKey = "token"

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
	if config.IsOpenApi(string(ctx.Path())) {
		return nil
	}
	jwtToken := string(ctx.Request.Header.Peek(AuthHeaderKey))
	if jwtToken == "" {
		logger.Error("empty Authorization")
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
	UserId uint
}

func MustGenAuthToken(ctx jet.Ctx, userId uint) string {
	token, err := GenAuthTokenByOpenIdAndUserId(userId)
	if err != nil {
		ctx.Logger().Infof("GenAuthTokenByUserName error:%v", err.Error())
		return ""
	}
	return token
}

func GenAuthTokenByOpenIdAndUserId(userId uint) (string, error) {
	authToken := &AuthToken{
		UserId: userId,
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

func MustGetUserId(ctx jet.Ctx) (userId uint) {
	authInfo, _ := ParseAuthTokenByCtx(ctx)
	return authInfo.UserId
}
