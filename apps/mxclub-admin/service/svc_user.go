package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/user/po"
	"mxclub/domain/user/repo"
	_ "mxclub/domain/user/repo"
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

func (svc UserService) GetUserById(ctx jet.Ctx, id int) (*vo.UserVO, error) {
	userPO, err := svc.userRepo.FindByID(id)
	return utils.MustCopy[vo.UserVO](ctx, userPO), err
}

func (svc UserService) CheckUser(ctx jet.Ctx, username string, password string) (*po.User, error) {
	userPO, err := svc.userRepo.QueryUserByAccount(username, password)
	if err != nil || userPO == nil || userPO.ID == 0 {
		ctx.Logger().Infof("user %s not exist", username)
		return nil, errors.New("账号或密码错误")
	}
	return userPO, nil
}
