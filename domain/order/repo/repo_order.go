package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/common/xredis"
)

func init() {
	jet.Provide(NewOrderRepo)
}

type IOrderRepo interface {
	xmysql.IBaseRepo[po.Order]
	ListByOrderStatus(ctx jet.Ctx, status enum.OrderStatus, params *api.PageParams) ([]*po.Order, int64, error)
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

func (repo *OrderRepo) ListByOrderStatus(ctx jet.Ctx, status enum.OrderStatus, params *api.PageParams) ([]*po.Order, int64, error) {
	// 根据页码参数生成唯一的缓存键
	cacheListKey := xredis.BuildListDataCacheKey(cachePrefix, params)
	cacheCountKey := xredis.BuildListCountCacheKey(listCachePrefix)

	return xredis.GetListOrDefault[po.Order](ctx, cacheListKey, cacheCountKey, func() ([]*po.Order, int64, error) {
		return repo.List(params.Page, params.PageSize, "order_status = ?", status)
	})
}
