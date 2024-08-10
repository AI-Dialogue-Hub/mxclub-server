package repo

import (
	"context"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"gorm.io/gorm"
	"mxclub/domain/order/entity/dto"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/pkg/common/xmysql"
	"strings"
	"time"
)

func init() {
	jet.Provide(NewWithdrawalRepo)
}

type IWithdrawalRepo interface {
	xmysql.IBaseRepo[po.WithdrawalRecord]
	// WithdrawnAmountNotReject 用户历史提现金额，包括通过和进行中的
	WithdrawnAmountNotReject(ctx jet.Ctx, dasherId int) (float64, error)
	ApproveWithdrawnAmount(ctx jet.Ctx, dasherId int) (float64, error)
	Withdrawn(ctx jet.Ctx, dasherId int, userId uint, dasherName string, amount float64) error
	ListWithdraw(ctx jet.Ctx, d *dto.WithdrawListDTO) ([]*po.WithdrawalRecord, error)
	// ApproveWithdrawnAmountByDasherIds 打手们运行提现的钱
	// @return 打手id -> 可以提现的钱
	ApproveWithdrawnAmountByDasherIds(ctx jet.Ctx, dasherIds []int) (map[int]float64, error)
}

func NewWithdrawalRepo(db *gorm.DB) IWithdrawalRepo {
	repo := new(WithdrawalRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.WithdrawalRecord)
	repo.Ctx = context.Background()
	return repo
}

type WithdrawalRepo struct {
	xmysql.BaseRepo[po.WithdrawalRecord]
}

// ====================================================

func (repo WithdrawalRepo) WithdrawnAmountNotReject(ctx jet.Ctx, dasherId int) (float64, error) {
	var amount float64

	sql := `select COALESCE(sum(withdrawal_amount), 0) 
			from withdrawal_records 
			where dasher_id = ? and withdrawal_status != ?`

	if err := repo.DB().Raw(sql, dasherId, enum.Reject()).Scan(&amount).Error; err != nil {
		ctx.Logger().Errorf("[WithdrawnAmountNotReject]ERROR:%v", err.Error())
		return 0, err
	}

	return amount, nil
}

func (repo WithdrawalRepo) ApproveWithdrawnAmount(ctx jet.Ctx, dasherId int) (float64, error) {
	var amount float64

	sql := `select COALESCE(sum(withdrawal_amount), 0) 
			from withdrawal_records 
			where dasher_id = ? and withdrawal_status = ?`

	if err := repo.DB().Raw(sql, dasherId, enum.Completed()).Scan(&amount).Error; err != nil {
		ctx.Logger().Errorf("[WithdrawnAmountNotReject]ERROR:%v", err.Error())
		return 0, err
	}

	return amount, nil
}

func (repo WithdrawalRepo) ApproveWithdrawnAmountByDasherIds(ctx jet.Ctx, dasherIds []int) (map[int]float64, error) {
	// 初始化结果map
	results := make(map[int]float64)

	// 将dasherIds数组转换为逗号分隔的字符串
	ids := make([]string, len(dasherIds))
	for i, id := range dasherIds {
		ids[i] = fmt.Sprintf("%d", id)
	}
	idsStr := strings.Join(ids, ",")

	// 编写SQL查询，使用IN子句
	sql := fmt.Sprintf(`SELECT dasher_id, COALESCE(SUM(withdrawal_amount), 0) AS amount 
                        FROM withdrawal_records 
                        WHERE dasher_id IN (%s) AND withdrawal_status = ? 
                        GROUP BY dasher_id`, idsStr)

	// 执行查询
	rows, err := repo.DB().Raw(sql, enum.Completed()).Rows()
	if err != nil {
		ctx.Logger().Errorf("[ApproveWithdrawnAmountByDasherIds]ERROR:%v", err.Error())
		return nil, err
	}
	defer rows.Close()

	// 遍历查询结果
	for rows.Next() {
		var dasherId int
		var amount float64
		if err = rows.Scan(&dasherId, &amount); err != nil {
			ctx.Logger().Errorf("[ApproveWithdrawnAmountByDasherIds]ERROR:%v", err.Error())
			return nil, err
		}
		results[dasherId] = amount
	}

	// 检查是否有行扫描错误
	if err = rows.Err(); err != nil {
		ctx.Logger().Errorf("[ApproveWithdrawnAmountByDasherIds]ERROR:%v", err.Error())
		return nil, err
	}

	return results, nil
}

func (repo WithdrawalRepo) Withdrawn(ctx jet.Ctx, dasherId int, userId uint, dasherName string, amount float64) error {
	return repo.InsertOne(&po.WithdrawalRecord{
		DasherID:         dasherId,
		DasherUserId:     userId,
		DasherName:       dasherName,
		WithdrawalAmount: amount,
		WithdrawalStatus: "initiated",
		ApplicationTime:  core.Time(time.Now()),
	})
}

func (repo WithdrawalRepo) ListWithdraw(ctx jet.Ctx, d *dto.WithdrawListDTO) ([]*po.WithdrawalRecord, error) {
	query := xmysql.NewMysqlQuery()
	query.SetPage(d.Page, d.PageSize)
	query.SetFilter("dasher_user_id = ? and created_at >= ? and created_at <= ?", d.UserId, d.Ge, d.Le)
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
