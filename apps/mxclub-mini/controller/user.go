package controller

import (
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/apps/mxclub-mini/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/utils"

	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
)

func init() {
	jet.Provide(NewUserController)
}

type UserController struct {
	jet.BaseJetController
	userService *service.UserService
}

func NewUserController(_userService *service.UserService) jet.ControllerResult {
	return jet.NewJetController(&UserController{
		userService: _userService,
	})
}

func (ctl UserController) GetV1User0(ctx jet.Ctx, args *jet.Args) (*api.Response, error) {
	if len(args.CmdArgs) == 0 {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, "userId is empty")
	}
	userId := args.CmdArgs[0]
	user, err := ctl.userService.GetUserById(ctx, utils.ParseInt(userId))
	return xjet.WrapperResult(ctx, user, err)
}

func (ctl UserController) GetV1UserOrderCount(ctx jet.Ctx, args *jet.Args) (*api.Response, error) {
	if len(args.CmdArgs) == 0 {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, "userId is empty")
	}
	userId := args.CmdArgs[0]
	user, err := ctl.userService.GetUserById(ctx, utils.ParseInt(userId))
	return xjet.WrapperResult(ctx, user, err)
}

type LoginParams struct {
	AuthCode   string `json:"authCode" form:"authCode"`
	ClientType string `json:"clientType" form:"clientType"`
}

func (ctl UserController) PostClientLoginWx(ctx jet.Ctx, param *LoginParams) (*api.Response, error) {
	token, err := ctl.userService.WxLogin(ctx, param.AuthCode)
	return xjet.WrapperResult(ctx, token, err)
}

func (ctl UserController) PostClientLoginMember(ctx jet.Ctx) (*api.Response, error) {
	tokenInfo := ctx.FastHttpCtx().UserValue("tokenInfo").(*middleware.AuthToken)
	userVO, err := ctl.userService.GetUserByOpenId(ctx, tokenInfo.OpenId)
	return xjet.WrapperResult(ctx, userVO, err)
}
