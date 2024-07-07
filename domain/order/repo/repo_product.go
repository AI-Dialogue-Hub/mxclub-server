package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/order/po"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/common/xredis"
)

func init() {
	jet.Provide(NewProductRepo)
}

type IProductRepo interface {
	xmysql.IBaseRepo[po.Product]
	ListAroundCache(ctx jet.Ctx, params *api.PageParams) ([]*po.Product, int64, error)
}

func NewProductRepo(db *gorm.DB) IProductRepo {
	repo := new(ProductRepo)
	repo.Db = db.Model(new(po.Product))
	repo.Ctx = context.Background()
	return repo
}

type ProductRepo struct {
	xmysql.BaseRepo[po.Product]
}

// ===========================================================

const productCachePrefix = "mini_product"
const productListCachePrefix = productCachePrefix + "_list"

func (repo ProductRepo) ListAroundCache(ctx jet.Ctx, params *api.PageParams) ([]*po.Product, int64, error) {
	// 根据页码参数生成唯一的缓存键
	cacheListKey := xredis.BuildListDataCacheKey(productCachePrefix, params)
	cacheCountKey := xredis.BuildListCountCacheKey(productListCachePrefix)

	list, count, err := xredis.GetListOrDefault[po.Product](ctx, cacheListKey, cacheCountKey, func() ([]*po.Product, int64, error) {
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
