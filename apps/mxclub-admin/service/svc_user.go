package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-admin/entity/req"
	"mxclub/apps/mxclub-admin/entity/vo"
	"mxclub/apps/mxclub-admin/middleware"
	"mxclub/domain/event"
	orderRepo "mxclub/domain/order/repo"
	"mxclub/domain/user/entity/enum"
	"mxclub/domain/user/po"
	"mxclub/domain/user/repo"
	_ "mxclub/domain/user/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewUserService)
}

type UserService struct {
	userRepo       repo.IUserRepo
	orderRepo      orderRepo.IOrderRepo
	messageService *MessageService
}

func NewUserService(repo repo.IUserRepo, orderRepo orderRepo.IOrderRepo, messageService *MessageService) *UserService {
	return &UserService{userRepo: repo, orderRepo: orderRepo, messageService: messageService}
}

// =============================================================

func (svc UserService) GetUserById(ctx jet.Ctx, id int) (*vo.UserLoginVO, error) {
	userPO, err := svc.userRepo.FindByID(id)
	return utils.MustCopyByCtx[vo.UserLoginVO](ctx, userPO), err
}

func (svc UserService) CheckUser(ctx jet.Ctx, username string, password string) (*po.User, error) {
	userPO, err := svc.userRepo.QueryUserByAccount(username, password)
	if err != nil || userPO == nil || userPO.ID == 0 {
		ctx.Logger().Infof("user %s not exist", username)
		return nil, errors.New("账号或密码错误")
	}
	return userPO, nil
}

func (svc UserService) List(ctx jet.Ctx, params *req.UserListReq) (*api.PageResult, error) {
	if params.UserType == "all" {
		params.UserType = ""
	}
	list, count, err := svc.userRepo.ListAroundCacheByUserTypeAndDasherId(
		ctx, params.PageParams, enum.RoleType(params.UserType), params.MemberNumber)

	if err != nil {
		ctx.Logger().Errorf("[UserService List] error:%v", err.Error())
		return nil, errors.New("查询失败")
	}
	userVOS := utils.CopySlice[*po.User, *vo.UserVO](list)
	utils.ForEach(userVOS, func(t *vo.UserVO) { t.DisPlayName = t.Role.DisPlayName() })
	return api.WrapPageResult(params.PageParams, userVOS, count), nil
}

func (svc UserService) Update(ctx jet.Ctx, req *req.UserReq) error {
	updateMap := utils.ObjToMap(req)
	updateMap["member_number"] = utils.SafeParseNumber[int](updateMap["member_number"])
	return svc.userRepo.UpdateUser(ctx, updateMap)
}

func (svc UserService) RemoveAssistant(ctx jet.Ctx, userId uint) error {
	// 最后再注销打手信息
	defer svc.RemoveAssistantEvent(ctx)
	ctx.Put("userId", userId)
	event.PublishEvent(event.EventRemoveDasher, ctx)
	userPO, _ := svc.userRepo.FindByIdAroundCache(ctx, userId)
	ctx.Logger().Infof("remove dasher, dasherId:%v, dasherName: %v", userPO.MemberNumber, utils.ObjToJsonStr(userPO))
	return nil
}

func (svc UserService) RemoveAssistantEvent(ctx jet.Ctx) error {
	userId := ctx.MustGet("userId").(uint)
	return svc.userRepo.RemoveDasher(ctx, userId)
}

func (svc UserService) AssistantOnline(ctx jet.Ctx) []*vo.AssistantOnlineVO {
	userPOList, err := svc.userRepo.AssistantOnline(ctx)
	if err != nil {
		ctx.Logger().Errorf("[AssistantOnline]ERROR:%v", err.Error())
		return nil
	}
	filterUserList := utils.Filter(userPOList, func(in *po.User) bool {
		return in.ID != middleware.MustGetUserInfo(ctx).ID
	})
	return utils.Map[*po.User, *vo.AssistantOnlineVO](filterUserList, func(in *po.User) *vo.AssistantOnlineVO {
		return &vo.AssistantOnlineVO{
			Id:     in.MemberNumber,
			UserId: in.ID,
			Name:   in.Name,
		}
	})
}

// RemoveDasher 移除队友
func (svc UserService) RemoveDasher(ctx jet.Ctx, req *req.UserRemoveReq) error {
	orderPO, err := svc.orderRepo.FindByOrderOrOrdersId(ctx, req.OrderId)
	if err != nil {
		return errors.New(fmt.Sprintf("订单不存在，订单Id:%v", req.OrderId))
	}
	for _, dasherId := range req.DasherIds {
		if orderPO.Executor2Id == dasherId {
			err = svc.orderRepo.RemoveDasher(ctx, req.OrderId, 2)
			userPO, _ := svc.userRepo.FindByMemberNumber(ctx, dasherId)
			_ = svc.messageService.PushSystemMessage(ctx, userPO.ID, fmt.Sprintf("您已被管理员移出订单:%v", req.OrderId))
			if err != nil {
				ctx.Logger().Errorf("[UserService#RemoveDasher] remove failed, err:%v", err)
				return errors.New("移除失败")
			}
		} else if orderPO.Executor3Id == dasherId {
			err = svc.orderRepo.RemoveDasher(ctx, req.OrderId, 3)
			userPO, _ := svc.userRepo.FindByMemberNumber(ctx, dasherId)
			_ = svc.messageService.PushSystemMessage(ctx, userPO.ID, fmt.Sprintf("您已被管理员移出订单:%v", req.OrderId))
			if err != nil {
				ctx.Logger().Errorf("[UserService#RemoveDasher] remove failed, err:%v", err)
				return errors.New("移除失败")
			}
		}
	}
	return nil
}
