package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/GoKit/future"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/bo"
	"mxclub/apps/mxclub-mini/entity/dto"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	miniUtil "mxclub/apps/mxclub-mini/utils"
	messageEnum "mxclub/domain/message/entity/enum"
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
	evaluationRepo orderRepo.IEvaluationRepo
	apRepo         repo.IAssistantApplicationRepo
	messageService *MessageService
}

func NewUserService(repo repo.IUserRepo,
	apRepo repo.IAssistantApplicationRepo,
	messageService *MessageService,
	orderRepo orderRepo.IOrderRepo,
	evaluationRepo orderRepo.IEvaluationRepo) *UserService {
	return &UserService{
		userRepo:       repo,
		apRepo:         apRepo,
		messageService: messageService,
		orderRepo:      orderRepo,
		evaluationRepo: evaluationRepo,
	}
}

func (svc UserService) GetUserById(ctx jet.Ctx, id uint) (*vo.UserVO, error) {
	// 用户信息
	userPO, err := svc.FindUserById(ctx, id)
	userVO := utils.MustCopyByCtx[vo.UserVO](ctx, userPO)
	if err != nil || userPO == nil {
		return nil, err
	}
	// 用户消费金额
	f1 := future.FutureFunc[float64](func() float64 {
		totalSpent, _ := svc.orderRepo.TotalSpent(ctx, id)
		return totalSpent
	})
	// 如果是打手，名字用打手名替换
	if userPO.Role == enum.RoleAssistant {
		userPO.WxName = fmt.Sprintf("%v 编号: %03d", userPO.Name, userPO.MemberNumber)
		// 获取打手评星
		staring, _ := svc.evaluationRepo.FindStaring(ctx, userPO.MemberNumber)
		userVO.DasherStaring = utils.RoundToTwoDecimalPlaces(staring)
	}
	totalSpent, _ := f1.Get()
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

func (svc UserService) FindUserById(ctx jet.Ctx, id uint) (*po.User, error) {
	return svc.userRepo.FindByIdAroundCache(ctx, id)
}

func (svc UserService) FindUserByDashId(ctx jet.Ctx, memberNumber int) (*po.User, error) {
	return svc.userRepo.FindByMemberNumber(ctx, memberNumber)
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
	userPOList, err := svc.userRepo.AssistantOnline(ctx)
	if err != nil {
		ctx.Logger().Errorf("[AssistantOnline]ERROR:%v", err.Error())
		return nil
	}
	filterUserList := utils.Filter(userPOList, func(in *po.User) bool {
		return in.ID != middleware.MustGetUserId(ctx)
	})
	return utils.Map[*po.User, *vo.AssistantOnlineVO](filterUserList, func(in *po.User) *vo.AssistantOnlineVO {
		return &vo.AssistantOnlineVO{
			Id:     in.MemberNumber,
			UserId: in.ID,
			Name:   in.Name,
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
		memberNumber := utils.ParseInt(handleReq.Ext)
		userPO, _ := svc.userRepo.FindByMemberNumber(ctx, memberNumber)
		message := fmt.Sprintf("您将被移除在进行中的订单，订单id:%v", handleReq.OrdersId)
		_ = svc.messageService.PushRemoveMessage(ctx, handleReq.OrdersId, userPO.ID, message)
	case 201:
		// 同意邀请
		svc.handleAcceptApplication(ctx, handleReq)
	case 301:
		// 接单拒绝，通知打手
		userPO, _ := svc.FindUserById(ctx, middleware.MustGetUserId(ctx))
		if handleReq.MessageType == messageEnum.REMOVE_MESSAGE {
			_ = svc.messageService.PushSystemMessage(ctx, userPO.ID,
				fmt.Sprintf("您移除打手:%v(%v)的申请已被拒绝，请联系相关打手", userPO.MemberNumber, userPO.Name))
		} else {
			_ = svc.messageService.PushSystemMessage(ctx, userPO.ID,
				fmt.Sprintf("您邀请打手:%v(%v)的申请已被拒绝，请联系其他打手", userPO.MemberNumber, userPO.Name))
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
	userPO, _ := svc.FindUserById(ctx, middleware.MustGetUserId(ctx))
	if orderPO.Executor2Id == -1 {
		_ = svc.messageService.PushSystemMessage(ctx, userPO.ID, fmt.Sprintf("您邀请打手:%v(%v)的申请已同意", userPO.MemberNumber, userPO.Name))
		// 更新角色
		_ = svc.orderRepo.UpdateOrderDasher2(ctx, orderId, userPO.MemberNumber, userPO.Name)
	} else if orderPO.Executor3Id == -1 {
		_ = svc.messageService.PushSystemMessage(ctx, userPO.ID, fmt.Sprintf("您邀请打手:%v(%v)的申请已同意", userPO.MemberNumber, userPO.Name))
		_ = svc.orderRepo.UpdateOrderDasher3(ctx, orderId, userPO.MemberNumber, userPO.Name)
	}
}

func (svc UserService) handleRemoveDasher(ctx jet.Ctx, handleReq *req.MessageHandleReq) {
	orderPO, _ := svc.orderRepo.FindByID(handleReq.OrdersId)
	userPO, _ := svc.FindUserById(ctx, middleware.MustGetUserId(ctx))
	if orderPO.Executor2Id == userPO.MemberNumber {
		_ = svc.orderRepo.UpdateOrderDasher2(ctx, orderPO.ID, -1, "")
	} else if orderPO.Executor3Id == userPO.MemberNumber {
		_ = svc.orderRepo.UpdateOrderDasher3(ctx, orderPO.ID, -1, "")
	}
	executorPO, _ := svc.userRepo.FindByMemberNumber(ctx, orderPO.ExecutorID)
	message := fmt.Sprintf("您移除打手:%v(%v)的申请已同意", userPO.MemberNumber, userPO.Name)
	_ = svc.messageService.PushSystemMessage(ctx, executorPO.ID, message)
}

func (svc UserService) checkUserGrade(ctx jet.Ctx, id uint) {
	spent, _ := svc.orderRepo.TotalSpent(ctx, id)
	gradeByScore := bo.GetGradeByScore(spent)
	m := map[string]any{
		"id":       id,
		"wx_grade": gradeByScore,
	}
	err := svc.userRepo.UpdateUser(ctx, m)

	if err != nil {
		ctx.Logger().Errorf("[checkUserGrade]ERROR:%v", err.Error())
	}
}

func (svc UserService) RemoveAssistant(ctx jet.Ctx) error {
	err := svc.userRepo.RemoveDasher(ctx, middleware.MustGetUserId(ctx))
	if err != nil {
		ctx.Logger().Errorf("[RemoveAssistant]ERROR:%v", err)
		return errors.New("注销打手失败")
	}
	return nil
}

func (svc UserService) PushInviteMessage(ctx jet.Ctx, req *req.OrderExecutorInviteReq) error {
	logger := ctx.Logger()
	// 检查是否在进行中订单
	orderPO, err := svc.orderRepo.FindByDasherId(ctx, req.ExecutorId)
	if err != nil || orderPO.ID > 0 {
		logger.Errorf("[FindByDasherId]ERROR:%v", err)
		return errors.New("打手有在进行中的订单，无法派单")
	}
	go func() {
		user, _ := svc.FindUserByDashId(ctx, req.ExecutorId)
		message := dto.NewDispatchMessage(user.ID, req.OrderId, req.GameRegion, req.RoleId, "")
		_ = svc.messageService.PushMessage(ctx, message)
	}()
	return nil
}
