package repo

import (
	"context"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/order/entity/dto"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
	"time"
)

func init() {
	jet.Provide(NewDeductionRepo)
}

type IDeductionRepo interface {
	xmysql.IBaseRepo[po.Deduction]
	TotalDeduct(ctx jet.Ctx, userId uint) (float64, error)
	ListDeduction(ctx jet.Ctx, d *dto.DeductionDTO) ([]*po.Deduction, int64, error)
	FindDeDuctListBeyondDuration(duration time.Duration) ([]*po.Deduction, error)
	// FindDeDuctListWithDurations 获取指定时间段超时的处罚记录
	//
	// Time Axis: |---- righted (e.g., -5 days ago) ----|---- left (e.g., -3 days ago) ----|---- now ----|
	//
	// Query: SELECT * FROM orders WHERE created_at BETWEEN (now - right) AND (now - left)
	FindDeDuctListWithDurations(left, right time.Duration) ([]*po.Deduction, error)
	// UpdateStatusByIds 批量更新处罚记录
	UpdateStatusByIds(ids []uint, status enum.DeductStatus) error
	FindByOrderNo(orderNo uint) (*po.Deduction, error)
}

func NewDeductionRepo(db *gorm.DB) IDeductionRepo {
	repo := new(DeductionRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.Deduction)
	repo.Ctx = context.Background()
	return repo
}

type DeductionRepo struct {
	xmysql.BaseRepo[po.Deduction]
}

func (repo DeductionRepo) TotalDeduct(ctx jet.Ctx, userId uint) (float64, error) {
	var result float64
	sql := fmt.Sprintf(
		"select COALESCE(sum(amount), 0) from %s where user_id = ? and status = ?",
		repo.ModelPO.TableName(),
	)
	err := repo.DB().Raw(sql, userId, enum.Deduct_SUCCESS).Scan(&result).Error
	if err != nil {
		ctx.Logger().Errorf("[DeductionRepo]TotalDeduct ERROR:%v", err)
		return 0, err
	}
	return result, nil
}

func (repo DeductionRepo) ListDeduction(ctx jet.Ctx, d *dto.DeductionDTO) ([]*po.Deduction, int64, error) {
	query := xmysql.NewMysqlQuery()
	query.SetPage(d.Page, d.PageSize)
	if d.UserId > 0 {
		query.SetFilter("user_id = ?", d.UserId)
	}
	if !utils.IsAnyBlank(d.Le, d.Ge) {
		query.SetFilter("created_at >= ? and created_at <= ?", d.Ge, d.Le)
	}
	query.SetSort("id desc")
	if d.Status != nil {
		query.SetFilter("status = ?", d.Status)
	}
	listNoCountByQuery, count, err := repo.ListByWrapper(ctx, query)
	if err != nil {
		ctx.Logger().Errorf("[DeductionRepo]ListDeduction ERROR:%v", err)
		return nil, 0, err
	}
	return listNoCountByQuery, count, nil
}

// FindDeDuctListBeyondDuration 查找超时指定时间 duration 的订单
func (repo DeductionRepo) FindDeDuctListBeyondDuration(duration time.Duration) ([]*po.Deduction, error) {
	query := xmysql.NewMysqlQuery()
	query.SetFilter("created_at <= ? and status = ?", time.Now().Add(-duration), enum.Deduct_PENDING)
	return repo.ListNoCountByQuery(query)
}

func (repo DeductionRepo) FindDeDuctListWithDurations(left, right time.Duration) ([]*po.Deduction, error) {
	var (
		query = xmysql.NewMysqlQuery()
		now   = time.Now()
	)
	query.SetFilter(
		"created_at <= ? and created_at >= ? and status = ?",
		now.Add(-left), now.Add(-right), enum.Deduct_PENDING)
	return repo.ListNoCountByQuery(query)
}

func (repo DeductionRepo) UpdateStatusByIds(ids []uint, status enum.DeductStatus) error {
	update := xmysql.NewMysqlUpdate()
	update.SetFilter("id in (?)", ids)
	update.Set("status", status)
	return repo.UpdateByWrapper(update)
}

func (repo DeductionRepo) FindByOrderNo(orderNo uint) (*po.Deduction, error) {
	return repo.FindOne("order_no = ?", orderNo)
}
