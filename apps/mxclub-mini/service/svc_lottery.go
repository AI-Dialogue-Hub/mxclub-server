package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/lottery/ability"
	"mxclub/domain/lottery/entity/dto"
	"mxclub/domain/lottery/entity/enum"
	"mxclub/domain/lottery/po"
	"mxclub/domain/lottery/repo"
	"mxclub/domain/lottery/strategy"
	userRepo "mxclub/domain/user/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewLotteryService)
}

type LotteryService struct {
	lotteryPrizeRepo   repo.ILotteryPrizeRepo
	lotteryAbility     ability.ILotteryAbility
	messageService     *MessageService
	lotteryRecordsRepo repo.ILotteryRecordsRepo
	userRepo           userRepo.IUserRepo
}

func NewLotteryService(
	lotteryPrizeRepo repo.ILotteryPrizeRepo,
	lotteryActivity ability.ILotteryAbility,
	messageService *MessageService,
	lotteryRecordsRepo repo.ILotteryRecordsRepo,
	userRepo userRepo.IUserRepo) *LotteryService {
	return &LotteryService{
		lotteryPrizeRepo:   lotteryPrizeRepo,
		lotteryAbility:     lotteryActivity,
		messageService:     messageService,
		lotteryRecordsRepo: lotteryRecordsRepo,
		userRepo:           userRepo,
	}
}

func (svc *LotteryService) ListLotteryPrize(ctx jet.Ctx, params *api.PageParams) ([]*vo.LotteryActivityPrizeVO, int64, error) {
	listActivity, count, err := svc.lotteryAbility.ListActivityPrize(ctx, params)
	if err != nil {
		ctx.Logger().Errorf("[LotteryService#ListLotteryPrize] ERROR:%v", err)
		return nil, 0, errors.New("活动获取错误")
	}
	lotteryActivityPrizeVOS := utils.Map(listActivity, func(activityDTO *dto.LotteryActivityDTO) *vo.LotteryActivityPrizeVO {
		return &vo.LotteryActivityPrizeVO{
			LotteryActivity: utils.MustCopy[vo.LotteryActivityVO](activityDTO.LotteryActivity),
			LotteryPrizes:   utils.CopySlice[*po.LotteryPrize, *vo.LotteryPrizeVO](activityDTO.LotteryPrizes),
		}
	})
	return lotteryActivityPrizeVOS, count, nil
}

func (svc *LotteryService) FindActivityPrizeByActivityId(ctx jet.Ctx, activityId int) (*vo.LotteryActivityPrizeVO, error) {
	activityDTO, err := svc.lotteryAbility.FindActivityPrizeByActivityId(ctx, uint(activityId))
	if err != nil {
		return nil, errors.New("活动获取错误")
	}
	return utils.MustCopy[vo.LotteryActivityPrizeVO](activityDTO), nil
}

func (svc *LotteryService) StartLottery(ctx jet.Ctx, req *req.LotteryStartReq) (*vo.LotteryVO, error) {
	if req.ActivityId <= 0 {
		ctx.Logger().Errorf("[LotteryService#StartLottery] ActivityId:%v is illegal", req.ActivityId)
		return nil, errors.New("请选择活动")
	}
	var (
		userId = middleware.MustGetUserId(ctx)
	)
	lotteryStrategy, _ := strategy.FetchLotteryStrategy(strategy.RandomStrategy)
	drawResultDTO, err := lotteryStrategy.DoDraw(ctx, &dto.LotteryStrategyDrawDTO{
		UserId:     middleware.MustGetUserId(ctx),
		ActivityId: req.ActivityId,
	})
	ctx.Logger().Infof("[LotteryService#StartLottery] DoDraw result:%v", utils.ObjToJsonStr(drawResultDTO))
	if err != nil || drawResultDTO == nil {
		ctx.Logger().Errorf("[LotteryService#StartLottery] DoDraw err:%v, drawResultDTO:%v", err, drawResultDTO)
		return nil, errors.New("抽奖失败，请联系客服")
	}
	prize := drawResultDTO.LotteryPrize
	// todo@lfy 这里如果发货失败了，会有问题
	go svc.DistributePrize(ctx, userId, req.ActivityId, prize)
	return &vo.LotteryVO{
		PrizeIndex: prize.SortOrder,
		WinMessage: prize.WinMessage,
		PrizeImage: prize.PrizeImage,
	}, nil
}

