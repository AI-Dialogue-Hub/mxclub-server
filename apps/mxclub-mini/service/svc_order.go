package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewOrderService)
}

type OrderService struct {
	orderRepo      repo.IOrderRepo
	withdrawalRepo repo.IWithdrawalRepo
	userService    *UserService
}

func NewOrderService(repo repo.IOrderRepo, withdrawalRepo repo.IWithdrawalRepo, userService *UserService) *OrderService {
	return &OrderService{
		orderRepo:      repo,
		withdrawalRepo: withdrawalRepo,
		userService:    userService,
	}
}

// ===============================================================

func (svc OrderService) Add(ctx jet.Ctx, req *req.OrderReq) error {
	return nil
}

func (svc OrderService) List(ctx jet.Ctx, req *req.OrderListReq) (*api.PageResult, error) {
	list, err := svc.orderRepo.ListByOrderStatus(ctx, req.OrderStatus, &req.PageParams, req.Ge, req.Le)
	if err != nil {
		ctx.Logger().Errorf("[OrderService]List ERROR, %v", err.Error())
		return nil, errors.New("查询不到数据")
	}
	return api.WrapPageResult(&req.PageParams, utils.CopySlice[*po.Order, *vo.OrderVO](list), 0), err
}

func (svc OrderService) HistoryWithDrawAmount(ctx jet.Ctx) (*vo.WithDrawVO, error) {
	userPO, err := svc.userService.FindUserById(middleware.MustGetUserInfo(ctx))
	if err != nil {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, cannot find user:%v", middleware.MustGetUserInfo(ctx))
		return nil, errors.New("cannot find user info")
	}
	withdrawnAmount, _ := svc.withdrawalRepo.WithdrawnAmount(ctx, userPO.MemberNumber)
	orderWithdrawAbleAmount, _ := svc.orderRepo.OrderWithdrawAbleAmount(ctx, userPO.MemberNumber)
	if withdrawnAmount > orderWithdrawAbleAmount {
		ctx.Logger().Errorf("[HistoryWithDrawAmount]ERROR, withdrawnAmount: %v gt orderWithdrawAbleAmount:%v", withdrawnAmount, orderWithdrawAbleAmount)
		return nil, errors.New("系统查询错误，请联系管理员")
	}
	return &vo.WithDrawVO{
		HistoryWithDrawAmount: withdrawnAmount,
		WithdrawAbleAmount:    orderWithdrawAbleAmount - withdrawnAmount,
		WithdrawRangeMax:      2000,
		WithdrawRangeMin:      200,
	}, nil
}
