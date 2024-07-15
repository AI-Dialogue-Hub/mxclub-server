package controller

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewProductController)
}

type ProductController struct {
	jet.BaseJetController
	ProductService *service.ProductService
}

func NewProductController(ProductService *service.ProductService) jet.ControllerResult {
	return jet.NewJetController(&ProductController{
		ProductService: ProductService,
	})
}

// ============================================================

func (ctr ProductController) GetV1Product0(ctx jet.Ctx, args *api.PathParam) (*api.Response, error) {
	productId := args.CmdArgs[0]
	if xjet.IsAnyEmpty(productId) {
		return xjet.WrapperResult(ctx, "", errors.New("product is empty"))
	}
	vo, err := ctr.ProductService.FindById(utils.ParseUint(productId))
	return xjet.WrapperResult(ctx, vo, err)
}

func (ctr ProductController) GetV1ProductType0(ctx jet.Ctx, args *api.PathParam) (*api.Response, error) {
	var (
		typeValue = "0"
	)
	if args != nil && len(args.CmdArgs) > 0 {
		typeValue = args.CmdArgs[0]
	}
	if xjet.IsAnyEmpty(typeValue) {
		return xjet.WrapperResult(ctx, "", errors.New("typeValue is empty"))
	}
	vo, err := ctr.ProductService.List(utils.ParseUint(typeValue))
	return xjet.WrapperResult(ctx, vo, err)
}