// DistributePrize 发放奖品
//
// 目前就代打订单和实物两种奖品
func (svc *LotteryService) DistributePrize(ctx jet.Ctx, userId uint, activityId uint, prize *po.LotteryPrize) {
	// 1. 执行发奖策略
	switch prize.PrizeType {
	case enum.Physical:
		_ = svc.messageService.PushSystemMessage(ctx, userId,
			fmt.Sprintf("恭喜您抽中%s，请联系客服进行领取", prize.PrizeName))
		// 扣减库存
	case enum.Virtual:
		_ = svc.messageService.PushSystemMessage(ctx, userId,
			fmt.Sprintf("恭喜您抽中单:%s，打手接单后会直接开始订单", prize.PrizeName))

	}

}

func (svc *LotteryService) handleAddPrizeToOrder(ctx jet.Ctx, userId uint, activityId uint, prize *po.LotteryPrize) {
	// 1. 查找用户最新一次的购买记录
	purchaseRecords, err := svc.lotteryAbility.FindPurchaseRecord(ctx, userId, activityId, false)
	if err != nil {
		ctx.Logger().Errorf("find purchase record error")
		_ = svc.messageService.PushSystemMessage(ctx, userId, "系统异常，抽奖失败，请联系客服")
	}
	ctx.Logger().Infof("purchase records: %v", purchaseRecords)
	//// 2. 创建订单
	//order := &orderPO.Order{
	//	OrderId:         orderTradeNo,
	//	PurchaseId:      userId,
	//	OrderName:       req.OrderName,
	//	OrderIcon:       req.OrderIcon,
	//	OrderStatus:     status,
	//	OriginalPrice:   preferentialVO.OriginalPrice,
	//	ProductID:       req.ProductId,
	//	Phone:           req.Phone,
	//	GameRegion:      req.GameRegion,
	//	RoleId:          req.RoleId,
	//	SpecifyExecutor: req.SpecifyExecutor,
	//	ExecutorID:      -1,
	//	Executor2Id:     -1,
	//	Executor3Id:     -1,
	//	ExecutorName:    "",
	//	Notes:           req.Notes,
	//	DiscountPrice:   preferentialVO.OriginalPrice - preferentialVO.DiscountedPrice,
	//	FinalPrice:      preferentialVO.DiscountedPrice,
	//	ExecutorPrice:   0,
	//	PurchaseDate:    utils.Ptr(time.Now()),
	//	GrabAt:          nil,
	//	OutRefundNo:     utils.ParseString(executorId), // 保存下用户选择的打手
	//}
}

// ============================================================

// ListLotteryRecords todo@lfy 随机找100条记录
func (svc *LotteryService) ListLotteryRecords(ctx jet.Ctx) ([]*vo.LotteryRecordsVO, int64, error) {
	list, count, err := svc.lotteryRecordsRepo.ListRecords(ctx, &api.PageParams{})
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

// ===========================================================

// CanDrawLottery 能否进行抽奖
func (svc *LotteryService) CanDrawLottery(ctx jet.Ctx, req *req.LotteryCanDrawReq) (bool, error) {
	userId := middleware.MustGetUserId(ctx)
	record, err := svc.lotteryAbility.FindPurchaseRecord(ctx, userId, req.ActivityId, true)
	if err != nil {
		ctx.Logger().Errorf("find purse chase record error: %v", err)
		return false, err
	}
	return len(record) > 0, nil
}
