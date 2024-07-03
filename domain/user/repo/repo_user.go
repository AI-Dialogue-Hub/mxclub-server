package repo

import (
	"context"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/user/entity/enum"
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
	AddUserByOpenId(ctx jet.Ctx, openId string) error
	FindByOpenId(ctx jet.Ctx, openId string) (*po.User, error)
	ExistsByOpenId(ctx jet.Ctx, openId string) bool
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

func (u *UserRepo) AddUserByOpenId(ctx jet.Ctx, openId string) error {
	user := &po.User{
		WxOpenId: openId,
		Role:     enum.RoleWxUser,
	}
	err := u.InsertOne(user)
	if err != nil {
		ctx.Logger().Errorf("insert user err:%v", err)
		return err
	}
	id := user.ID
	err = u.DB().Where("id = ?", id).Updates(map[string]interface{}{"wx_name": fmt.Sprintf("用户: %v", id)}).Error
	if err != nil {
		ctx.Logger().Errorf("update user err:%v", err)
		return err
	}
	return nil
}

func (u *UserRepo) FindByOpenId(ctx jet.Ctx, openId string) (*po.User, error) {
	one, err := u.FindOne("wx_open_id = ?", openId)
	if err != nil {
		ctx.Logger().Errorf("find user err:%v", err)
		return nil, err
	}
	return one, err
}

func (u *UserRepo) ExistsByOpenId(ctx jet.Ctx, openId string) bool {
	count, err := u.Count("wx_open_id = ?", openId)
	if err != nil {
		ctx.Logger().Errorf("ExistsByOpenId err:%v", err)
		return false
	}
	return count > 0
}
