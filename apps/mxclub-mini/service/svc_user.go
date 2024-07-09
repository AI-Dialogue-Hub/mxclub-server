package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	miniUtil "mxclub/apps/mxclub-mini/utils"
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
	user, _ := svc.userRepo.FindByUserId(ctx, openId)
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
	userPO, err := svc.userRepo.FindByUserId(ctx, openId)
	return utils.MustCopyByCtx[vo.User](ctx, userPO), err
}

func (svc UserService) FindUserById(id uint) (*po.User, error) {
	return svc.userRepo.FindByID(id)
}
