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
)

func init() {
	jet.Provide(NewDeductionRepo)
}

type IDeductionRepo interface {
	xmysql.IBaseRepo[po.Deduction]
	TotalDeduct(ctx jet.Ctx, userId uint) (float64, error)
	ListDeduction(ctx jet.Ctx, d *dto.DeductionDTO) ([]*po.Deduction, error)
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
	sql := fmt.Sprintf("select COALESCE(sum(executor_price), 0) from %s where user_id = ? status = %v", repo.ModelPO.TableName(), enum.Deduct_SUCCESS)
	err := repo.DB().Raw(sql, userId).Scan(&result).Error
	if err != nil {
		ctx.Logger().Errorf("[DeductionRepo]TotalDeduct ERROR:%v", err)
		return 0, err
	}
	return result, nil
}

func (repo DeductionRepo) ListDeduction(ctx jet.Ctx, d *dto.DeductionDTO) ([]*po.Deduction, error) {
	query := xmysql.NewMysqlQuery()
	query.SetPage(d.Page, d.PageSize)
	query.SetFilter("created_at >= ? and created_at <= ?", d.Ge, d.Le)
	if d.Status != nil {
		query.SetFilter("status = ?", d.Status)
	}
	listNoCountByQuery, err := repo.ListNoCountByQuery(query)
	if err != nil {
		ctx.Logger().Errorf("[DeductionRepo]ListDeduction ERROR:%v", err)
		return nil, err
	}
	return listNoCountByQuery, nil
}
