package repo

import (
	"context"
	"fmt"
	"github.com/fengyuan-liang/GoKit/collection/maps"
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
	FindByMemberNumber(ctx jet.Ctx, memberNumber int) (*po.User, error)
	FindGradeByUserIdList(userIdList []uint) (maps.IMap[uint, string], error)
	ExistsByOpenId(ctx jet.Ctx, openId string) bool
	ListAroundCache(ctx jet.Ctx, params *api.PageParams) ([]*po.User, int64, error)
	ListAroundCacheByUserType(ctx jet.Ctx, params *api.PageParams, userType enum.RoleType) ([]*po.User, int64, error)
	UpdateUser(ctx jet.Ctx, updateMap map[string]any) error
	// UpdateUserIconAndNickName 如果为空会给默认的头像和昵称
	UpdateUserIconAndNickName(ctx jet.Ctx, id uint, icon, nickName, userInfoJson string) error
	ToBeAssistant(ctx jet.Ctx, userId uint, phone string, memberNumber int64, name string) error
	ExistsAssistant(ctx jet.Ctx, phone string, memberNumber int64) bool
	AssistantOnline(ctx jet.Ctx) ([]*po.User, error)
	CheckAssistantStatus(ctx jet.Ctx, memberNumber int) bool
	UpdateAssistantStatus(ctx jet.Ctx, userId uint, status enum.MemberStatus) error
	UpdateUserPhone(ctx jet.Ctx, id uint, phone string) error
	RemoveDasher(ctx jet.Ctx, id uint) error
	FindByIdAroundCache(ctx jet.Ctx, id uint) (*po.User, error)
}

func NewUserRepo(db *gorm.DB) IUserRepo {
	userRepo := new(UserRepo)
	userRepo.SetDB(db)
	userRepo.ModelPO = new(po.User)
	userRepo.Ctx = context.Background()
	return userRepo
}

type UserRepo struct {
	xmysql.BaseRepo[po.User]
}

// =============================================================================

const userCachePrefix = "mini_user"
const userListCachePrefix = "mini_user_list"

func (repo UserRepo) QueryUserByAccount(username string, password string) (*po.User, error) {
	return repo.FindOne("name = ? and password = ?", username, utils.EncryptPassword(password))
}

func (repo UserRepo) AddUserByOpenId(ctx jet.Ctx, openId string) (uint, error) {
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
	return user.ID, nil
}

func (repo UserRepo) FindByOpenId(ctx jet.Ctx, openId string) (*po.User, error) {
	one, err := repo.FindOne("wx_open_id = ?", openId)
	if err != nil {
		ctx.Logger().Errorf("find user err:%v", err)
		return nil, err
	}
	return one, err
}

func (repo UserRepo) FindByMemberNumber(ctx jet.Ctx, memberNumber int) (*po.User, error) {
	cacheKey := fmt.Sprintf("%v_%v_%v", userCachePrefix, "FindByMemberNumber", memberNumber)
	return xredis.GetOrDefault[po.User](ctx, cacheKey, func() (*po.User, error) {
		return repo.FindOne("member_number = ? and role = ?", memberNumber, enum.RoleAssistant.String())
	})
}

func (repo UserRepo) ExistsByOpenId(ctx jet.Ctx, openId string) bool {
	count, err := repo.Count("wx_open_id = ?", openId)
	if err != nil {
		ctx.Logger().Errorf("ExistsByOpenId err:%v", err)
		return false
	}
	return count > 0
}

func (repo UserRepo) ListAroundCache(ctx jet.Ctx, params *api.PageParams) ([]*po.User, int64, error) {
	return repo.ListAroundCacheByUserType(ctx, params, "")
}

func (repo UserRepo) ListAroundCacheByUserType(ctx jet.Ctx, params *api.PageParams, userType enum.RoleType) ([]*po.User, int64, error) {
	// 根据页码参数生成唯一的缓存键
	cacheListKey := fmt.Sprintf("%s_%s", xredis.BuildListDataCacheKey(userCachePrefix, params), userType)
	cacheCountKey := fmt.Sprintf("%s_%s", xredis.BuildListCountCacheKey(userListCachePrefix), userType)

	list, count, err := xredis.GetListOrDefault[po.User](ctx, cacheListKey, cacheCountKey, func() ([]*po.User, int64, error) {
		// 如果缓存中未找到，则从数据库中获取
		query := xmysql.NewMysqlQuery()
		query.SetPage(params.Page, params.PageSize)
		if userType != "" {
			query.SetFilter("role = ?", userType)
		}
		list, count, err := repo.ListByWrapper(ctx, query)
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

func (repo UserRepo) ToBeAssistant(ctx jet.Ctx, userId uint, phone string, memberNumber int64, name string) error {
	_ = xredis.DelMatchingKeys(ctx, userCachePrefix)
	return repo.UpdateUser(ctx, map[string]any{
		"id":            userId,
		"member_number": memberNumber,
		"phone":         phone,
		"name":          name,
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
	count, _ := repo.Count(
		"member_number = ? and member_status = ? and role = ?",
		memberNumber, enum.Online, enum.RoleAssistant.String(),
	)
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
	cacheKey := fmt.Sprintf("%v_%v", userCachePrefix, userId)
	_ = xredis.Del(cacheKey)
	return repo.UpdateUser(ctx, map[string]any{
		"id":            userId,
		"member_status": status,
	})
}

func (repo UserRepo) UpdateUserIconAndNickName(ctx jet.Ctx, id uint, icon, nickName, userInfoJson string) error {
	_ = xredis.DelMatchingKeys(ctx, userCachePrefix)
	if nickName == "" {
		nickName = fmt.Sprintf("用户%v", 30000+id)
	}
	if icon == "" {
		icon = "https://mx.fengxianhub.top/v1/file/2024071622064557713.jpg"
	}
	return repo.UpdateUser(ctx, map[string]any{
		"id":           id,
		"wx_icon":      icon,
		"wx_name":      nickName,
		"wx_user_info": userInfoJson,
	})
}

func (repo UserRepo) RemoveDasher(ctx jet.Ctx, id uint) error {
	_ = xredis.DelMatchingKeys(ctx, userCachePrefix)
	update := xmysql.NewMysqlUpdate()
	update.SetFilter("id = ?", id)
	update.Set("role", enum.RoleWxUser)
	update.Set("member_number", -1)
	update.Set("phone", 0)
	update.Set("name", "")
	return repo.UpdateByWrapper(update)
}

func (repo UserRepo) FindGradeByUserIdList(userIdList []uint) (maps.IMap[uint, string], error) {
	type Pair struct {
		Id      uint
		WxGrade string `gorm:"column:wx_grade"`
	}
	var result []*Pair
	err := repo.DB().
		Raw(fmt.Sprintf("select id, wx_grade from %s where id in (?)", repo.ModelPO.TableName()), userIdList).
		Scan(&result).
		Error
	if err != nil {
		return nil, err
	}
	var m = maps.NewHashMap[uint, string]()
	for _, pair := range result {
		m.Put(pair.Id, utils.GetOrDefault(pair.WxGrade, "LV0"))
	}
	return m, nil
}

func (repo UserRepo) FindByIdAroundCache(ctx jet.Ctx, id uint) (*po.User, error) {
	cacheKey := fmt.Sprintf("%v_%v", userCachePrefix, id)
	return xredis.GetOrDefault[po.User](ctx, cacheKey, func() (*po.User, error) {
		return repo.FindByID(id)
	})
}
