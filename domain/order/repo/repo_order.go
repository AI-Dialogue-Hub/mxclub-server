package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	traceUtil "github.com/fengyuan-liang/jet-web-fasthttp/pkg/utils"
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
	ListByOrderStatus(ctx jet.Ctx, d *dto.ListByOrderStatusDTO) ([]*po.Order, error)
	ListAroundCache(ctx jet.Ctx, params *api.PageParams, ge, le string, status enum.OrderStatus) ([]*po.Order, int64, error)
	// OrderWithdrawAbleAmount 查询打手获得的总金额
	OrderWithdrawAbleAmount(ctx jet.Ctx, dasherId int) (float64, error)
	TotalSpent(ctx jet.Ctx, userId uint) (float64, error)
	FinishOrder(ctx jet.Ctx, d *dto.FinishOrderDTO) error
	// FindByOrderId orderId 订单流水号
	FindByOrderId(ctx jet.Ctx, orderId uint) (*po.Order, error)
	QueryOrderByStatus(ctx jet.Ctx, processing enum.OrderStatus) ([]*po.Order, error)
	QueryOrderWithDelayTime(ctx jet.Ctx, status enum.OrderStatus, thresholdTime time.Time) ([]*po.Order, error)
	// UpdateOrderStatus 这里的orderId为订单流水号
	UpdateOrderStatus(ctx jet.Ctx, orderId uint64, status enum.OrderStatus) error
	RemoveAssistant(ctx jet.Ctx, executorDTO *dto.OrderExecutorDTO) error
	AddAssistant(ctx jet.Ctx, executorDTO *dto.OrderExecutorDTO) error
	GrabOrder(ctx jet.Ctx, ordersId uint, executorId int) error
	// UpdateOrderByDasher 通过车头进行更新
	UpdateOrderByDasher(ctx jet.Ctx, ordersId uint, executorId int, status enum.OrderStatus) error
	UpdateOrderDasher2(ctx jet.Ctx, ordersId uint, executor2Id int, executor2Name string) error
	UpdateOrderDasher3(ctx jet.Ctx, ordersId uint, executor3Id int, executor3Name string) error
	DoneEvaluation(id uint) error
	RemoveByTradeNo(orderNo string) error
	FindByDasherId(ctx jet.Ctx, dasherId int) (*po.Order, error)
}

func NewOrderRepo(db *gorm.DB) IOrderRepo {
	repo := new(OrderRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.Order)
	repo.Ctx = context.Background()
	return repo
}

type OrderRepo struct {
	xmysql.BaseRepo[po.Order]
}

const cachePrefix = "_order_CachePrefix"
const listCachePrefix = "_order_configListCachePrefix"

func (repo OrderRepo) ListByOrderStatus(ctx jet.Ctx, d *dto.ListByOrderStatusDTO) ([]*po.Order, error) {
	// 根据页码参数生成唯一的缓存键
	//cacheListKey := xredis.BuildListDataCacheKey(cachePrefix + ge, params)
	//cacheCountKey := xredis.BuildListCountCacheKey(listCachePrefix + le)
	//
	//return xredis.GetListOrDefault(ctx, cacheListKey, cacheCountKey, func() ([]*po.Order, int64, error) {
	//	return repo.List(params.Page, params.PageSize, "order_status = ?", status)
	//})
	query := new(xmysql.MysqlQuery)
	if d.IsDasher {
		// 这里有可能是三个打手中的任意一个
		query.SetFilter(
			"(executor_id = ? or executor2_id = ? or executor3_id = ?)",
			d.MemberNumber, d.MemberNumber, d.MemberNumber,
		)
	} else {
		query.SetFilter("purchase_id = ?", d.UserId)
	}
	query.SetPage(d.PageParams.Page, d.PageParams.PageSize)
	query.SetFilter("purchase_date >= ?", d.Ge)
	query.SetFilter("purchase_date <= ?", d.Le)
	if d.Status != 0 {
		query.SetFilter("order_status = ?", d.Status)
	}
	return repo.ListNoCountByQuery(query)
}

