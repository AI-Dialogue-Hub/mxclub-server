package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/user/po"
	"mxclub/pkg/common/xmysql"
)

type IUserRepo interface {
	xmysql.IBaseRepo[po.User]
}

func init() {
	jet.Provide(NewUserRepo)
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
