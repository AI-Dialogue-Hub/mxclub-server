package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
	"mxclub/pkg/common/wxpay"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewWxPayService)
}

type WxPayService struct {
	orderRepo         repo.IOrderRepo
	wxpayCallbackRepo repo.IWxPayCallbackRepo
	userService       *UserService
}

func NewWxPayService(orderRepo repo.IOrderRepo, userService *UserService, wxpayCallbackRepo repo.IWxPayCallbackRepo) *WxPayService {
	return &WxPayService{orderRepo: orderRepo, userService: userService, wxpayCallbackRepo: wxpayCallbackRepo}
}

func (s WxPayService) Prepay(ctx jet.Ctx, id uint, amount float64) (*wxpay.PrePayDTO, error) {
	userPO, _ := s.userService.FindUserById(ctx, id)
	prePayRequestDTO := wxpay.NewPrepayRequest(amount, userPO.WxOpenId)
	prepayDTO, err := wxpay.Prepay(ctx, prePayRequestDTO)
	if err != nil {
		ctx.Logger().Errorf("[WxPayService]prepay ERROR: %v\nprepayDTO:%v", err.Error(), utils.ObjToJsonStr(prepayDTO))
		return nil, errors.New("申请微信支付失败")
	}
	ctx.Logger().Infof("用户:%v 付款：%v，进行中，prepayDTO：%v", id, amount, utils.ObjToJsonStr(prepayDTO))
	return prepayDTO, nil
}

func (s WxPayService) HandleWxpayNotify(ctx jet.Ctx) {
	defer utils.RecoverAndLogError(ctx)
	// 解析回调参数
	transaction, err := wxpay.DecryptWxpayCallBack(ctx)
	if err != nil || transaction == nil {
		ctx.Logger().Errorf("[DecryptWxpayCallBack] %v", err)
		return
	}
	objToMap := utils.ObjToMap(*transaction)
	err = s.wxpayCallbackRepo.InsertOne(&po.WxPayCallback{
		OutTradeNo: *transaction.OutTradeNo,
		RawData:    objToMap,
	})
	ctx.Logger().Infof("%v", utils.ObjToJsonStr(objToMap))
	if err != nil {
		ctx.Logger().Errorf("[DecryptWxpayCallBack] %v", err)
	}
	return
}
