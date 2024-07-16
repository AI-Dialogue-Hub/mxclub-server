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
	FindByOpenId(ctx jet.Ctx, userId string) (*po.User, error)
	ExistsByOpenId(ctx jet.Ctx, openId string) bool
	ListAroundCache(ctx jet.Ctx, params *api.PageParams) ([]*po.User, int64, error)
	UpdateUser(ctx jet.Ctx, updateMap map[string]any) error
	ToBeAssistant(ctx jet.Ctx, userId uint, phone string, memberNumber int64) error
	ExistsAssistant(ctx jet.Ctx, phone string, memberNumber int64) bool
	AssistantOnline(ctx jet.Ctx) ([]*po.User, error)
	CheckAssistantStatus(ctx jet.Ctx, memberNumber int) bool
	UpdateAssistantStatus(ctx jet.Ctx, userId uint, status enum.MemberStatus) error
	UpdateUserPhone(ctx jet.Ctx, id uint, phone string) error
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

func (repo UserRepo) QueryUserByAccount(username string, password string) (*po.User, error) {
	return repo.FindOne("name = ? and password = ?", username, utils.EncryptPassword(password))
}

func (repo UserRepo) AddUserByOpenId(ctx jet.Ctx, openId string) (uint, error) {
	_ = xredis.DelMatchingKeys(ctx, userCachePrefix)
	user := &po.User{
		WxOpenId: openId,
		Role:     enum.RoleWxUser,
		WxIcon:   "https://mx.fengxianhub.top/v1/file/2024071621495340739.png",
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

func (repo UserRepo) FindByOpenId(ctx jet.Ctx, openId string) (*po.User, error) {
	one, err := repo.FindOne("wx_open_id = ?", openId)
	if err != nil {
		ctx.Logger().Errorf("find user err:%v", err)
		return nil, err
	}
	return one, err
}

func (repo UserRepo) ExistsByOpenId(ctx jet.Ctx, openId string) bool {
	count, err := repo.Count("wx_open_id = ?", openId)
	if err != nil {
		ctx.Logger().Errorf("ExistsByOpenId err:%v", err)
		return false
	}
	return count > 0
}

const userCachePrefix = "mini_user"
const userListCachePrefix = "mini_user_list"

func (repo UserRepo) ListAroundCache(ctx jet.Ctx, params *api.PageParams) ([]*po.User, int64, error) {
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

func (repo UserRepo) UpdateUser(ctx jet.Ctx, updateMap map[string]any) error {
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

func (repo UserRepo) ToBeAssistant(ctx jet.Ctx, userId uint, phone string, memberNumber int64) error {
	_ = xredis.DelMatchingKeys(ctx, userCachePrefix)
	return repo.UpdateUser(ctx, map[string]any{
		"id":            userId,
		"member_number": memberNumber,
		"phone":         phone,
	})
}

func (repo UserRepo) ExistsAssistant(ctx jet.Ctx, phone string, memberNumber int64) bool {
	count, _ := repo.Count("phone = ? or member_number = ?", phone, memberNumber)
	return count >= 1
}

func (repo UserRepo) AssistantOnline(ctx jet.Ctx) ([]*po.User, error) {
	return repo.Find("member_status = ?", enum.Online)
}

// CheckAssistantStatus 检查打手是否可以接单, 必须是在线状态
func (repo UserRepo) CheckAssistantStatus(ctx jet.Ctx, memberNumber int) bool {
	count, _ := repo.Count("member_number = ? and member_status = ?", memberNumber, enum.Online)
	return count >= 1
}

func (repo UserRepo) UpdateUserPhone(ctx jet.Ctx, userId uint, phone string) error {
	_ = xredis.DelMatchingKeys(ctx, userCachePrefix)
	return repo.UpdateUser(ctx, map[string]any{
		"id":    userId,
		"phone": phone,
	})
}

func (repo UserRepo) UpdateAssistantStatus(ctx jet.Ctx, userId uint, status enum.MemberStatus) error {
	_ = xredis.DelMatchingKeys(ctx, userCachePrefix)
	return repo.UpdateUser(ctx, map[string]any{
		"id":            userId,
		"member_status": status,
	})
}