func (repo OrderRepo) ListAroundCache(ctx jet.Ctx, params *api.PageParams, ge, le string, status enum.OrderStatus) ([]*po.Order, int64, error) {
	// 根据页码参数生成唯一的缓存键
	//cacheListKey := xredis.BuildListDataCacheKey(cachePrefix+ge+le+status.String(), params)
	//cacheCountKey := xredis.BuildListCountCacheKey(listCachePrefix + ge + le + status.String())

	// list, count, err := xredis.GetListOrDefault[po.Order](ctx, cacheListKey, cacheCountKey, func() (list []*po.Order, count int64, err error) {
	// 	// 如果缓存中未找到，则从数据库中获取
	// 	if status == 0 {
	// 		list, count, err = repo.List(params.Page, params.PageSize, "purchase_date >= ? and purchase_date <= ?", ge, le)
	// 	} else {
	// 		list, count, err = repo.List(params.Page, params.PageSize, "purchase_date >= ? and purchase_date <= ? and order_status = ?", ge, le, status)
	// 	}
	// 	if err != nil {
	// 		return nil, 0, err
	// 	}
	// 	return list, count, nil
	// })

	query := xmysql.NewMysqlQuery()

	query.SetPage(params.Page, params.PageSize)

	query.SetFilter("purchase_date >= ? and purchase_date <= ?", ge, le)

	if status != 0 {
		query.SetFilter("order_status = ?", status)
	}

	list, count, err := repo.ListByWrapper(ctx, query)

	if err != nil {
		ctx.Logger().Errorf("ListAroundCache 错误: %v", err)
		return nil, 0, err
	}
	return list, count, nil

}

func (repo OrderRepo) OrderWithdrawAbleAmount(ctx jet.Ctx, dasherId int) (result float64, err error) {
	defer traceUtil.TraceElapsedByName(time.Now(), "OrderWithdrawAbleAmount")
	var totalAmount float64
	var amount1, amount2, amount3 float64
	c1 := make(chan struct{})
	c2 := make(chan struct{})
	c3 := make(chan struct{})
	// 查询executor_id匹配的金额
	go func() {
		defer traceUtil.TraceElapsedByName(time.Now(), "c1")
		defer func() { c1 <- struct{}{} }()
		sql1 := `select COALESCE(sum(executor_price), 0) 
			 from orders 
			 where executor_id = ? and order_status = ?`
		if err = repo.DB().Raw(sql1, dasherId, enum.SUCCESS).Scan(&amount1).Error; err != nil {
			ctx.Logger().Errorf("[OrderWithdrawAbleAmount] ERROR in sql1: %v", err.Error())
			result = 0
		}
	}()

	// 查询executor2_id匹配的金额
	go func() {
		defer traceUtil.TraceElapsedByName(time.Now(), "c2")
		defer func() { c2 <- struct{}{} }()
		sql2 := `select COALESCE(sum(executor2_price), 0) 
			 from orders 
			 where executor2_id = ? and order_status = ?`
		if err = repo.DB().Raw(sql2, dasherId, enum.SUCCESS).Scan(&amount2).Error; err != nil {
			ctx.Logger().Errorf("[OrderWithdrawAbleAmount] ERROR in sql2: %v", err.Error())
			result = 0
		}
	}()

	// 查询executor3_id匹配的金额
	go func() {
		defer traceUtil.TraceElapsedByName(time.Now(), "c3")
		defer func() { c3 <- struct{}{} }()
		sql3 := `select COALESCE(sum(executor3_price), 0) 
			 from orders 
			 where executor3_id = ? and order_status = ?`
		if err = repo.DB().Raw(sql3, dasherId, enum.SUCCESS).Scan(&amount3).Error; err != nil {
			ctx.Logger().Errorf("[OrderWithdrawAbleAmount] ERROR in sql3: %v", err.Error())
			result = 0
		}
	}()

	<-c1
	<-c2
	<-c3

	// 总金额
	totalAmount = amount1 + amount2 + amount3

	return totalAmount, nil
}

func (repo OrderRepo) TotalSpent(ctx jet.Ctx, userId uint) (float64, error) {
	sql := `select COALESCE(SUM(final_price), 0) AS total_price
			from orders
			where purchase_id = ? and order_status = ?`

	var totalAmount float64

	if err := repo.DB().Raw(sql, userId, enum.SUCCESS).Scan(&totalAmount).Error; err != nil {
		ctx.Logger().Errorf("[OrderWithdrawAbleAmount]ERROR:%v", err.Error())
		return 0, err
	}
	return totalAmount, nil
}

func (repo OrderRepo) FinishOrder(ctx jet.Ctx, d *dto.FinishOrderDTO) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	update := xmysql.NewMysqlUpdate()
	update.SetFilter("id = ?", d.Id)
	update.Set("detail_images", xmysql.StringArray(d.Images))
	update.Set("completion_date", core.Time(time.Now()))
	update.Set("order_status", enum.SUCCESS)
	update.Set("executor_price", d.ExecutorPrice)
	update.Set("cut_rate", d.CutRate)
	if d.ExecutorNum == 2 {
		update.Set("executor2_price", d.ExecutorPrice)
	} else if d.ExecutorNum == 3 {
		update.Set("executor2_price", d.ExecutorPrice)
		update.Set("executor3_price", d.ExecutorPrice)
	}
	return repo.UpdateByWrapper(update)
}

