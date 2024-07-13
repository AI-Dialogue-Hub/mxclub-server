package controller

import (
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/apps/mxclub-mini/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/captcha"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/utils"

	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
)

func init() {
	jet.Provide(NewUserController)
}

type UserController struct {
	jet.BaseJetController
	userService    *service.UserService
	messageService *service.MessageService
}

func NewUserController(userService *service.UserService, messageService *service.MessageService) jet.ControllerResult {
	return jet.NewJetController(&UserController{
		userService:    userService,
		messageService: messageService,
	})
}

// ============================================================================================================

func (ctl UserController) GetV1User0(ctx jet.Ctx, args *jet.Args) (*api.Response, error) {
	if len(args.CmdArgs) == 0 {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, "userId is empty")
	}
	userId := args.CmdArgs[0]
	user, err := ctl.userService.GetUserById(ctx, utils.ParseUint(userId))
	return xjet.WrapperResult(ctx, user, err)
}

func (ctl UserController) GetV1UserOrderCount(ctx jet.Ctx, args *jet.Args) (*api.Response, error) {
	if len(args.CmdArgs) == 0 {
		return nil, api.ErrorBadRequest(ctx.Logger().ReqId, "userId is empty")
	}
	userId := args.CmdArgs[0]
	user, err := ctl.userService.GetUserById(ctx, utils.ParseUint(userId))
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
	userVO, err := ctl.userService.GetUserById(ctx, tokenInfo.UserId)
	return xjet.WrapperResult(ctx, userVO, err)
}

func (ctl UserController) PostV1MessageList(ctx jet.Ctx, params *api.PageParams) (*api.Response, error) {
	userId := middleware.MustGetUserInfo(ctx)
	pageResponse, err := ctl.messageService.List(ctx, userId, params)
	return xjet.WrapperResult(ctx, pageResponse, err)
}

func (ctl UserController) PostV1Message(ctx jet.Ctx, params *api.PageParams) (*api.Response, error) {
	userId := middleware.MustGetUserInfo(ctx)
	pageResponse, err := ctl.messageService.List(ctx, userId, params)
	return xjet.WrapperResult(ctx, pageResponse, err)
}

func (ctl UserController) PostV1MessageUnreadCount(ctx jet.Ctx) (*api.Response, error) {
	countUnReadMessage, err := ctl.messageService.CountUnReadMessage(ctx)
	return xjet.WrapperResult(ctx, countUnReadMessage, err)
}

func (ctl UserController) PostV1MessageRead(ctx jet.Ctx) (*api.Response, error) {
	err := ctl.messageService.ReadAllMessage(ctx)
	return xjet.WrapperResult(ctx, "ok", err)
}

// ====== 验证码 ========

var stringCaptcha = captcha.NewCaptcha()

func (ctl UserController) GetV1CaptchaGenerate(ctx jet.Ctx) (*api.Response, error) {
	id, base64, answer := stringCaptcha.Generate()
	captchaVO := &vo.CaptchaVO{
		CaptchaId: id,
		B64s:      base64,
		Answer:    answer,
	}
	return xjet.WrapperResult(ctx, captchaVO, nil)
}

func (ctl UserController) PostV1CaptchaVerify(ctx jet.Ctx, req *req.CaptchaReq) (*api.Response, error) {
	result := stringCaptcha.Verify(req.CaptchaId, req.Answer)
	captchaVerifyVO := &vo.CaptchaVerifyVO{
		CaptchaId: req.CaptchaId,
		Answer:    req.Answer,
		Result:    result,
	}
	return xjet.WrapperResult(ctx, captchaVerifyVO, nil)
}

// ========================

func (ctl UserController) PostV1UserAssistant(ctx jet.Ctx, assistantReq req.AssistantReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "ok", ctl.userService.ToBeAssistant(ctx, assistantReq))
}
