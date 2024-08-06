package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/order/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewWxPayCallbackRepo)
}

type IWxPayCallbackRepo interface {
	xmysql.IBaseRepo[po.WxPayCallback]
	FindByTraceNo(orderTradeNo string) (*po.WxPayCallback, error)
}

func NewWxPayCallbackRepo(db *gorm.DB) IWxPayCallbackRepo {
	repo := new(WxPayCallbackRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.WxPayCallback)
	repo.Ctx = context.Background()
	return repo
}

type WxPayCallbackRepo struct {
	xmysql.BaseRepo[po.WxPayCallback]
}

// ========================================================

func (repo WxPayCallbackRepo) FindByTraceNo(orderTradeNo string) (*po.WxPayCallback, error) {
	return repo.FindOne("out_trade_no = ?", orderTradeNo)
}
