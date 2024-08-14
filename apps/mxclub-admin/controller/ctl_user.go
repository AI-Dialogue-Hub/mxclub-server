package controller

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/apps/mxclub-admin/middleware"
	"mxclub/apps/mxclub-admin/service"
	"mxclub/domain/user/entity/enum"
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

// =========================================================================

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
	if xjet.IsAnyEmpty(param.Username, param.Password) {
		return xjet.WrapperResult(ctx, nil, errors.New("username or password is empty"))
	}
	user, err := ctl.userService.CheckUser(ctx, param.Username, param.Password)
	if err != nil {
		return xjet.WrapperResult(ctx, nil, err)
	}
	userVO := &vo.UserLoginVO{
		Name:     user.Name,
		Role:     user.Role,
		JwtToken: middleware.MustGenAuthToken(ctx, user),
	}
	return xjet.WrapperResult(ctx, userVO, nil)
}

func (ctl UserController) PostV1UserList(ctx jet.Ctx, params *req.UserListReq) (*api.Response, error) {
	pageResult, err := ctl.userService.List(ctx, params)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (ctl UserController) GetV1UserTypeList(ctx jet.Ctx) (*api.Response, error) {
	return xjet.WrapperResult(ctx, vo.WrapUserTypeVOS(enum.RoleDisPlayNameMap), nil)
}

func (ctl UserController) PostV1UserUpdate(ctx jet.Ctx, userReq *req.UserReq) (*api.Response, error) {
	if err := middleware.MustGetUserInfo(ctx).Role.CheckPermission(enum.PermissionAdminRead); err != nil {
		return xjet.WrapperResult(ctx, nil, errors.New("权限不够"))
	}
	err := ctl.userService.Update(ctx, userReq)
	return xjet.WrapperResult(ctx, "Ok", err)
}

func (ctl UserController) DeleteV1UserDasher0(ctx jet.Ctx, param *api.PathParam) (*api.Response, error) {
	userId, _ := param.GetInt64(0)
	return xjet.WrapperResult(ctx, "Ok", ctl.userService.RemoveAssistant(ctx, uint(userId)))
}

// GetV1AssistantOnline 获取在线打手，不包括打手自己
func (ctl UserController) GetV1AssistantOnline(ctx jet.Ctx) (*api.Response, error) {
	return xjet.WrapperResult(ctx, ctl.userService.AssistantOnline(ctx), nil)
}