func (repo OrderRepo) QueryOrderByStatus(ctx jet.Ctx, status enum.OrderStatus) ([]*po.Order, error) {
	return repo.Find("order_status = ? and specify_executor = ?", status, false)
}

func (repo OrderRepo) QueryOrderWithDelayTime(
	ctx jet.Ctx, status enum.OrderStatus, thresholdTime time.Time) ([]*po.Order, error) {
	return repo.Find("order_status = ? and specify_executor = ? and purchase_date < ?", status, false, thresholdTime)
}

func (repo OrderRepo) UpdateOrderStatus(ctx jet.Ctx, orderId uint64, status enum.OrderStatus) error {
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

func (repo OrderRepo) GrabOrder(ctx jet.Ctx, ordersId uint, executorId int) error {
	defer traceUtil.TraceElapsedByName(time.Now(), fmt.Sprintf("[%s]orderRepo GrabOrder", ctx.Logger().ReqId))
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	tx := repo.DB().Begin()
	// 1. 读取一个未被抢单的订单，并锁定该行（读取锁）
	var lockOrderId uint
	row := tx.Raw(
		"SELECT id FROM orders WHERE id = ? and order_status = ? and specify_executor = 0 LIMIT 1 FOR UPDATE",
		ordersId, enum.PROCESSING,
	)
	if err := row.Scan(&lockOrderId).Error; err != nil || lockOrderId <= 0 {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Logger().Errorf("ERROR:%v", err.Error())
			return errors.New("no pending orders available")
		}
		ctx.Logger().Errorf("lockOrderId:%v", lockOrderId)
		return errors.New("no pending orders available")
	}
	// 2. 更新该订单的状态为已抢单，并设置执行者 Id
	err := tx.Exec(
		"UPDATE orders SET order_status = ?, executor_id = ?, specify_executor = ? WHERE id = ?",
		enum.PROCESSING, executorId, true, ordersId,
	).Error
	if err != nil {
		tx.Rollback()
		ctx.Logger().Errorf("ERROR:%v", err.Error())
		return errors.New("update orders failed")
	}
	// 3. 提交事物
	tx.Commit()
	ctx.Logger().Infof("Order %d has been claimed by executor %d\n", ordersId, executorId)
	return nil
}

func (repo OrderRepo) UpdateOrderByDasher(ctx jet.Ctx, ordersId uint, executorId int, status enum.OrderStatus) error {
	update := xmysql.NewMysqlUpdate()
	update.SetFilter("id = ?", ordersId)
	update.Set("executor_id", executorId)
	update.Set("order_status", status)
	return repo.UpdateByWrapper(update)
}
func (repo OrderRepo) UpdateOrderDasher2(ctx jet.Ctx, ordersId uint, executor2Id int, executor2Name string) error {
	update := xmysql.NewMysqlUpdate()
	update.SetFilter("id = ?", ordersId)
	update.Set("executor2_id", executor2Id)
	update.Set("executor2_name", executor2Name)
	return repo.UpdateByWrapper(update)
}
func (repo OrderRepo) UpdateOrderDasher3(ctx jet.Ctx, ordersId uint, executor3Id int, executor3Name string) error {
	update := xmysql.NewMysqlUpdate()
	update.SetFilter("id = ?", ordersId)
	update.Set("executor3_id", executor3Id)
	update.Set("executor3_name", executor3Name)
	return repo.UpdateByWrapper(update)
}

func (repo OrderRepo) FindByOrderId(ctx jet.Ctx, orderId uint) (*po.Order, error) {
	return repo.FindOne("order_id = ?", orderId)
}

func (repo OrderRepo) DoneEvaluation(id uint) error {
	update := xmysql.NewMysqlUpdate()
	update.SetFilter("id = ?", id)
	update.Set("is_evaluation", true)
	return repo.UpdateByWrapper(update)
}

func (repo OrderRepo) RemoveByTradeNo(orderNo string) error {
	return repo.RemoveOne("order_id = ?", orderNo)
}

func (repo OrderRepo) FindByDasherId(ctx jet.Ctx, dasherId int) (*po.Order, error) {
	query := xmysql.NewMysqlQuery()
	// 这里有可能是三个打手中的任意一个
	query.SetFilter(
		"order_status = 2 and (executor_id = ? or executor2_id = ? or executor3_id = ?)",
		dasherId, dasherId, dasherId,
	)
	return repo.FindByWrapper(query)
}
