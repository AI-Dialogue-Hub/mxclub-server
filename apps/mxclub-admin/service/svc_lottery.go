package service

import (
	"errors"
	"github.com/fengyuan-liang/GoKit/collection/stream"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/lottery/activity"
	"mxclub/domain/lottery/entity/enum"
	"mxclub/domain/lottery/po"
	"mxclub/domain/lottery/repo"
	productRepo "mxclub/domain/product/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewLotteryService)
}

type LotteryService struct {
	lotteryPrizeRepo    repo.ILotteryPrizeRepo
	lotteryActivityRepo repo.ILotteryActivityRepo
	lotteryRepo         repo.ILotteryRepo
	productRepo         productRepo.IProductRepo
	lotteryActivity     activity.ILotteryActivity
}

func NewLotteryService(
	lotteryPrizeRepo repo.ILotteryPrizeRepo,
	lotteryActivityRepo repo.ILotteryActivityRepo,
	lotteryRepo repo.ILotteryRepo,
	lotteryActivity activity.ILotteryActivity,
	productRepo productRepo.IProductRepo) *LotteryService {
	return &LotteryService{
		lotteryPrizeRepo:    lotteryPrizeRepo,
		lotteryActivityRepo: lotteryActivityRepo,
		lotteryRepo:         lotteryRepo,
		lotteryActivity:     lotteryActivity,
		productRepo:         productRepo,
	}
}

func (svc *LotteryService) FetchLotteryPrizeType() *vo.LotteryTypeVO {
	prizeType := svc.lotteryActivity.FetchLotteryPrizeType()
	options := make([]vo.Option, 0)
	prizeType.PrizeType.ForEach(func(key enum.PrizeTypeEnum, value string) {
		options = append(options, vo.Option{Label: value, Value: string(key)})
	})
	return &vo.LotteryTypeVO{LotteryType: options}
}

func (svc *LotteryService) ListPrize(ctx jet.Ctx, params *api.PageParams) ([]*vo.LotteryPrizeVO, int64, error) {
	list, count, err := svc.lotteryPrizeRepo.List(params.Page, params.PageSize, nil)
	if err != nil {
		ctx.Logger().Errorf("[LotteryService#ListPrize] list error, %v", err)
		return nil, 0, errors.New("查找失败")
	}
	if list == nil || len(list) == 0 {
		return make([]*vo.LotteryPrizeVO, 0, 0), count, nil
	}
	virtualPrizeIds := stream.Of[*po.LotteryPrize, uint64](list).
		Filter(func(ele *po.LotteryPrize) bool { return ele.PrizeType == enum.Virtual }).
		Map(func(ele *po.LotteryPrize) uint64 { return ele.ProductAttributeID }).
		CollectToSlice()
	productId2ProductMap, err := svc.productRepo.FindByIds(ctx, virtualPrizeIds)
	if err != nil {
		ctx.Logger().Errorf("[LotteryService#ListPrize] find order error, %v", err)
		return nil, 0, errors.New("查找失败")
	}
	if productId2ProductMap != nil && len(productId2ProductMap) > 0 {
		// 如果是虚拟类型的奖品，需要拼接商品信息
		for _, prizePO := range list {
			if prizePO.PrizeType != enum.Physical {
				continue
			}
			if product, ok := productId2ProductMap[prizePO.ProductAttributeID]; ok {
				prizePO.PrizeName = product.ShortDescription
				prizePO.PrizeValue = product.Price
				prizePO.PrizeImage = product.Thumbnail
			}
		}
	}
	return utils.CopySlice[*po.LotteryPrize, *vo.LotteryPrizeVO](list), count, nil
}

func (svc *LotteryService) AddOrUpdatePrize(ctx jet.Ctx, req *req.LotteryPrizeReq) error {
	if req.PrizeType == enum.Virtual && req.ProductAttributeID <= 0 {
		return errors.New("请选择商品")
	}
	if req.ActivityId <= 0 {
		return errors.New("请选择活动")
	}
	if req.Id > 0 {
		// 修改
		lotteryPrizePO := wrap2PO(req)
		if err := svc.lotteryPrizeRepo.UpdateById(lotteryPrizePO, req.Id); err != nil {
			ctx.Logger().Errorf("ERROR:%v", err.Error())
			return errors.New("更新失败")
		}
		return nil
	} else {
		// 新增
		return svc.AddPrize(ctx, req)
	}
}

func (svc *LotteryService) AddPrize(ctx jet.Ctx, req *req.LotteryPrizeReq) error {
	if req.PrizeType == enum.Virtual && req.ProductAttributeID <= 0 {
		return errors.New("请选择商品")
	}
	if req.ActivityId <= 0 {
		return errors.New("请选择活动")
	}
	// 添加
	wrappedPO := wrap2PO(req)
	err := svc.lotteryActivity.AddPrize(ctx, req.ActivityId, wrappedPO)
	if err != nil {
		ctx.Logger().Errorf(
			"[LotteryService#AddOrUpdatePrize] LotteryPrizeReq is %v, insert error, %v", utils.ObjToJsonStr(req), err)
		return errors.New("添加失败")
	}
	return nil
}

func (svc *LotteryService) DelPrize(ctx jet.Ctx, req *req.LotteryPrizeReq) error {
	if err := svc.lotteryActivity.DelPrize(ctx, req.Id); err != nil {
		ctx.Logger().Errorf(
			"[LotteryService#DelPrize] LotteryPrizeReq is %v, delete error, %v", utils.ObjToJsonStr(req), err)
		return errors.New("删除失败")
	}
	return nil
}

func wrap2PO(req *req.LotteryPrizeReq) *po.LotteryPrize {
	wrappedPO := &po.LotteryPrize{
		ProductAttributeID:    req.ProductAttributeID,
		PrizeLevel:            req.PrizeLevel,
		PrizeName:             req.PrizeName,
		PrizeType:             req.PrizeType,
		PrizeValue:            req.PrizeValue,
		TotalQuantity:         req.TotalQuantity,
		RemainingQuantity:     req.RemainingQuantity,
		DailyLimit:            req.DailyLimit,
		UserDailyLimit:        req.UserDailyLimit,
		UserTotalLimit:        req.UserTotalLimit,
		PrizeImage:            req.PrizeImage,
		WinMessage:            req.WinMessage,
		DisplayProbability:    req.DisplayProbability,
		ActualProbability:     req.ActualProbability,
		ProbabilityAdjustment: req.ProbabilityAdjustment,
		SortOrder:             req.SortOrder,
		IsActive:              req.IsActive,
		StartTime:             req.StartTime,
		EndTime:               req.EndTime,
	}
	return wrappedPO
}
