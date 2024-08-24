package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/order/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewTransferRepo)
}

type ITransferRepo interface {
	xmysql.IBaseRepo[po.OrderTransfer]
}

func NewTransferRepo(db *gorm.DB) ITransferRepo {
	repo := new(TransferRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.OrderTransfer)
	repo.Ctx = context.Background()
	return repo
}

type TransferRepo struct {
	xmysql.BaseRepo[po.OrderTransfer]
}
