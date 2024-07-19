package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"gorm.io/gorm"
	"mxclub/domain/order/entity/dto"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/common/xredis"
	"time"
)

func init() {
	jet.Provide(NewOrderRepo)
}

type IOrderRepo interface {
	xmysql.IBaseRepo[po.Order]
	ListByOrderStatus(
		ctx jet.Ctx,
		status enum.OrderStatus,
		params *api.PageParams,
		ge, le string,
		memberNumber int,
		userId uint) ([]*po.Order, error)
	ListAroundCache(ctx jet.Ctx, params *api.PageParams, ge, le string, status enum.OrderStatus) ([]*po.Order, int64, error)
	// OrderWithdrawAbleAmount 查询打手获得的总金额
	OrderWithdrawAbleAmount(ctx jet.Ctx, dasherId int) (float64, error)
	TotalSpent(ctx jet.Ctx, userId uint) (float64, error)
	FinishOrder(ctx jet.Ctx, id uint, images []string) error
	QueryOrderByStatus(ctx jet.Ctx, processing enum.OrderStatus) ([]*po.Order, error)
	UpdateOrderStatus(ctx jet.Ctx, orderId uint, status enum.OrderStatus) error
	RemoveAssistant(ctx jet.Ctx, executorDTO *dto.OrderExecutorDTO) error
	AddAssistant(ctx jet.Ctx, executorDTO *dto.OrderExecutorDTO) error
}

func NewOrderRepo(db *gorm.DB) IOrderRepo {
	repo := new(OrderRepo)
	repo.Db = db
	repo.ModelPO = new(po.Order)
	repo.Ctx = context.Background()
	return repo
}

type OrderRepo struct {
	xmysql.BaseRepo[po.Order]
}

const cachePrefix = "_order_CachePrefix"
const listCachePrefix = "_order_configListCachePrefix"

func (repo OrderRepo) ListByOrderStatus(ctx jet.Ctx, status enum.OrderStatus, params *api.PageParams, ge, le string, memberNumber int, userId uint) ([]*po.Order, error) {
	// 根据页码参数生成唯一的缓存键
	//cacheListKey := xredis.BuildListDataCacheKey(cachePrefix + ge, params)
	//cacheCountKey := xredis.BuildListCountCacheKey(listCachePrefix + le)
	//
	//return xredis.GetListOrDefault(ctx, cacheListKey, cacheCountKey, func() ([]*po.Order, int64, error) {
	//	return repo.List(params.Page, params.PageSize, "order_status = ?", status)
	//})
	query := new(xmysql.MysqlQuery)
	if memberNumber > 0 {
		query.SetFilter("executor_id = ?", memberNumber)
	} else {
		query.SetFilter("purchase_id = ?", userId)
	}
	query.SetPage(int32(params.Page), int32(params.PageSize))
	query.SetFilter("purchase_date >= ?", ge)
	query.SetFilter("purchase_date <= ?", le)
	if status != 0 {
		query.SetFilter("order_status = ?", status)
	}
	return repo.ListNoCountByQuery(query)
}

func (repo OrderRepo) ListAroundCache(ctx jet.Ctx, params *api.PageParams, ge, le string, status enum.OrderStatus) ([]*po.Order, int64, error) {
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

func (repo OrderRepo) OrderWithdrawAbleAmount(ctx jet.Ctx, dasherId int) (float64, error) {
	var totalAmount float64

	sql := "select COALESCE(sum(executor_price), 0) from orders where executor_id = ? and order_status = ?"

	if err := repo.DB().Raw(sql, dasherId, enum.SUCCESS).Scan(&totalAmount).Error; err != nil {
		ctx.Logger().Errorf("[OrderWithdrawAbleAmount]ERROR:%v", err.Error())
		return 0, err
	}
	return totalAmount, nil
}

func (repo OrderRepo) TotalSpent(ctx jet.Ctx, userId uint) (float64, error) {
	sql := `SELECT SUM(final_price) AS total_price
		 FROM (
		 		 SELECT DISTINCT order_id, final_price
		 		 FROM orders
		 		 WHERE purchase_id = ? AND order_status = ?
		 	 ) AS unique_orders
		 GROUP BY order_id;`

	var totalAmount float64

	if err := repo.DB().Raw(sql, userId, enum.SUCCESS).Scan(&totalAmount).Error; err != nil {
		ctx.Logger().Errorf("[OrderWithdrawAbleAmount]ERROR:%v", err.Error())
		return 0, err
	}
	return totalAmount, nil
}

func (repo OrderRepo) FinishOrder(ctx jet.Ctx, orderId uint, images []string) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	updateMap := map[string]any{
		"detail_images":   xmysql.JSON(images),
		"completion_date": core.Time(time.Now()),
		"order_status":    enum.SUCCESS,
	}
	return repo.Update(updateMap, "order_id = ?", orderId)
}

func (repo OrderRepo) QueryOrderByStatus(ctx jet.Ctx, status enum.OrderStatus) ([]*po.Order, error) {
	return repo.Find("order_status = ? and specify_executor = ?", status, false)
}

func (repo OrderRepo) UpdateOrderStatus(ctx jet.Ctx, orderId uint, status enum.OrderStatus) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	updateMap := map[string]any{
		"order_status": status,
	}
	return repo.Update(updateMap, "order_id = ?", orderId)
}

func (repo OrderRepo) RemoveAssistant(ctx jet.Ctx, executorDTO *dto.OrderExecutorDTO) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	updateWrap := xmysql.NewMysqlUpdate()
	updateWrap.SetFilter("order_id = ?", executorDTO.OrderId)
	executorType := executorDTO.ExecutorType
	if executorType <= 0 {
		return errors.New("executorType cannot empty")
	}
	updateWrap.Set(fmt.Sprintf("executor%v_id", executorType), 0)
	updateWrap.Set(fmt.Sprintf("executor%v_name", executorType), "")
	return repo.UpdateByWrapper(updateWrap)
}

func (repo OrderRepo) AddAssistant(ctx jet.Ctx, executorDTO *dto.OrderExecutorDTO) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	updateWrap := xmysql.NewMysqlUpdate()
	updateWrap.SetFilter("order_id = ?", executorDTO.OrderId)
	executorType := executorDTO.ExecutorType
	if executorType <= 0 {
		executorType = 2
	}
	updateWrap.Set(fmt.Sprintf("executor%v_id", executorType), executorDTO.ExecutorId)
	updateWrap.Set(fmt.Sprintf("executor%v_name", executorType), executorDTO.ExecutorName)
	return repo.UpdateByWrapper(updateWrap)
}
