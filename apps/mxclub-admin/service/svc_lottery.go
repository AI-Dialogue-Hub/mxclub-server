package service

import (
	"github.com/fengyuan-liang/GoKit/collection/stream"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/pkg/errors"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/lottery/ability"
	"mxclub/domain/lottery/entity/dto"
	"mxclub/domain/lottery/entity/enum"
	"mxclub/domain/lottery/po"
	"mxclub/domain/lottery/repo"
	productRepo "mxclub/domain/product/repo"
	userRepo "mxclub/domain/user/repo"
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
	lotteryAbility      ability.ILotteryAbility
	lotteryRecordsRepo  repo.ILotteryRecordsRepo
	userRepo            userRepo.IUserRepo
}

func NewLotteryService(
	lotteryPrizeRepo repo.ILotteryPrizeRepo,
	lotteryActivityRepo repo.ILotteryActivityRepo,
	lotteryRepo repo.ILotteryRepo,
	lotteryActivity ability.ILotteryAbility,
	productRepo productRepo.IProductRepo,
	lotteryRecordsRepo repo.ILotteryRecordsRepo,
	userRepo userRepo.IUserRepo) *LotteryService {
	return &LotteryService{
		lotteryPrizeRepo:    lotteryPrizeRepo,
		lotteryActivityRepo: lotteryActivityRepo,
		lotteryRepo:         lotteryRepo,
		lotteryAbility:      lotteryActivity,
		productRepo:         productRepo,
		lotteryRecordsRepo:  lotteryRecordsRepo,
		userRepo:            userRepo,
	}
}

func (svc *LotteryService) FetchLotteryPrizeType() *vo.LotteryTypeVO {
	prizeType := svc.lotteryAbility.FetchLotteryPrizeType()
	options := make([]vo.Option, 0)
	prizeType.PrizeType.ForEach(func(key enum.PrizeTypeEnum, value string) {
		options = append(options, vo.Option{Label: value, Value: string(key)})
	})
	return &vo.LotteryTypeVO{LotteryType: options}
}

func (svc *LotteryService) ListPrize(ctx jet.Ctx, params *req.LotteryPrizePageReq) ([]*vo.LotteryPrizeVO, int64, error) {
	var (
		list  []*po.LotteryPrize
		count int64
		err   error
	)
	if params.ActivityId <= 0 {
		list, count, err = svc.lotteryPrizeRepo.List(params.Page, params.PageSize, nil)
	} else {
		list, count, err = svc.lotteryPrizeRepo.ListByActivityId(ctx, params.ActivityId, params.PageParams)
	}
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
	if req.PrizeType == enum.Virtual {
		if req.ProductAttributeID <= 0 {
			return errors.New("请选择商品")
		}
		if productPO, err := svc.productRepo.FindByID(req.ProductAttributeID); err == nil {
			req.PrizeName = productPO.Title
		}
	}
	if req.ActivityId <= 0 {
		return errors.New("请选择活动")
	}
	// 添加
	wrappedPO := wrap2PO(req)
	err := svc.lotteryAbility.AddPrize(ctx, req.ActivityId, wrappedPO)
	if err != nil {
		ctx.Logger().Errorf(
			"[LotteryService#AddOrUpdatePrize] LotteryPrizeReq is %v, insert error, %v", utils.ObjToJsonStr(req), err)
		return errors.New("添加失败")
	}
	return nil
}

