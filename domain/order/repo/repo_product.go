package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/order/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewProductRepo)
}

type IProductRepo interface {
	xmysql.IBaseRepo[po.Product]
}

func NewProductRepo(db *gorm.DB) IProductRepo {
	repo := new(ProductRepo)
	repo.Db = db.Model(new(po.Product))
	repo.Ctx = context.Background()
	return repo
}

type ProductRepo struct {
	xmysql.BaseRepo[po.Product]
}
