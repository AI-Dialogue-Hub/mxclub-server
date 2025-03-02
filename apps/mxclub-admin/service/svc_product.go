package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/product/po"
	"mxclub/domain/product/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewProductService)
}

type ProductService struct {
	productRepo      repo.IProductRepo
	productSalesRepo repo.IProductSalesRepo
}

func NewProductService(repo repo.IProductRepo, productSalesRepo repo.IProductSalesRepo) *ProductService {
	return &ProductService{
		productRepo:      repo,
		productSalesRepo: productSalesRepo,
	}
}

// =============================================================

func (s ProductService) List(ctx jet.Ctx, params *req.ProductListReq) (*api.PageResult, error) {
	pageParams := &api.PageParams{Page: params.Page, PageSize: params.PageSize}
	list, count, err := s.productRepo.ListAroundCache(ctx, pageParams, utils.ParseUint(params.ProductType))
	if err != nil {
		return nil, err
	}
	productVOS := utils.CopySlice[*po.Product, *vo.ProductVO](list)
	productIds := utils.Map[*vo.ProductVO, uint64](productVOS, func(in *vo.ProductVO) uint64 {
		return in.ID
	})
	// 销量
	id2ProductSaleMap, err := s.productSalesRepo.FindByProductIds(ctx, productIds)
	if err == nil && id2ProductSaleMap != nil && !id2ProductSaleMap.IsEmpty() {
		utils.ForEach(productVOS, func(ele *vo.ProductVO) {
			if productSalePO, ok := id2ProductSaleMap.Get(ele.ID); ok {
				ele.Sale = int(productSalePO.SalesVolume)
			}
		})
	}
	return api.WrapPageResult(pageParams, productVOS, count), nil
}

func (s ProductService) Update(ctx jet.Ctx, req *req.ProductReq) error {
	updateMap := utils.ObjToMap(req)
	finalPrice := req.Price - req.DiscountPrice
	if finalPrice <= 0 {
		return errors.New("价格错误，折扣价不能大于等于原价")
	}
	updateMap["final_price"] = finalPrice
	updateMap["images"] = req.Images
	updateMap["detail_images"] = req.DetailImages
	return s.productRepo.UpdateProduct(ctx, updateMap)
}

func (s ProductService) UpdateHotInfo(ctx jet.Ctx, req *req.ProductHotReq) error {
	updateMap := utils.ObjToMap(req)
	return s.productRepo.UpdateProduct(ctx, updateMap)
}

func (s ProductService) DeleteById(ctx jet.Ctx, id int64) error {
	return s.productRepo.DeleteById(ctx, id)
}

func (s ProductService) Add(ctx jet.Ctx, productReq *req.ProductReq) error {
	product := utils.MustCopy[po.Product](productReq)
	product.FinalPrice = product.Price - product.DiscountPrice
	return s.productRepo.Add(ctx, product)
}

func (s ProductService) UpdateSales(ctx jet.Ctx, req *req.ProductSaleReq) error {
	if err := s.productSalesRepo.ReplaceSale(ctx, req.ProductId, req.Sale); err != nil {
		ctx.Logger().Errorf("UpdateSales failed, err:%v", err)
		return errors.New("更新失败")
	}
	return nil
}
