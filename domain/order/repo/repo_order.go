package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewOrderRepo)
}

type IOrderRepo interface {
	xmysql.IBaseRepo[po.Order]
	ListByOrderStatus(ctx jet.Ctx, status enum.OrderStatus, params *api.PageParams, ge, le string) ([]*po.Order, error)
	ListAroundCache(ctx jet.Ctx, params *api.PageParams, ge, le string, status enum.OrderStatus) ([]*po.Order, int64, error)
	// OrderWithdrawAbleAmount 查询打手获得的总金额
	OrderWithdrawAbleAmount(ctx jet.Ctx, dasherId int) (float64, error)
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
