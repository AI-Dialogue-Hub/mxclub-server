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
	productRepo repo.IProductRepo
}

func NewProductService(repo repo.IProductRepo) *ProductService {
	return &ProductService{productRepo: repo}
}

// =============================================================

func (s ProductService) List(ctx jet.Ctx, params *req.ProductListReq) (*api.PageResult, error) {
	pageParams := &api.PageParams{Page: params.Page, PageSize: params.PageSize}
	list, count, err := s.productRepo.ListAroundCache(ctx, pageParams, utils.ParseUint(params.ProductType))
	if err != nil {
		return nil, err
	}
	return api.WrapPageResult(pageParams, utils.CopySlice[*po.Product, *vo.ProductVO](list), count), nil
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

func (s ProductService) DeleteById(ctx jet.Ctx, id int64) error {
	return s.productRepo.DeleteById(ctx, id)
}

func (s ProductService) Add(ctx jet.Ctx, productReq *req.ProductReq) error {
	product := utils.MustCopy[po.Product](productReq)
	product.FinalPrice = product.Price - product.DiscountPrice
	return s.productRepo.Add(ctx, product)
}
