package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/domain/product/po"
	"mxclub/domain/product/repo"
	"mxclub/pkg/common/xmysql"
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
	productPO, err := svc.ProductRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return utils.Copy[vo.ProductVO](productPO)
}

func (svc ProductService) List(typeValue uint) ([]*vo.ProductVO, error) {
	var (
		list []*po.Product
		err  error
	)
	query := xmysql.NewMysqlQuery()
	query.SetPage(1, 1000)
	query.SetSort("created_at DESC")
	if typeValue == 101 {
		query.SetFilter("isHot = ?", true)
	} else {
		query.SetFilter("type = ?", typeValue)
	}
	list, err = svc.ProductRepo.ListNoCountByQuery(query)
	if err != nil {
		return nil, err
	}
	productVOS := utils.CopySlice[*po.Product, *vo.ProductVO](list)
	return productVOS, err
}
