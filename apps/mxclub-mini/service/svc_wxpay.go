package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/domain/order/repo"
	"mxclub/pkg/common/wxpay"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewWxPayService)
}

type WxPayService struct {
	orderRepo   repo.IOrderRepo
	userService *UserService
}

func NewWxPayService(orderRepo repo.IOrderRepo, userService *UserService) *WxPayService {
	return &WxPayService{orderRepo: orderRepo, userService: userService}
}

func (s WxPayService) Prepay(ctx jet.Ctx, id uint, amount float64) (*wxpay.PrePayDTO, error) {
	userPO, _ := s.userService.FindUserById(id)
	prePayRequestDTO := wxpay.NewPrepayRequest(amount, userPO.WxOpenId)
	prepayDTO, err := wxpay.Prepay(ctx, prePayRequestDTO)
	if err != nil {
		ctx.Logger().Errorf("[WxPayService]prepay ERROR: %v\nprepayDTO:%v", err.Error(), utils.ObjToJsonStr(prepayDTO))
		return nil, errors.New("申请微信支付失败")
	}
	return prepayDTO, nil
}
