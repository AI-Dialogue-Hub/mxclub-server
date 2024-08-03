package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/product/po"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/common/xredis"
)

func init() {
	jet.Provide(NewProductRepo)
}

type IProductRepo interface {
	xmysql.IBaseRepo[po.Product]
	ListAroundCache(ctx jet.Ctx, params *api.PageParams, productType uint) ([]*po.Product, int64, error)
	UpdateProduct(ctx jet.Ctx, updateMap map[string]any) error
	DeleteById(ctx jet.Ctx, id int64) error
	Add(ctx jet.Ctx, po *po.Product) error
}

func NewProductRepo(db *gorm.DB) IProductRepo {
	repo := new(ProductRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.Product)
	repo.Ctx = context.Background()
	return repo
}

type ProductRepo struct {
	xmysql.BaseRepo[po.Product]
}

// ===========================================================

const productCachePrefix = "mini_product"
const productListCachePrefix = productCachePrefix + "_list"

func (repo ProductRepo) ListAroundCache(ctx jet.Ctx, params *api.PageParams, productType uint) ([]*po.Product, int64, error) {
	// 根据页码参数生成唯一的缓存键
	cacheListKey := xredis.BuildListDataCacheKey(fmt.Sprintf("%v_type_%v", productCachePrefix, productType), params)
	cacheCountKey := xredis.BuildListCountCacheKey(fmt.Sprintf("%v_type_%v", productListCachePrefix, productType))

	list, count, err := xredis.GetListOrDefault[po.Product](ctx, cacheListKey, cacheCountKey, func() ([]*po.Product, int64, error) {
		// 如果缓存中未找到，则从数据库中获取
		query := xmysql.NewMysqlQuery()
		query.SetPage(params.Page, params.PageSize)
		if productType != 0 {
			query.SetFilter("type = ?", productType)
		}
		list, count, err := repo.ListByWrapper(ctx, query)
		if err != nil {
			return nil, 0, err
		}
		return list, count, nil
	})
	if err != nil {
		ctx.Logger().Errorf("ListAroundCache 错误: %v", err)
		return nil, 0, err
	}

	return list, count, nil
}

func (repo ProductRepo) UpdateProduct(ctx jet.Ctx, updateMap map[string]any) error {
	_ = xredis.DelMatchingKeys(ctx, productCachePrefix)
	id := updateMap["id"]
	delete(updateMap, "id")
	return repo.Update(updateMap, "id = ?", id)
}

func (repo ProductRepo) DeleteById(ctx jet.Ctx, id int64) error {
	_ = xredis.DelMatchingKeys(ctx, productCachePrefix)
	err := repo.RemoveByID(id)
	if err != nil {
		ctx.Logger().Errorf("[ProductRepo]DeleteById ERROR:%v", err.Error())
		return errors.New("删除失败")
	}
	return nil
}

func (repo ProductRepo) Add(ctx jet.Ctx, po *po.Product) error {
	_ = xredis.DelMatchingKeys(ctx, productCachePrefix)
	err := repo.InsertOne(po)
	if err != nil {
		ctx.Logger().Errorf("[ProductRepo]Add ERROR:%v", err.Error())
		return errors.New("添加失败")
	}
	return nil
}
