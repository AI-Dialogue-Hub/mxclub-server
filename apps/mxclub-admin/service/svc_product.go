package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
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

func (s ProductService) List(ctx jet.Ctx, params *api.PageParams) (*api.PageResult, error) {
	list, count, err := s.productRepo.ListAroundCache(ctx, params)
	if err != nil {
		return nil, err
	}
	return api.WrapPageResult(params, utils.CopySlice[*po.Product, *vo.ProductVO](list), count), nil
}
