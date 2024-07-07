package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/service"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xjet"
)

func init() {
	jet.Provide(NewProductController)
}

type ProductController struct {
	jet.BaseJetController
	productService *service.ProductService
}

func NewProductController(productService *service.ProductService) jet.ControllerResult {
	return jet.NewJetController(&ProductController{
		productService: productService,
	})
}

// =========================================================================

func (ctl ProductController) GetV1ProductList(ctx jet.Ctx, params *api.PageParams) (*api.Response, error) {
	pageResult, err := ctl.productService.List(ctx, params)
	return xjet.WrapperResult(ctx, pageResult, err)
}
