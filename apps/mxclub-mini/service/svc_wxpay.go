package service

import (
	"errors"
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
	"mxclub/pkg/common/wxpay"
	"mxclub/pkg/utils"
	"time"
)

func init() {
	jet.Provide(NewWxPayService)
}

type WxPayService struct {
	wxpayCallbackRepo repo.IWxPayCallbackRepo
	userService       *UserService
	orderService      *OrderService
}

func NewWxPayService(
	orderService *OrderService,
	userService *UserService,
	wxpayCallbackRepo repo.IWxPayCallbackRepo) *WxPayService {
	return &WxPayService{
		orderService:      orderService,
		userService:       userService,
		wxpayCallbackRepo: wxpayCallbackRepo,
	}
}

func (s WxPayService) Prepay(ctx jet.Ctx, id uint, orderReq *req.OrderReq) (*wxpay.PrePayDTO, error) {
	// 给用户创建订单
	orderPO, err := s.orderService.AddByOrderStatus(ctx, orderReq, enum.PrePay)
	if err != nil {
		return nil, err
	}
	userPO, _ := s.userService.FindUserById(ctx, id)
	prePayRequestDTO := wxpay.NewPrepayRequest(orderPO.FinalPrice, userPO.WxOpenId, utils.ParseString(orderPO.OrderId))
	prepayDTO, err := wxpay.Prepay(ctx, prePayRequestDTO)
	if err != nil {
		ctx.Logger().Errorf("[WxPayService]prepay ERROR: %v\nprepayDTO:%v", err.Error(), utils.ObjToJsonStr(prepayDTO))
		return nil, errors.New("申请微信支付失败")
	}
	ctx.Logger().Infof("用户: %v 付款：%v，进行中，prepayDTO：%v", id, orderPO.FinalPrice, utils.ObjToJsonStr(prepayDTO))

	return prepayDTO, nil
}

func (s WxPayService) addRawOrder(ctx jet.Ctx, outTradeNo string, preferentialVO *vo.PreferentialVO, productId uint) {
	// 插入订单数据
	// 2. 创建订单
	order := &po.Order{
		OrderId:         utils.SafeParseUint64(outTradeNo),
		PurchaseId:      middleware.MustGetUserId(ctx),
		OrderName:       "",
		OrderIcon:       "",
		OrderStatus:     enum.PROCESSING,
		OriginalPrice:   preferentialVO.OriginalPrice,
		ProductID:       productId,
		Phone:           "",
		GameRegion:      "",
		RoleId:          "",
		SpecifyExecutor: false,
		ExecutorID:      -1,
		Executor2Id:     -1,
		Executor3Id:     -1,
		ExecutorName:    "",
		Notes:           "",
		DiscountPrice:   preferentialVO.OriginalPrice - preferentialVO.DiscountedPrice,
		FinalPrice:      preferentialVO.DiscountedPrice,
		ExecutorPrice:   0,
		PurchaseDate:    utils.Ptr(time.Now()),
	}
	if err := s.orderService.orderRepo.InsertOne(order); err != nil {
		ctx.Logger().Errorf("addRawOrder ERROR:%v", err)
	}
}

func (s WxPayService) HandleWxpayNotify(ctx jet.Ctx, params *maps.LinkedHashMap[string, any]) {
	defer utils.RecoverAndLogError(ctx)
	// 解析回调参数
	transaction, err := wxpay.DecryptWxpayCallBack(ctx)
	if err != nil || transaction == nil {
		ctx.Logger().Errorf("[DecryptWxpayCallBack] %v", err)
		// 失败了，直接解析参数
		if transaction, err = wxpay.DecryptWxpayCallBackByParams(ctx, params); err != nil {
			ctx.Logger().Errorf("[DecryptWxpayCallBackByParams]ERROR %v", err)
			return
		}
	}

	callbackInfo, err := s.wxpayCallbackRepo.FindByTraceNo(*transaction.OutTradeNo)

	// 幂等保护
	if err == nil && callbackInfo != nil && callbackInfo.ID > 0 {
		ctx.Logger().Errorf("duplicate callback, %v", utils.ObjToJsonStr(*transaction))
		return
	}

	// 修改订单状态为支付成功
	_ = s.orderService.PaySuccessOrder(ctx, utils.SafeParseUint64(*transaction.OutTradeNo))

	objToMap := utils.ObjToMap(*transaction)
	err = s.wxpayCallbackRepo.InsertOne(&po.WxPayCallback{
		OutTradeNo: *transaction.OutTradeNo,
		RawData:    objToMap,
	})
	ctx.Logger().Infof("HandleWxpayNotify:%v", utils.ObjToJsonStr(objToMap))
	if err != nil {
		ctx.Logger().Errorf("[DecryptWxpayCallBack] %v", err)
	}
	return
}
