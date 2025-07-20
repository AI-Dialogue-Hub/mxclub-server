package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/config"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/lottery/ability"
	lotteryPO "mxclub/domain/lottery/po"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/product/po"
	"mxclub/domain/product/repo"
	userPOInfo "mxclub/domain/user/po"
	userRepo "mxclub/domain/user/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/utils"
	"slices"
	"sort"
)

func init() {
	jet.Provide(NewProductService)
}

type ProductService struct {
	productRepo      repo.IProductRepo
	userRepo         userRepo.IUserRepo
	productSalesRepo repo.IProductSalesRepo
	lotteryActivity  ability.ILotteryAbility
}

func NewProductService(
	repo repo.IProductRepo,
	userRepo userRepo.IUserRepo,
	productSalesRepo repo.IProductSalesRepo,
	lotteryActivity ability.ILotteryAbility) *ProductService {
	return &ProductService{
		productRepo:      repo,
		userRepo:         userRepo,
		productSalesRepo: productSalesRepo,
		lotteryActivity:  lotteryActivity,
	}
}

func (svc ProductService) FindById(ctx jet.Ctx, id uint) (*vo.ProductVO, error) {
	// 查找产品信息
	productPO, err := svc.productRepo.FindByID(id)
	if err != nil {
		ctx.Logger().Errorf("cannot find product, productId is:%v", id)
		return nil, errors.New("查找商品出错或商品已下架，请联系客服")
	}

	// 复制 productPO 到 productVO
	productVO, err := utils.Copy[vo.ProductVO](productPO)
	if err != nil {
		ctx.Logger().Errorf("cannot copy product, productId is:%v", id)
		return nil, errors.New("查找商品出错或商品已下架，请联系客服")
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

	// 设置用户游戏Id
	productVO.GameId = userPO.GameId

	// 订单金额大于100才设置折扣
	if productVO.Price > 100 {
		// 查找用户会员优惠金额
		discountRate := enum.FetchDiscountByGrade(userPO.WxGrade)
		discountedPrice := utils.RoundToTwoDecimalPlaces(productVO.Price * discountRate)

		// 计算最终价格和优惠金额
		productVO.FinalPrice = discountedPrice
		productVO.DiscountPrice = utils.RoundToTwoDecimalPlaces(productVO.Price - productVO.FinalPrice)
	}

	return productVO, nil
}

func (svc ProductService) List(ctx jet.Ctx, typeValue uint) ([]*vo.ProductVO, error) {
	var (
		list   []*po.Product
		err    error
		userId = middleware.MustGetUserId(ctx)
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
	if err != nil || list == nil || len(list) == 0 {
		ctx.Logger().Errorf("list product failed, err:%v", err)
		return nil, errors.New("查询失败")
	}
	productVOS := utils.CopySlice[*po.Product, *vo.ProductVO](list)
	// 老板已经保存电话了，选用上一次老板保存的电话
	userPO, _ := svc.userRepo.FindByIdAroundCache(ctx, userId)
	// 老板保存的gameId
	utils.ForEach(productVOS, func(ele *vo.ProductVO) {
		ele.Phone = userPO.Phone
		ele.GameId = userPO.GameId
	})
	if config.GetConfig().WxPayConfig.IsBaoZaoClub() {
		sort.Slice(productVOS, func(i, j int) bool {
			return productVOS[i].FinalPrice < productVOS[j].FinalPrice
		})
	}
	// 查询商品的销量
	productIds := utils.Map[*po.Product, uint64](list, func(in *po.Product) uint64 { return in.ID })
	id2ProductMap, err := svc.productSalesRepo.FindByProductIds(ctx, productIds)
	if err == nil && id2ProductMap != nil && !id2ProductMap.IsEmpty() {
		utils.ForEach(productVOS, func(ele *vo.ProductVO) {
			if value, ok := id2ProductMap.Get(uint64(ele.ID)); ok {
				ele.SalesVolume = int(value.SalesVolume)
			}
		})
	}
	if !config.GetConfig().WxPayConfig.IsBaoZaoClub() {
		// 销量高的放在前面
		slices.SortFunc(productVOS, func(a, b *vo.ProductVO) int {
			if a == nil || b == nil {
				return 0
			}
			return b.SalesVolume - a.SalesVolume
		})
	}
	productVOS = svc.wrapLotteryPrize(ctx, typeValue, productVOS, userPO)
	return productVOS, nil
}

func (svc ProductService) wrapLotteryPrize(
	ctx jet.Ctx, typeValue uint, productVOS []*vo.ProductVO, userPO *userPOInfo.User,
) []*vo.ProductVO {
	if typeValue != 101 {
		return productVOS
	}
	// 抽奖活动
	lotteryActivityPOList, count, err := svc.lotteryActivity.ListActivity(
		ctx, &api.PageParams{Page: 1, PageSize: 1000})
	if err != nil {
		ctx.Logger().Errorf("list lottery activity failed, err:%v", err)
	} else {
		if count > 0 && lotteryActivityPOList != nil && len(lotteryActivityPOList) > 0 {
			lotteryProduct := utils.Map(lotteryActivityPOList,
				func(lotteryActivity *lotteryPO.LotteryActivity) *vo.ProductVO {
					return &vo.ProductVO{
						ID:               lotteryActivity.ID,
						Title:            lotteryActivity.ActivityTitle,
						Price:            lotteryActivity.ActivityPrice,
						DiscountRuleID:   0,
						DiscountPrice:    0,
						FinalPrice:       lotteryActivity.ActivityPrice,
						Description:      lotteryActivity.ActivityDesc,
						ShortDescription: lotteryActivity.ActivitySubtitle,
						Images:           utils.ToSlice(lotteryActivity.EntryImage),
						DetailImages:     utils.ToSlice(lotteryActivity.EntryImage),
						Thumbnail:        lotteryActivity.EntryImage,
						Phone:            userPO.Phone,
						GameId:           userPO.GameId,
						SalesVolume:      0,
						LotteryActivity:  true,
					}
				})
			productVOS = append(lotteryProduct, productVOS...)
		}
	}
	return productVOS
}
