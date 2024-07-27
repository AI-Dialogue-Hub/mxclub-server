package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	miniUtil "mxclub/apps/mxclub-mini/utils"
	messageEnum "mxclub/domain/message/entity/enum"
	orderEnum "mxclub/domain/order/entity/enum"
	orderRepo "mxclub/domain/order/repo"
	"mxclub/domain/user/entity/enum"
	"mxclub/domain/user/po"
	"mxclub/domain/user/repo"
	_ "mxclub/domain/user/repo"
	"mxclub/pkg/common/xjet"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewUserService)
}

type UserService struct {
	userRepo       repo.IUserRepo
	orderRepo      orderRepo.IOrderRepo
	apRepo         repo.IAssistantApplicationRepo
	messageService *MessageService
}

func NewUserService(repo repo.IUserRepo,
	apRepo repo.IAssistantApplicationRepo,
	messageService *MessageService,
	orderRepo orderRepo.IOrderRepo) *UserService {
	return &UserService{userRepo: repo, apRepo: apRepo, messageService: messageService, orderRepo: orderRepo}
}

func (svc UserService) GetUserById(ctx jet.Ctx, id uint) (*vo.UserVO, error) {
	// 用户信息
	userPO, err := svc.userRepo.FindByID(id)
	// 用户消费金额
	totalSpent, _ := svc.orderRepo.TotalSpent(ctx, id)
	userVO := utils.MustCopyByCtx[vo.UserVO](ctx, userPO)
	userVO.SetCurrentPoints(totalSpent)
	return userVO, err
}

func (svc UserService) WxLogin(ctx jet.Ctx, code string) (string, error) {
	openId, err := miniUtil.GetWxOpenId(code)
	if err != nil || xjet.IsAnyEmpty(openId) {
		ctx.Logger().Errorf("get user id err:%v", err)
		return "", err
	}
	var jwtToken string
	user, _ := svc.userRepo.FindByOpenId(ctx, openId)
	if user == nil || user.ID <= 0 {
		id, err := svc.userRepo.AddUserByOpenId(ctx, openId)
		if err != nil {
			ctx.Logger().Errorf("get user id err:%v", err)
			return "", errors.New("登录失败")
		}
		jwtToken, _ = middleware.GenAuthTokenByOpenIdAndUserId(id)
		return jwtToken, nil
	}
	jwtToken, _ = middleware.GenAuthTokenByOpenIdAndUserId(user.ID)
	return jwtToken, err
}

func (svc UserService) UpdateWxUserInfo(ctx jet.Ctx, userInfo *req.UserInfoReq) (*vo.UserVO, error) {
	var (
		userId   = middleware.MustGetUserId(ctx)
		imageURL string
		err      error
	)
	if userInfo != nil && userInfo.AvatarUrlBase64 != "" {
		imageURL, err = miniUtil.DecodeBase64ToImage(userInfo.AvatarUrlBase64)
	}
	if err != nil {
		imageURL = ""
	}
	err = svc.userRepo.UpdateUserIconAndNickName(ctx, userId, imageURL, "", "")
	if err != nil {
		ctx.Logger().Errorf("[UpdateWxUserInfo]ERROR:%v", err.Error())
		return nil, errors.New("用户信息更新失败")
	}
	userPO, _ := svc.GetUserById(ctx, userId)
	return utils.MustCopyByCtx[vo.UserVO](ctx, userPO), nil
}

func (svc UserService) GetUserByOpenId(ctx jet.Ctx, openId string) (*vo.UserVO, error) {
	userPO, err := svc.userRepo.FindByOpenId(ctx, openId)
	return utils.MustCopyByCtx[vo.UserVO](ctx, userPO), err
}

func (svc UserService) FindUserById(id uint) (*po.User, error) {
	return svc.userRepo.FindByID(id)
}

func (svc UserService) FindUserByDashId(memberNumber uint) (*po.User, error) {
	return svc.userRepo.FindByMemberNumber(memberNumber)
}

func (svc UserService) ToBeAssistant(ctx jet.Ctx, req req.AssistantReq) error {
	if svc.userRepo.ExistsAssistant(ctx, req.Phone, req.MemberNumber) {
		return errors.New("电话或id已被使用")
	}
	// 获取当前用户ID
	userID := middleware.MustGetUserId(ctx)

	// 创建打手申请记录
	err := svc.apRepo.CreateAssistantApplication(ctx, userID, req.Phone, req.MemberNumber, req.Name)
	if err != nil {
		ctx.Logger().Errorf("[ToBeAssistant]ERROR:%v", err.Error())
		return errors.New("提交打手申请失败，请联系客服")
	}
	// 发送申请消息
	err = svc.messageService.PushSystemMessage(ctx, userID, "您成为打手的申请已提交，请联系管理员审核")
	if err != nil {
		ctx.Logger().Errorf("[ToBeAssistant]消息发送失败:%v", err.Error())
	}
	return nil
}

func (svc UserService) PassAssistantApplication(ctx jet.Ctx, id uint) error {
	application, _ := svc.apRepo.FindByID(id)
	// 提交申请
	err := svc.userRepo.ToBeAssistant(ctx, application.UserID, application.Phone, application.MemberNumber, application.Name)
	if err != nil {
		ctx.Logger().Errorf("[ToBeAssistant]ERROR:%v", err.Error())
		return errors.New("转换身份失败，请联系客服")
	}
	return nil
}

