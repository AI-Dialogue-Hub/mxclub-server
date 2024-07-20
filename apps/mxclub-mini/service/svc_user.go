package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	miniUtil "mxclub/apps/mxclub-mini/utils"
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
	apRepo         repo.IAssistantApplicationRepo
	messageService *MessageService
}

func NewUserService(repo repo.IUserRepo, apRepo repo.IAssistantApplicationRepo, messageService *MessageService) *UserService {
	return &UserService{userRepo: repo, apRepo: apRepo, messageService: messageService}
}

func (svc UserService) GetUserById(ctx jet.Ctx, id uint) (*vo.User, error) {
	userPO, err := svc.userRepo.FindByID(id)
	return utils.MustCopyByCtx[vo.User](ctx, userPO), err
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

func (svc UserService) UpdateWxUserInfo(ctx jet.Ctx, userInfo *req.UserInfoReq) (*vo.User, error) {
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
	return utils.MustCopyByCtx[vo.User](ctx, userPO), nil
}

func (svc UserService) GetUserByOpenId(ctx jet.Ctx, openId string) (*vo.User, error) {
	userPO, err := svc.userRepo.FindByOpenId(ctx, openId)
	return utils.MustCopyByCtx[vo.User](ctx, userPO), err
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
	err := svc.userRepo.ToBeAssistant(ctx, application.UserID, application.Phone, application.MemberNumber)
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
	case 201:
		// 邀请打手操作，ext为订单id
	}
	return nil
}
