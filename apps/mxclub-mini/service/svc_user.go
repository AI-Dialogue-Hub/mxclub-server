package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	miniUtil "mxclub/apps/mxclub-mini/utils"
	"mxclub/domain/user/entity/enum"
	"mxclub/domain/user/po"
	"mxclub/domain/user/repo"
	_ "mxclub/domain/user/repo"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewUserService)
}

type UserService struct {
	userRepo repo.IUserRepo
}

func NewUserService(repo repo.IUserRepo) *UserService {
	return &UserService{userRepo: repo}
}

func (svc UserService) GetUserById(ctx jet.Ctx, id uint) (*vo.User, error) {
	userPO, err := svc.userRepo.FindByID(id)
	return utils.MustCopyByCtx[vo.User](ctx, userPO), err
}

func (svc UserService) WxLogin(ctx jet.Ctx, code string) (string, error) {
	openId, err := miniUtil.GetWxOpenId(code)
	if err != nil || xjet.IsAnyEmpty(openId) {
		ctx.Logger().Errorf("get user id err:%v", err)
		return "", err
	}
	var jwtToken string
	user, _ := svc.userRepo.FindByOpenId(ctx, openId)
	if xjet.IsNil(user) || user.ID <= 0 {
		id, err := svc.userRepo.AddUserByOpenId(ctx, openId)
		if err != nil {
			ctx.Logger().Errorf("get user id err:%v", err)
			return "", errors.New("登录失败")
		}
		jwtToken, _ = middleware.GenAuthTokenByOpenIdAndUserId(id)

	}
	jwtToken, _ = middleware.GenAuthTokenByOpenIdAndUserId(user.ID)
	return jwtToken, err
}

func (svc UserService) GetUserByOpenId(ctx jet.Ctx, openId string) (*vo.User, error) {
	userPO, err := svc.userRepo.FindByOpenId(ctx, openId)
	return utils.MustCopyByCtx[vo.User](ctx, userPO), err
}

func (svc UserService) FindUserById(id uint) (*po.User, error) {
	return svc.userRepo.FindByID(id)
}

func (svc UserService) ToBeAssistant(ctx jet.Ctx, req req.AssistantReq) error {
	if svc.userRepo.ExistsAssistant(ctx, req.Phone, req.MemberNumber) {
		return errors.New("电话或id已被使用")
	}
	err := svc.userRepo.ToBeAssistant(ctx, middleware.MustGetUserId(ctx), req.Phone, req.MemberNumber)
	if err != nil {
		ctx.Logger().Errorf("[ToBeAssistant]ERROR:%v", err.Error())
		return errors.New("转换身份失败，请联系客服")
	}
	return nil
}

func (svc UserService) AssistantOnline(ctx jet.Ctx) []*vo.AssistantOnlineVO {
	userPO, err := svc.userRepo.AssistantOnline(ctx)
	if err != nil {
		ctx.Logger().Errorf("[AssistantOnline]ERROR:%v", err.Error())
		return nil
	}
	return utils.Map[*po.User, *vo.AssistantOnlineVO](userPO, func(in *po.User) *vo.AssistantOnlineVO {
		return &vo.AssistantOnlineVO{
			Id:   in.MemberNumber,
			Name: in.Name,
		}
	})
}

func (svc UserService) CheckAssistantStatus(ctx jet.Ctx, memberNumber int) bool {
	return svc.userRepo.CheckAssistantStatus(ctx, memberNumber)
}

func (svc UserService) SwitchAssistantStatus(ctx jet.Ctx, status enum.MemberStatus) error {
	err := svc.userRepo.UpdateAssistantStatus(ctx, middleware.MustGetUserId(ctx), status)
	if err != nil {
		ctx.Logger().Errorf("[SwitchAssistantStatus]ERROR:%v", err.Error())
		return errors.New("修改在线状态失败，请联系客服")
	}
	return nil
}

func (svc UserService) AssistantStatus(ctx jet.Ctx) string {
	userPO, _ := svc.userRepo.FindByID(middleware.MustGetUserId(ctx))
	return string(userPO.MemberStatus)
}
