package controller

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
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

func (ctl ProductController) GetV1ProductList(ctx jet.Ctx, req *req.ProductListReq) (*api.Response, error) {
	pageResult, err := ctl.productService.List(ctx, req)
	return xjet.WrapperResult(ctx, pageResult, err)
}

func (ctl ProductController) PostV1Product(ctx jet.Ctx, req *req.ProductReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "ok", ctl.productService.Update(ctx, req))
}

func (ctl ProductController) PostV1ProductHot(ctx jet.Ctx, req *req.ProductHotReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "ok", ctl.productService.UpdateHotInfo(ctx, req))
}

func (ctl ProductController) DeleteV1Product0(ctx jet.Ctx, param *api.PathParam) (*api.Response, error) {
	got, err := param.GetInt64(0)
	if err != nil {
		return xjet.WrapperResult(ctx, "FUCK YOU", err)
	}
	return xjet.WrapperResult(ctx, "ok", ctl.productService.DeleteById(ctx, got))
}

func (ctl ProductController) PutV1Product(ctx jet.Ctx, req *req.ProductReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "ok", ctl.productService.Add(ctx, req))
}

func (ctl ProductController) PostV1ProductSale(ctx jet.Ctx, req *req.ProductSaleReq) (*api.Response, error) {
	return xjet.WrapperResult(ctx, "ok", ctl.productService.UpdateSales(ctx, req))
}
