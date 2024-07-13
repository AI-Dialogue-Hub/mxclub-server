package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/domain/order/repo"
)

func init() {
	jet.Provide(NewWxPayService)
}

type WxPayService struct {
	orderRepo repo.IOrderRepo
}

func NewWxPayService(orderRepo repo.IOrderRepo) *WxPayService {
	return &WxPayService{orderRepo: orderRepo}
}