func (svc UserService) AssistantOnline(ctx jet.Ctx) []*vo.AssistantOnlineVO {
	userPO, err := svc.userRepo.AssistantOnline(ctx)
	if err != nil {
		ctx.Logger().Errorf("[AssistantOnline]ERROR:%v", err.Error())
		return nil
	}
	return utils.Map[*po.User, *vo.AssistantOnlineVO](userPO, func(in *po.User) *vo.AssistantOnlineVO {
		return &vo.AssistantOnlineVO{
			Id:   in.MemberNumber,
			Name: in.Name,
		}
	})
}

func (svc UserService) CheckAssistantStatus(ctx jet.Ctx, memberNumber int) bool {
	return svc.userRepo.CheckAssistantStatus(ctx, memberNumber)
}

func (svc UserService) SwitchAssistantStatus(ctx jet.Ctx, status enum.MemberStatus) error {
	err := svc.userRepo.UpdateAssistantStatus(ctx, middleware.MustGetUserId(ctx), status)
	if err != nil {
		ctx.Logger().Errorf("[SwitchAssistantStatus]ERROR:%v", err.Error())
		return errors.New("修改在线状态失败，请联系客服")
	}
	return nil
}

func (svc UserService) AssistantStatus(ctx jet.Ctx) string {
	userPO, _ := svc.userRepo.FindByID(middleware.MustGetUserId(ctx))
	return string(userPO.MemberStatus)
}

func (svc UserService) HandleMessage(ctx jet.Ctx, handleReq *req.MessageHandleReq) error {
	switch handleReq.MessageTypeNumber {
	case 101:
		// 订单进行中 移除队友操作 ext为打手编号
		memberNumber := utils.ParseUint(handleReq.Ext)
		userPO, _ := svc.userRepo.FindByMemberNumber(memberNumber)
		message := fmt.Sprintf("您将被移除在进行中的订单，订单id:%v", handleReq.OrdersId)
		_ = svc.messageService.PushRemoveMessage(ctx, handleReq.OrdersId, userPO.ID, message)
	case 201:
		// 同意邀请
		svc.handleAcceptApplication(ctx, handleReq)
	case 301:
		// 接单拒绝，通知打手
		userPO, _ := svc.FindUserById(middleware.MustGetUserId(ctx))
		if handleReq.MessageType == messageEnum.REMOVE_MESSAGE {
			_ = svc.messageService.PushSystemMessage(ctx, userPO.ID, fmt.Sprintf("您移除打手:%v(%v)的申请已被拒绝，请联系相关打手", userPO.MemberNumber, userPO.Name))
		} else {
			_ = svc.messageService.PushSystemMessage(ctx, userPO.ID, fmt.Sprintf("您邀请打手:%v(%v)的申请已被拒绝，请联系其他打手", userPO.MemberNumber, userPO.Name))
		}
	case 401:
		// 同意移除
		svc.handleRemoveDasher(ctx, handleReq)
	}
	return nil
}

func (svc UserService) handleAcceptApplication(ctx jet.Ctx, handleReq *req.MessageHandleReq) {
	orderId := handleReq.OrdersId
	orderPO, _ := svc.orderRepo.FindByID(orderId)
	userPO, _ := svc.FindUserById(middleware.MustGetUserId(ctx))
	// 需要的打手
	needExecutorNum := utils.ParseUint(handleReq.Ext)
	if orderPO.Executor2Id == 0 {
		_ = svc.messageService.PushSystemMessage(ctx, userPO.ID, fmt.Sprintf("您邀请打手:%v(%v)的申请已同意", userPO.MemberNumber, userPO.Name))
		// 更新角色
		_ = svc.orderRepo.UpdateOrderDasher2(ctx, orderId, userPO.MemberNumber, userPO.Name)
		if needExecutorNum == 1 {
			// 开始订单
			_ = svc.orderRepo.UpdateOrderStatus(ctx, orderPO.OrderId, orderEnum.RUNNING)
		}
	} else if orderPO.Executor3Id == 0 {
		_ = svc.messageService.PushSystemMessage(ctx, userPO.ID, fmt.Sprintf("您邀请打手:%v(%v)的申请已同意", userPO.MemberNumber, userPO.Name))
		_ = svc.orderRepo.UpdateOrderDasher3(ctx, orderId, userPO.MemberNumber, userPO.Name)
		_ = svc.orderRepo.UpdateOrderStatus(ctx, orderPO.OrderId, orderEnum.RUNNING)
	}
}

func (svc UserService) handleRemoveDasher(ctx jet.Ctx, handleReq *req.MessageHandleReq) {
	orderPO, _ := svc.orderRepo.FindByID(handleReq.OrdersId)
	userPO, _ := svc.FindUserById(middleware.MustGetUserId(ctx))
	if orderPO.Executor2Id == userPO.MemberNumber {
		_ = svc.orderRepo.UpdateOrderDasher2(ctx, orderPO.ID, 0, "")
	} else if orderPO.Executor3Id == userPO.MemberNumber {
		_ = svc.orderRepo.UpdateOrderDasher3(ctx, orderPO.ID, 0, "")
	}
	executorPO, _ := svc.userRepo.FindByMemberNumber(orderPO.ExecutorID)
	message := fmt.Sprintf("您移除打手:%v(%v)的申请已同意", userPO.MemberNumber, userPO.Name)
	_ = svc.messageService.PushSystemMessage(ctx, executorPO.ID, message)
}
