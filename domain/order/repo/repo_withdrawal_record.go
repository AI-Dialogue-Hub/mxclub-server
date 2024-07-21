package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewWithdrawalRepo)
}

type IWithdrawalRepo interface {
	xmysql.IBaseRepo[po.WithdrawalRecord]
	// WithdrawnAmount 用户历史提现金额
	WithdrawnAmount(ctx jet.Ctx, dasherId uint) (float64, error)
	Withdrawn(ctx jet.Ctx, dasherId uint, amount float64) error
}

func NewWithdrawalRepo(db *gorm.DB) IWithdrawalRepo {
	repo := new(WithdrawalRepo)
	repo.Db = db
	repo.ModelPO = new(po.WithdrawalRecord)
	repo.Ctx = context.Background()
	return repo
}

type WithdrawalRepo struct {
	xmysql.BaseRepo[po.WithdrawalRecord]
}

// ====================================================

func (repo WithdrawalRepo) WithdrawnAmount(ctx jet.Ctx, dasherId uint) (float64, error) {
	var amount float64

	sql := "select COALESCE(sum(withdrawal_amount), 0) from withdrawal_records where dasher_id = ? and withdrawal_status = ?"

	if err := repo.DB().Raw(sql, dasherId, enum.Completed()).Scan(&amount).Error; err != nil {
		ctx.Logger().Errorf("[WithdrawnAmount]ERROR:%v", err.Error())
		return 0, err
	}

	return amount, nil
}

func (repo WithdrawalRepo) Withdrawn(ctx jet.Ctx, dasherId uint, amount float64) error {
	return repo.InsertOne(&po.WithdrawalRecord{
		DasherID:         dasherId,
		WithdrawalAmount: amount,
		WithdrawalStatus: "initiated",
	})
}
