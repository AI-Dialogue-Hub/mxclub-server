package repo

import (
	"fmt"
	"github.com/fengyuan-liang/GoKit/collection/maps"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/product/po"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
	"time"
)

const Default_Sale_Volume = 1

func init() {
	jet.Provide(NewProductSalesRepo)
}

type IProductSalesRepo interface {
	xmysql.IBaseRepo[po.ProductSale]

	AddOrUpdateSale(ctx jet.Ctx, productId uint, salesVolume int) (err error)
	FindByProductIds(ctx jet.Ctx, ids []uint64) (maps.IMap[uint64, *po.ProductSale], error)
	ReplaceSale(ctx jet.Ctx, productId uint, salesVolume int) error
}

func NewProductSalesRepo(db *gorm.DB) IProductSalesRepo {
	repo := new(ProductSalesRepoImpl)
	repo.SetDB(db)
	repo.ModelPO = new(po.ProductSale)
	return repo
}

type ProductSalesRepoImpl struct {
	xmysql.BaseRepo[po.ProductSale]
}

// ==============================================================

func (repo ProductSalesRepoImpl) AddOrUpdateSale(ctx jet.Ctx, productId uint, salesVolume int) (err error) {
	sql := fmt.Sprintf(
		`INSERT INTO %v (product_id, sales_volume, sale_date) VALUES (?, ?, ?) 
            ON DUPLICATE KEY UPDATE sales_volume = sales_volume + VALUES(sales_volume)`, repo.ModelPO.TableName(),
	)

	if err = repo.DB().Exec(sql, productId, salesVolume, time.Now()).Error; err != nil {
		ctx.Logger().Errorf("[ProductSales#AddOrUpdateSale]ERROR:%v", err)
	}

	return
}

func (repo ProductSalesRepoImpl) FindByProductIds(ctx jet.Ctx, ids []uint64) (maps.IMap[uint64, *po.ProductSale], error) {
	query := xmysql.NewMysqlQuery()
	query.SetFilter("product_id in (?)", ids)
	productSales, err := repo.ListNoCountByQuery(query)
	if err != nil || productSales == nil || len(productSales) == 0 {
		ctx.Logger().Errorf("find product ids failed, ids is %v", ids)
		return nil, err
	}
	sliceToMap := utils.SliceToSingleMap[*po.ProductSale, uint64](productSales, func(ele *po.ProductSale) uint64 {
		return uint64(ele.ProductID)
	})
	return sliceToMap, nil
}

func (repo ProductSalesRepoImpl) ReplaceSale(ctx jet.Ctx, productId uint, salesVolume int) (err error) {
	sql := fmt.Sprintf(
		`INSERT INTO %v (product_id, sales_volume, sale_date) VALUES (?, ?, ?) 
            ON DUPLICATE KEY UPDATE sales_volume = VALUES(sales_volume)`, repo.ModelPO.TableName(),
	)

	if err = repo.DB().Exec(sql, productId, salesVolume, time.Now()).Error; err != nil {
		ctx.Logger().Errorf("[ProductSales#AddOrUpdateSale]ERROR:%v", err)
	}

	return
}
