package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mxclub/domain/lottery/entity/dto"
	"mxclub/domain/lottery/po"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewLotteryRecordsRepo)
}

type ILotteryRecordsRepo interface {
	xmysql.IBaseRepo[po.LotteryRecords]
	ListRecords(ctx jet.Ctx, params *api.PageParams) ([]*dto.LotteryRecordsDTO, int64, error)
}

func NewLotteryRecordsRepo(db *gorm.DB) ILotteryRecordsRepo {
	repo := new(LotteryRecordsRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.LotteryRecords)
	repo.Ctx = context.Background()
	return repo
}

type LotteryRecordsRepo struct {
	xmysql.BaseRepo[po.LotteryRecords]
}

func (repo *LotteryRecordsRepo) ListRecords(
	ctx jet.Ctx, params *api.PageParams,
) ([]*dto.LotteryRecordsDTO, int64, error) {
	var (
		dataSQL = `
		SELECT lottery_records.id,
		       activity.id       AS activity_id,
			   activity.activity_title,
			   activity.activity_price,
			   lottery_prizes.id AS prize_id,
			   lottery_prizes.prize_name,
			   lottery_records.user_id,
			   lottery_records.activity_prize_snapshot,
			   lottery_records.created_at
		FROM lottery_records
				 LEFT JOIN lottery_activities activity ON activity.id = lottery_records.activity_id
				 LEFT JOIN lottery_prizes ON lottery_records.prize_id = lottery_prizes.id
		ORDER BY lottery_records.created_at DESC
		LIMIT ? OFFSET ?
		;
	`
		countSQL = `
		SELECT COUNT(1)
		FROM lottery_records
				 LEFT JOIN lottery_activities activity ON activity.id = lottery_records.activity_id
				 LEFT JOIN lottery_prizes ON lottery_records.prize_id = lottery_prizes.id
		;
	`
	)
	var (
		lotteryRecordsDTOS = make([]*dto.LotteryRecordsDTO, 0)
		countResult        int64
	)
	err := repo.DB().Raw(dataSQL, params.Limit(), params.Offset()).Scan(&lotteryRecordsDTOS).Error
	if err != nil {
		ctx.Logger().Errorf("LotteryRecordsRepo.ListAll error: %v", err)
		return nil, 0, errors.Wrap(err, "LotteryRecordsRepo.ListAll")
	}
	err = repo.DB().Raw(countSQL).Scan(&countResult).Error
	if err != nil {
		ctx.Logger().Errorf("LotteryRecordsRepo.ListAll error: %v", err)
		return nil, 0, errors.Wrap(err, "LotteryRecordsRepo.ListAll")
	}
	return lotteryRecordsDTOS, countResult, nil
}
