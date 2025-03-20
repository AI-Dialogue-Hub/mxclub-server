package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/config"
	"mxclub/apps/mxclub-mini/entity/dto"
	"mxclub/apps/mxclub-mini/middleware"
	messageRepo "mxclub/domain/message/repo"
	notifyRepo "mxclub/domain/notify/repo"
	"mxclub/domain/user/repo"
	"mxclub/pkg/common/wxnotify"
	"mxclub/pkg/utils"
	"time"
)

func init() {
	jet.Provide(NewWxNotifyService)
}

type WxNotifyService struct {
	appId              string
	userRepo           repo.IUserRepo
	notifyTokenService wxnotify.INotifyTokenService
	messageRepo        messageRepo.IMessageRepo
	subNotifyRepo      notifyRepo.ISubNotifyRepo
}

func NewWxNotifyService(userRepo repo.IUserRepo,
	messageRepo messageRepo.IMessageRepo,
	subNotifyRepo notifyRepo.ISubNotifyRepo) *WxNotifyService {
	conf := config.GetConfig()
	notifyTokenService := wxnotify.NewNotifyTokenService(conf.WxConfig.Ak, conf.WxConfig.Sk)
	return &WxNotifyService{
		appId:              conf.WxConfig.Ak,
		notifyTokenService: notifyTokenService,
		userRepo:           userRepo,
		messageRepo:        messageRepo,
		subNotifyRepo:      subNotifyRepo,
	}
}

const (
	// MESSAGE_URI POST 小程序消息发送地址
	MESSAGE_URI = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%v"
)

const (
	TEMPLATE_COMMON = "5j8tGGE2F70S-G4gquPikNLk-DIpwTVz5Ophp2nbR5E"
)

func (s WxNotifyService) SendMessage(ctx jet.Ctx, userId uint, message string) error {
	// 1. 查找对应用户的openId
	userPO, err := s.userRepo.FindByIdAroundCache(ctx, userId)
	if err != nil || userPO == nil || userPO.ID <= 0 {
		ctx.Logger().Errorf("find user failed, userId:%v, err:%v", userId, err)
		return fmt.Errorf("cannot find userId %v", userId)
	}
	// 查找未读消息数量
	unReadMessageCount, err := s.messageRepo.CountUnReadMessageById(ctx, userId)
	if err != nil {
		ctx.Logger().Errorf("CountUnReadMessageById failed, userId:%v, err:%v", userId, err)
		unReadMessageCount = 1
	}
	openId := userPO.WxOpenId
	messageSendDTO := &dto.WxNotifyMessageSendDTO{
		Touser:     openId,
		TemplateID: TEMPLATE_COMMON,
		Page:       "mp.weixin.qq.com",
		Lang:       "zh_CN",
		Miniprogram: &dto.MiniProgram{
			Appid:    s.appId,
			Pagepath: "page/pages/my/comsumer",
		},
		Data: map[string]dto.DataValue{
			"time1":        {Value: time.Now().Format("2006-01-02 15:04:05")},
			"short_thing2": {Value: unReadMessageCount},
			"thing3":       {Value: message},
			"thing4":       {Value: "系统发送"},
		},
	}
	token, err := s.notifyTokenService.FetchToken()
	if err != nil || token == "" {
		ctx.Logger().Errorf("fetch token error:%v", err)
		return errors.New("token 获取失败")
	}
	// 2. 发送消息
	response, err := utils.Post(fmt.Sprintf(MESSAGE_URI, token), messageSendDTO)
	ctx.Logger().Infof("response,%v", utils.ObjToJsonStr(response))
	ctx.Logger().Infof("err, %v", err)
	return nil
}

func (s WxNotifyService) FindSubStatus(ctx jet.Ctx, templateId string) bool {
	var (
		userId = middleware.MustGetUserId(ctx)
	)
	return s.subNotifyRepo.ExistsSubNotifyRecord(ctx, userId, templateId)
}

func (s WxNotifyService) AddSubNotifyRecord(ctx jet.Ctx, templateId string) error {
	var (
		userId = middleware.MustGetUserId(ctx)
	)
	err := s.subNotifyRepo.AddSubNotifyRecord(ctx, userId, templateId)
	if err != nil {
		ctx.Logger().Errorf("AddSubNotifyRecord ERROR, %v", err)
		return errors.New("添加订阅状态失败")
	}
	return nil
}
