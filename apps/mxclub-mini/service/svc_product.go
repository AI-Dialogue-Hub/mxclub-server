package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/product/po"
	"mxclub/domain/product/repo"
	userRepo "mxclub/domain/user/repo"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewProductService)
}

type ProductService struct {
	ProductRepo repo.IProductRepo
	userRepo    userRepo.IUserRepo
}

func NewProductService(repo repo.IProductRepo, userRepo userRepo.IUserRepo) *ProductService {
	return &ProductService{ProductRepo: repo, userRepo: userRepo}
}

func (svc ProductService) FindById(ctx jet.Ctx, id uint) (*vo.ProductVO, error) {
	productPO, err := svc.ProductRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	productVO, _ := utils.Copy[vo.ProductVO](productPO)
	productVO.Description = productVO.ShortDescription + "\n" + productVO.Description
	userPO, _ := svc.userRepo.FindByIdAroundCache(ctx, middleware.MustGetUserId(ctx))
	productVO.Phone = userPO.Phone
	return productVO, nil
}

func (svc ProductService) List(ctx jet.Ctx, typeValue uint) ([]*vo.ProductVO, error) {
	var (
		list []*po.Product
		err  error
	)
	query := xmysql.NewMysqlQuery()
	query.SetPage(1, 1000)
	query.SetSort("created_at DESC")
	if typeValue == 101 {
		query.SetFilter("isHot = ?", true)
	} else if typeValue != 0 {
		query.SetFilter("type = ?", typeValue)
	}
	list, err = svc.ProductRepo.ListNoCountByQuery(query)
	if err != nil {
		return nil, err
	}
	productVOS := utils.CopySlice[*po.Product, *vo.ProductVO](list)
	// 老板已经保存电话了，选用上一次老板保存的电话
	userPO, _ := svc.userRepo.FindByIdAroundCache(ctx, middleware.MustGetUserId(ctx))
	utils.ForEach(productVOS, func(ele *vo.ProductVO) { ele.Phone = userPO.Phone })
	return productVOS, err
}
