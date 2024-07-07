package repo

import (
	"context"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/user/entity/enum"
	"mxclub/domain/user/po"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/common/xredis"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewUserRepo)
}

type IUserRepo interface {
	xmysql.IBaseRepo[po.User]
	QueryUserByAccount(username string, password string) (*po.User, error)
	AddUserByOpenId(ctx jet.Ctx, openId string) (uint, error)
	FindByUserId(ctx jet.Ctx, userId string) (*po.User, error)
	ExistsByOpenId(ctx jet.Ctx, openId string) bool
	ListAroundCache(ctx jet.Ctx, params *api.PageParams) ([]*po.User, int64, error)
	UpdateUser(ctx jet.Ctx, updateMap map[string]any) error
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

func (repo *UserRepo) QueryUserByAccount(username string, password string) (*po.User, error) {
	return repo.FindOne("name = ? and password = ?", username, utils.EncryptPassword(password))
}

func (repo *UserRepo) AddUserByOpenId(ctx jet.Ctx, openId string) (uint, error) {
	_ = xredis.DelMatchingKeys(ctx, userCachePrefix)
	user := &po.User{
		WxOpenId: openId,
		Role:     enum.RoleWxUser,
	}
	err := repo.InsertOne(user)
	if err != nil {
		ctx.Logger().Errorf("insert user err:%v", err)
		return 0, err
	}
	id := user.ID
	err = repo.DB().Where("id = ?", id).Updates(map[string]interface{}{"wx_name": fmt.Sprintf("用户: %v", id)}).Error
	if err != nil {
		ctx.Logger().Errorf("update user err:%v", err)
		return 0, err
	}
	return id, nil
}

func (repo *UserRepo) FindByUserId(ctx jet.Ctx, openId string) (*po.User, error) {
	one, err := repo.FindOne("wx_open_id = ?", openId)
	if err != nil {
		ctx.Logger().Errorf("find user err:%v", err)
		return nil, err
	}
	return one, err
}

func (repo *UserRepo) ExistsByOpenId(ctx jet.Ctx, openId string) bool {
	count, err := repo.Count("wx_open_id = ?", openId)
	if err != nil {
		ctx.Logger().Errorf("ExistsByOpenId err:%v", err)
		return false
	}
	return count > 0
}

const userCachePrefix = "mini_user"
const userListCachePrefix = "mini_user_list"

func (repo *UserRepo) ListAroundCache(ctx jet.Ctx, params *api.PageParams) ([]*po.User, int64, error) {
	// 根据页码参数生成唯一的缓存键
	cacheListKey := xredis.BuildListDataCacheKey(userCachePrefix, params)
	cacheCountKey := xredis.BuildListCountCacheKey(userListCachePrefix)

	list, count, err := xredis.GetListOrDefault[po.User](ctx, cacheListKey, cacheCountKey, func() ([]*po.User, int64, error) {
		// 如果缓存中未找到，则从数据库中获取
		list, count, err := repo.List(params.Page, params.PageSize, nil)
		if err != nil {
			ctx.Logger().Errorf("[repo list]error:%v", err.Error())
			return nil, 0, err
		}
		return list, count, nil
	})

	if err != nil {
		ctx.Logger().Errorf("[UserRepo]ListAroundCache ERROR: %v", err)
		return nil, 0, err
	}

	return list, count, nil
}

func (repo *UserRepo) UpdateUser(ctx jet.Ctx, updateMap map[string]any) error {
	_ = xredis.DelMatchingKeys(ctx, userCachePrefix)
	id := updateMap["id"]
	delete(updateMap, "id")
	err := repo.Update(updateMap, "id = ?", id)
	if err != nil {
		ctx.Logger().Errorf("[UserRepo]UpdateUser ERROR:%v", err.Error())
		return err
	}
	return nil
}
