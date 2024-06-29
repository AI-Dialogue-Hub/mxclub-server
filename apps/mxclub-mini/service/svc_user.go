package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/domain/user/repo"
	_ "mxclub/domain/user/repo"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewUserService)
}

type UserService struct {
	xjet.JetContext
	userRepo repo.IUserRepo
}

func NewUserService(repo repo.IUserRepo) *UserService {
	return &UserService{userRepo: repo}
}

func (svc UserService) GetUserById(id int) (*vo.User, error) {
	userPO, err := svc.userRepo.FindByID(id)
	return utils.MustCopy[vo.User](svc.Ctx, userPO), err
}
