package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/apps/mxclub-admin/middleware"
	"mxclub/apps/mxclub-admin/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/utils"
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

func (ctl UserController) GetOverview(ctx jet.Ctx) (*api.Response, error) {
	authInfo, err := middleware.ParseAuthTokenByCtx(ctx)
	if err != nil {
		return xjet.WrapperResult(ctx, nil, err)
	}
	return xjet.WrapperResult(ctx, map[string]any{"username": authInfo.UserName}, nil)
}

type loginParams struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

func (ctl UserController) GetLogin(ctx jet.Ctx, param *loginParams) (*api.Response, error) {
	if err := xjet.IsAnyEmpty(param.Username, param.Password); err != nil {
		return nil, err
	}
	user, err := ctl.userService.CheckUser(ctx, param.Username, param.Password)
	if err != nil {
		return xjet.WrapperResult(ctx, nil, err)
	}
	userVO := &vo.UserVO{
		Name:     user.Name,
		Role:     user.Role,
		JwtToken: middleware.MustGenAuthToken(ctx, user.Name),
	}
	return xjet.WrapperResult(ctx, userVO, nil)
}
