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
	"sync"
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
		WHERE lottery_records.deleted_at IS NULL
		ORDER BY lottery_records.created_at DESC
		LIMIT ? OFFSET ?
		;
	`
		countSQL = `
		SELECT COUNT(1)
		FROM lottery_records
				 LEFT JOIN lottery_activities activity ON activity.id = lottery_records.activity_id
				 LEFT JOIN lottery_prizes ON lottery_records.prize_id = lottery_prizes.id
		WHERE lottery_records.deleted_at IS NULL
		;
	`
	)
	var (
		lotteryRecordsDTOS = make([]*dto.LotteryRecordsDTO, 0)
		countResult        int64
		wg                 = new(sync.WaitGroup)
		dataErr, countErr  error
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		dataErr = repo.DB().Raw(dataSQL, params.Limit(), params.Offset()).Scan(&lotteryRecordsDTOS).Error
		if dataErr != nil {
			ctx.Logger().Errorf("LotteryRecordsRepo.ListAll error: %v", dataErr)
		}
	}()
	go func() {
		defer wg.Done()
		countErr = repo.DB().Raw(countSQL).Scan(&countResult).Error
		if countErr != nil {
			ctx.Logger().Errorf("LotteryRecordsRepo.ListAll error: %v", countErr)
		}
	}()
	wg.Wait()
	if dataErr != nil {
		ctx.Logger().Errorf("LotteryRecordsRepo.ListAll error: %v", dataErr)
		return nil, 0, errors.Wrap(dataErr, "LotteryRecordsRepo.ListAll")
	}
	if countErr != nil {
		ctx.Logger().Errorf("LotteryRecordsRepo.ListAll error: %v", countErr)
		return nil, 0, errors.Wrap(countErr, "LotteryRecordsRepo.ListAll")
	}
	return lotteryRecordsDTOS, countResult, nil
}