func (svc *LotteryService) DelPrize(ctx jet.Ctx, req *req.LotteryPrizeReq) error {
	if err := svc.lotteryAbility.DelPrize(ctx, req.Id); err != nil {
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

// =====================================================================

func (svc *LotteryService) AddOrUpdateActivity(ctx jet.Ctx, req *req.LotteryActivityReq) error {
	if err := svc.lotteryAbility.AddOrUpdateActivity(ctx, wrapActivity(req)); err != nil {
		ctx.Logger().Errorf("AddOrUpdateActivity, err:%v", err)
		return errors.New("添加失败")
	}
	return nil
}

func wrapActivity(req *req.LotteryActivityReq) *po.LotteryActivity {
	return &po.LotteryActivity{
		ID:                  req.ID,
		FallbackPrizeId:     req.FallbackPrizeId,
		ActivityPrice:       req.ActivityPrice,
		ActivityTitle:       req.ActivityTitle,
		ActivitySubtitle:    req.ActivitySubtitle,
		ActivityDesc:        req.ActivityDesc,
		EntryURL:            req.EntryURL,
		EntryImage:          req.EntryImage,
		BannerImage:         req.BannerImage,
		BackgroundImage:     req.BackgroundImage,
		ActivityRules:       req.ActivityRules,
		PrizePoolID:         req.PrizePoolID,
		StartTime:           req.StartTime,
		EndTime:             req.EndTime,
		ParticipateTimes:    req.ParticipateTimes,
		ShareAddTimes:       req.ShareAddTimes,
		TotalPrizeCount:     req.TotalPrizeCount,
		RemainingPrizeCount: req.RemainingPrizeCount,
		ActivityStatus:      req.ActivityStatus,
		DisplayOrder:        req.DisplayOrder,
		IsFeatured:          req.IsFeatured,
		IsHot:               req.IsHot,
		SalesVolume:         req.SalesVolume,
	}
}

func (svc *LotteryService) ListActivity(ctx jet.Ctx, params *api.PageParams) ([]*vo.LotteryActivityVO, int64, error) {
	list, count, err := svc.lotteryAbility.ListActivity(ctx, params)
	if err != nil {
		ctx.Logger().Errorf("ListActivity error: %v", err)
		return nil, 0, errors.Wrap(err, "ListActivity 错误")
	}
	vos := utils.CopySlice[*po.LotteryActivity, *vo.LotteryActivityVO](list)
	fallbackIds := stream.Of[*po.LotteryActivity, uint](list).
		Filter(func(ele *po.LotteryActivity) bool { return ele.FallbackPrizeId > 0 }).
		Map(func(ele *po.LotteryActivity) uint { return ele.FallbackPrizeId }).
		CollectToSlice()
	if fallbackIds != nil && len(fallbackIds) > 0 {
		if id2PrizeMap, err := svc.lotteryPrizeRepo.FindByIds(ctx, fallbackIds); err == nil && id2PrizeMap != nil {
			for _, v := range vos {
				if fallbackPrize, ok := id2PrizeMap[v.FallbackPrizeId]; ok {
					v.FallbackPrizeName = fallbackPrize.PrizeLevel.String() + " " + fallbackPrize.PrizeName
				}
			}
		}
	}
	return vos, count, nil
}

func (svc *LotteryService) FindActivityById(ctx jet.Ctx, activityId uint) (*vo.LotteryActivityPrizeVO, error) {
	data, err := svc.lotteryAbility.FindActivityPrizeByActivityId(ctx, activityId)
	if err != nil {
		return nil, errors.New("活动获取错误")
	}
	activityPrizeVO := &vo.LotteryActivityPrizeVO{
		LotteryActivity: utils.MustCopy[vo.LotteryActivityVO](data.LotteryActivity),
		LotteryPrizes:   utils.CopySlice[*po.LotteryPrize, *vo.LotteryPrizeVO](data.LotteryPrizes),
	}
	prizes := data.LotteryPrizes
	productIds := utils.Map[*po.LotteryPrize, uint64](prizes, func(in *po.LotteryPrize) uint64 {
		return in.ProductAttributeID
	})
	if productList, err := svc.productRepo.FindByIds(ctx, productIds); err == nil {
		for _, prizeVO := range activityPrizeVO.LotteryPrizes {
			if product, ok := productList[prizeVO.ProductAttributeID]; ok {
				prizeVO.PrizeInfo = product.Description
			}
		}
	}
	return activityPrizeVO, nil
}

func (svc *LotteryService) UpdateActivityStatus(ctx jet.Ctx, req *req.LotteryActivityStatusReq) error {
	// 如果改为进行中，需要检查要有8个奖品，并且奖品概率小于1
	if req.LotteryActivityStatus == enum.Ongoing {
		activityPrizeDTO, err := svc.lotteryAbility.FindActivityPrizeByActivityId(ctx, req.LotteryActivityId)
		if err != nil || activityPrizeDTO == nil || activityPrizeDTO.LotteryActivity == nil {
			ctx.Logger().Errorf("FindActivityPrizeByActivityId error: %v", err)
			return errors.New("活动查询失败")
		}
		lotteryPrizes := activityPrizeDTO.LotteryPrizes
		if lotteryPrizes == nil || len(lotteryPrizes) < 8 {
			ctx.Logger().Errorf("prizes number less 8")
			return errors.New("奖品数量不足")
		}
		var (
			displayProbability, actualProbability float64
		)
		for _, lotteryPrize := range lotteryPrizes {
			actualProbability += lotteryPrize.ActualProbability
			displayProbability += lotteryPrize.DisplayProbability
		}
		if actualProbability > 1.0 {
			ctx.Logger().Errorf("prize actualProbability greater than 1")
			return errors.New("奖品池实际概率需要小于1")
		}
		if displayProbability > 1.0 {
			ctx.Logger().Errorf("prize displayProbability greater than 1")
			return errors.New("奖品池展示概率需要小于1")
		}
		if activityPrizeDTO.LotteryActivity.FallbackPrizeId <= 0 {
			ctx.Logger().Errorf("activity:%v, fallbackPrizeId is empty", activityPrizeDTO.LotteryActivity.ID)
			return errors.New("无抽奖三次必中奖品")
		}
	}
	if err := svc.lotteryActivityRepo.UpdateStatus(ctx, req.LotteryActivityId, req.LotteryActivityStatus); err != nil {
		ctx.Logger().Errorf("UpdateActivityStatus error: %v", err)
		return errors.New("更新失败")
	}
	return nil
}

func (svc *LotteryService) DelActivity(ctx jet.Ctx, req *req.LotteryActivityReq) error {
	if err := svc.lotteryAbility.DelActivity(ctx, req.ID); err != nil {
		ctx.Logger().Errorf("DelActivity error: %v", err)
		return errors.New("删除失败")
	}
	return nil
}

// =======================================================================

func (svc *LotteryService) ListLotteryRecords(ctx jet.Ctx, params *api.PageParams) ([]*vo.LotteryRecordsVO, int64, error) {
	list, count, err := svc.lotteryRecordsRepo.ListRecords(ctx, params)
	if err != nil {
		return nil, 0, errors.New("查找失败")
	}
	vos := utils.CopySlice[*dto.LotteryRecordsDTO, *vo.LotteryRecordsVO](list)
	for _, recordsVO := range vos {
		if userPO, err := svc.userRepo.FindByIdAroundCache(ctx, recordsVO.UserId); err == nil && userPO != nil {
			recordsVO.AvatarUrl = userPO.WxIcon
		}
	}
	return vos, count, nil
}
