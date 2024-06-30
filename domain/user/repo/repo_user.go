package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/user/po"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewUserRepo)
}

type IUserRepo interface {
	xmysql.IBaseRepo[po.User]
	QueryUserByAccount(username string, password string) (*po.User, error)
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	userRepo := new(UserRepo)
	userRepo.Db = db.Model(new(po.User))
	userRepo.Ctx = context.Background()
	return userRepo
}

type UserRepo struct {
	xmysql.BaseRepo[po.User]
}

func (u *UserRepo) QueryUserByAccount(username string, password string) (*po.User, error) {
	return u.FindOne("name = ? and password = ?", username, utils.EncryptPassword(password))
}
