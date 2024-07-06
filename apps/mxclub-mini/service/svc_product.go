package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewProductService)
}

type ProductService struct {
	ProductRepo repo.IProductRepo
}

func NewProductService(repo repo.IProductRepo) *ProductService {
	return &ProductService{ProductRepo: repo}
}

func (svc ProductService) FindById(id uint) (*vo.ProductVO, error) {
	po, err := svc.ProductRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return utils.Copy[vo.ProductVO](po)
}

func (svc ProductService) List(typeValue uint) ([]*vo.ProductVO, error) {
	list, err := svc.ProductRepo.ListNoCount(1, 1000, "type = ?", typeValue)
	if err != nil {
		return nil, err
	}
	productVOS := utils.CopySlice[*po.Product, *vo.ProductVO](list)
	return productVOS, err
}
