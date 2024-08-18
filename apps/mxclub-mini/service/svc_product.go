package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/order/entity/enum"
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
	productRepo repo.IProductRepo
	userRepo    userRepo.IUserRepo
}

func NewProductService(repo repo.IProductRepo, userRepo userRepo.IUserRepo) *ProductService {
	return &ProductService{productRepo: repo, userRepo: userRepo}
}

func (svc ProductService) FindById(ctx jet.Ctx, id uint) (*vo.ProductVO, error) {
	// 查找产品信息
	productPO, err := svc.productRepo.FindByID(id)
	if err != nil {
		ctx.Logger().Errorf("cannot find product, productId is:%v", id)
		return nil, errors.New("查找商品出错，请联系客服")
	}

	// 复制 productPO 到 productVO
	productVO, err := utils.Copy[vo.ProductVO](productPO)
	if err != nil {
		ctx.Logger().Errorf("cannot copy product, productId is:%v", id)
		return nil, errors.New("查找商品出错，请联系客服")
	}

	// 拼接 Description 字段
	productVO.Description = productVO.ShortDescription + "\n" + productVO.Description

	// 查找用户信息
	userPO, err := svc.userRepo.FindByIdAroundCache(ctx, middleware.MustGetUserId(ctx))
	if err != nil {
		return nil, err
	}

	// 设置用户手机号
	productVO.Phone = userPO.Phone

	// 查找用户会员优惠金额
	discountRate := enum.FetchDiscountByGrade(userPO.WxGrade)
	discountedPrice := utils.RoundToTwoDecimalPlaces(productVO.Price * discountRate)

	// 计算最终价格和优惠金额
	productVO.FinalPrice = discountedPrice
	productVO.DiscountPrice = utils.RoundToTwoDecimalPlaces(productVO.Price - productVO.FinalPrice)

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
	list, err = svc.productRepo.ListNoCountByQuery(query)
	if err != nil {
		return nil, err
	}
	productVOS := utils.CopySlice[*po.Product, *vo.ProductVO](list)
	// 老板已经保存电话了，选用上一次老板保存的电话
	userPO, _ := svc.userRepo.FindByIdAroundCache(ctx, middleware.MustGetUserId(ctx))
	utils.ForEach(productVOS, func(ele *vo.ProductVO) { ele.Phone = userPO.Phone })
	return productVOS, err
}
