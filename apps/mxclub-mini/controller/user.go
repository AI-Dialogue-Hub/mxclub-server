package controller

import (
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
