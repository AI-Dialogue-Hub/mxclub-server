package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mxclub/domain/lottery/po"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
	"sync"
)

func init() {
	jet.Provide(NewLotteryPrizeRepo)
}

type ILotteryPrizeRepo interface {
	xmysql.IBaseRepo[po.LotteryPrize]
	RemoveByPrizeIds(ctx jet.Ctx, ids []uint) (delCount int64, err error)
	ListByActivityId(ctx jet.Ctx, activityId uint, params *api.PageParams) ([]*po.LotteryPrize, int64, error)
	FindByIds(ctx jet.Ctx, ids []uint) (map[uint]*po.LotteryPrize, error)
}

func NewLotteryPrizeRepo(db *gorm.DB) ILotteryPrizeRepo {
	repo := new(LotteryPrizeRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.LotteryPrize)
	repo.Ctx = context.Background()
	return repo
}

type LotteryPrizeRepo struct {
	xmysql.BaseRepo[po.LotteryPrize]
}

func (repo *LotteryPrizeRepo) RemoveByPrizeIds(ctx jet.Ctx, ids []uint) (delCount int64, err error) {
	tx := repo.DB().Where("id in (?)", ids).Delete(repo.ModelPO)
	tx.Scan(&delCount)
	return delCount, tx.Error
}

func (repo *LotteryPrizeRepo) ListByActivityId(
	ctx jet.Ctx, activityId uint, params *api.PageParams,
) ([]*po.LotteryPrize, int64, error) {
	var (
		dataSQL = `SELECT * FROM lottery_prizes
				   LEFT JOIN activity_prize_relations
				   ON lottery_prizes.id = activity_prize_relations.prize_id
				   WHERE activity_prize_relations.activity_id = ? AND lottery_prizes.deleted_at IS NULL
				   LIMIT ? OFFSET ?
				   ;`
		countSQL = `SELECT COUNT(1) FROM lottery_prizes
					LEFT JOIN activity_prize_relations
					ON lottery_prizes.id = activity_prize_relations.prize_id
					WHERE activity_prize_relations.activity_id = ? AND lottery_prizes.deleted_at IS NULL
				   ;`
	)
	var (
		datas = make([]*po.LotteryPrize, 0)
		count int64
		wg    = new(sync.WaitGroup)
		err   error
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		err = repo.DB().Raw(dataSQL, activityId, params.Limit(), params.Offset()).Scan(&datas).Error
		if err != nil {
			ctx.Logger().Errorf("LotteryPrizeRepo.ListPrizeByActivityId Raw error: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		err = repo.DB().Raw(countSQL, activityId).Scan(&count).Error
		if err != nil {
			ctx.Logger().Errorf("LotteryPrizeRepo.ListPrizeByActivityId Raw error: %v", err)
		}
	}()
	wg.Wait()
	if err != nil {
		ctx.Logger().Errorf("LotteryPrizeRepo.ListPrizeByActivityId Raw error: %v", err)
		return nil, 0, errors.Wrap(err, "LotteryPrizeRepo.ListPrizeByActivityId")
	}
	return datas, count, nil
}

func (repo *LotteryPrizeRepo) FindByIds(ctx jet.Ctx, ids []uint) (map[uint]*po.LotteryPrize, error) {
	query := xmysql.NewMysqlQuery()
	query.SetFilter("id in ?", ids)
	data, err := repo.ListNoCountByQuery(query)
	if err != nil {
		ctx.Logger().Errorf("LotteryPrizeRepo.FindByIds Raw error: %v", err)
		return nil, errors.Wrap(err, "LotteryPrizeRepo.FindByIds")
	}
	return utils.SliceToRawMap(data,
		func(ele *po.LotteryPrize) uint { return ele.ID }), nil
}
