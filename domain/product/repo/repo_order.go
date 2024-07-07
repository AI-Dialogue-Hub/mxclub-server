package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/product/entity/enum"
	"mxclub/domain/product/po"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/common/xredis"
)

func init() {
	jet.Provide(NewOrderRepo)
}

type IOrderRepo interface {
	xmysql.IBaseRepo[po.Order]
	ListByOrderStatus(ctx jet.Ctx, status enum.OrderStatus, params *api.PageParams, ge, le string) ([]*po.Order, error)
	ListAroundCache(ctx jet.Ctx, params *api.PageParams, ge, le string, status enum.OrderStatus) ([]*po.Order, int64, error)
}

func NewOrderRepo(db *gorm.DB) IOrderRepo {
	repo := new(OrderRepo)
	repo.Db = db.Model(new(po.Order))
	repo.Ctx = context.Background()
	return repo
}

type OrderRepo struct {
	xmysql.BaseRepo[po.Order]
}

const cachePrefix = "_order_CachePrefix"
const listCachePrefix = "_order_configListCachePrefix"

func (repo *OrderRepo) ListByOrderStatus(ctx jet.Ctx, status enum.OrderStatus, params *api.PageParams, ge, le string) ([]*po.Order, error) {
	// 根据页码参数生成唯一的缓存键
	//cacheListKey := xredis.BuildListDataCacheKey(cachePrefix + ge, params)
	//cacheCountKey := xredis.BuildListCountCacheKey(listCachePrefix + le)
	//
	//return xredis.GetListOrDefault(ctx, cacheListKey, cacheCountKey, func() ([]*po.Order, int64, error) {
	//	return repo.List(params.Page, params.PageSize, "order_status = ?", status)
	//})
	if status == 0 {
		return repo.ListNoCount(params.Page, params.PageSize, "purchase_date >= ? and purchase_date <= ?", ge, le)
	} else {
		return repo.ListNoCount(params.Page, params.PageSize, "purchase_date >= ? and purchase_date <= ? and order_status = ?", ge, le, status)
	}
}

func (repo *OrderRepo) ListAroundCache(ctx jet.Ctx, params *api.PageParams, ge, le string, status enum.OrderStatus) ([]*po.Order, int64, error) {
	// 根据页码参数生成唯一的缓存键
	cacheListKey := xredis.BuildListDataCacheKey(cachePrefix+ge+le+status.String(), params)
	cacheCountKey := xredis.BuildListCountCacheKey(listCachePrefix + ge + le + status.String())

	list, count, err := xredis.GetListOrDefault[po.Order](ctx, cacheListKey, cacheCountKey, func() (list []*po.Order, count int64, err error) {
		// 如果缓存中未找到，则从数据库中获取
		if status == 0 {
			list, count, err = repo.List(params.Page, params.PageSize, "purchase_date >= ? and purchase_date <= ?", ge, le)
		} else {
			list, count, err = repo.List(params.Page, params.PageSize, "purchase_date >= ? and purchase_date <= ? and order_status = ?", ge, le, status)
		}
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
