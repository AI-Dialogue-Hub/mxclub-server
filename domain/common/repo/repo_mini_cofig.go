package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"gorm.io/gorm"
	"mxclub/domain/common/po"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/common/xredis"
)

func init() {
	jet.Provide(NewIMiniConfigRepo)
}

type IMiniConfigRepo interface {
	xmysql.IBaseRepo[po.MiniConfig]
	FindConfigByName(ctx jet.Ctx, configName string) (*po.MiniConfig, error)
	FindSwiperConfig(ctx jet.Ctx) (*po.MiniConfig, error)
	AddConfig(ctx jet.Ctx, configName string, content []map[string]any) error
	ExistConfig(ctx jet.Ctx, configName string) (bool, error)
	ListAroundCache(ctx jet.Ctx, params *api.PageParams) ([]*po.MiniConfig, int64, error)
	DeleteById(ctx jet.Ctx, id string) error
	UpdateConfigByConfigName(ctx jet.Ctx, configName string, content []map[string]any) error
}

func NewIMiniConfigRepo(db *gorm.DB) IMiniConfigRepo {
	repo := new(MiniConfigRepo)
	repo.Db = db.Model(new(po.MiniConfig))
	repo.Ctx = context.Background()
	return repo
}

type MiniConfigRepo struct {
	xmysql.BaseRepo[po.MiniConfig]
}

const configCachePrefix = "mini_config_"
const configListCachePrefix = "mini_config_list_"

func (repo MiniConfigRepo) FindConfigByName(ctx jet.Ctx, configName string) (*po.MiniConfig, error) {
	got, err := xredis.GetOrDefault[po.MiniConfig](ctx, configCachePrefix+configName, func() (*po.MiniConfig, error) {
		return repo.FindOne("config_name = ?", configName)
	})
	if err != nil {
		ctx.Logger().Errorf("FindConfigByName error: %v", err.Error())
	}
	return got, err
}

func (repo MiniConfigRepo) FindSwiperConfig(ctx jet.Ctx) (*po.MiniConfig, error) {
	return repo.FindConfigByName(ctx, "swiper")
}

func (repo MiniConfigRepo) AddConfig(ctx jet.Ctx, configName string, content []map[string]any) error {
	// 删除分页缓存
	err := xredis.DelMatchingKeys(ctx, configCachePrefix)
	if err != nil {
		xlog.Errorf(err.Error())
	}
	err = repo.InsertOne(&po.MiniConfig{ConfigName: configName, Content: xmysql.JSONArray(content)})
	if err != nil {
		ctx.Logger().Errorf("InsertOne error: %v", err.Error())
		return err
	}
	return nil
}

func (repo MiniConfigRepo) ExistConfig(ctx jet.Ctx, configName string) (bool, error) {
	count, err := repo.Count("config_name = ?", configName)
	if err != nil {
		ctx.Logger().Errorf("Count error: %v", err.Error())
		return false, err
	}
	return count > 0, nil
}

func (repo MiniConfigRepo) ListAroundCache(ctx jet.Ctx, params *api.PageParams) ([]*po.MiniConfig, int64, error) {
	// 根据页码参数生成唯一的缓存键
	cacheListKey := xredis.BuildListDataCacheKey(configCachePrefix, params)
	cacheCountKey := xredis.BuildListCountCacheKey(configListCachePrefix)

	list, count, err := xredis.GetListOrDefault[po.MiniConfig](ctx, cacheListKey, cacheCountKey, func() ([]*po.MiniConfig, int64, error) {
		// 如果缓存中未找到，则从数据库中获取
		list, count, err := repo.List(params.Page, params.PageSize, nil)
		if err != nil {
			return nil, 0, err
		}
		return list, count, nil
	})
	if err != nil {
		ctx.Logger().Errorf("ListAroundCache 错误: %v", err)
		return nil, 0, err
	}

	return list, count, nil
}

func (repo MiniConfigRepo) DeleteById(ctx jet.Ctx, id string) error {
	_ = xredis.DelMatchingKeys(ctx, configCachePrefix)
	return repo.RemoveByID(id)
}

func (repo MiniConfigRepo) UpdateConfigByConfigName(ctx jet.Ctx, configName string, content []map[string]any) error {
	_ = xredis.DelMatchingKeys(ctx, configCachePrefix)
	updateMap := map[string]any{"content": xmysql.JSONArray(content)}
	err := repo.Update(updateMap, "config_name = ?", configName)
	if err != nil {
		ctx.Logger().Errorf("[UpdateConfigByConfigName]error: %v", err.Error())
	}
	return err
}
