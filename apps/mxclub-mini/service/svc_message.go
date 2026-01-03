package service

import (
	"errors"
	"mxclub/apps/mxclub-mini/entity/constant"
	"mxclub/apps/mxclub-mini/entity/dto"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/event"
	"mxclub/domain/message/entity/enum"
	"mxclub/domain/message/po"
	"mxclub/domain/message/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"

	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
)

func init() {
	jet.Provide(NewMessageService)
	jet.Invoke(func(svc *MessageService) {
		event.RegisterEvent("MessageService", event.EventRemoveDasher, svc.RemoveAllMessage)
	})
}

type MessageService struct {
	messageRepo repo.IMessageRepo
}

func NewMessageService(repo repo.IMessageRepo) *MessageService {
	return &MessageService{messageRepo: repo}
}

func (svc MessageService) List(ctx jet.Ctx, id uint, params *api.PageParams) (*api.PageResult, error) {
	messages, count, err := svc.messageRepo.ListAroundCache(ctx, params, id)
	if err != nil {
		return nil, errors.New("查询失败")
	}
	messageVOS := utils.CopySlice[*po.Message, *vo.MessageVO](messages)
	utils.ForEach(messageVOS, func(ele *vo.MessageVO) {
		ele.MessageDisPlayType = ele.MessageType.ParseDisPlayName()
	})
	return api.WrapPageResult(params, messageVOS, count), err
}

func (svc MessageService) ReadAllMessage(ctx jet.Ctx) error {
	userId := middleware.MustGetUserId(ctx)
	if err := svc.messageRepo.ReadAllMessage(ctx, userId); err != nil {
		return errors.New("标记已读失败")
	}
	return nil
}

func (svc MessageService) ReadByMessageId(ctx jet.Ctx, req *req.MessageReadReq) error {
	userId := middleware.MustGetUserId(ctx)
	// 接单的消息

	if err := svc.messageRepo.ReadByMessageId(ctx, userId, req.Id); err != nil {
		return errors.New("标记已读失败")
	}
	return nil
}

func (svc MessageService) CountUnReadMessage(ctx jet.Ctx) (int64, error) {
	count, err := svc.messageRepo.CountUnReadMessageById(ctx, middleware.MustGetUserId(ctx))
	if err != nil {
		return 0, errors.New("获取失败")
	}
	return count, nil
}

func (svc MessageService) PushMessage(ctx jet.Ctx, messageDTO *dto.MessageDTO) error {
	defer svc.messageRepo.ClearCache(ctx)
	return svc.messageRepo.InsertOne(&po.Message{
		MessageType:   messageDTO.MessageType,
		Title:         messageDTO.Title,
		Content:       messageDTO.Content,
		MessageFrom:   messageDTO.MessageFrom,
		MessageTo:     messageDTO.MessageTo,
		MessageStatus: messageDTO.MessageStatus,
		OrderId:       messageDTO.OrdersId,
		Ext:           messageDTO.Ext,
	})
}

func (svc MessageService) PushSystemMessage(ctx jet.Ctx, messageTo uint, content string) error {
	return svc.messageRepo.PushNormalMessage(ctx, enum.SYSTEM_NOTIFICATION, messageTo, "系统通知", content)
}

func (svc MessageService) PushLotteryMessage(ctx jet.Ctx, messageTo uint, content string) error {
	return svc.messageRepo.PushNormalMessage(ctx, enum.SYSTEM_NOTIFICATION, messageTo, "抽奖通知", content)
}

func (svc MessageService) PushRemoveMessage(ctx jet.Ctx, ordersId uint, messageTo uint, content string) error {
	return svc.messageRepo.PushOrderMessage(ctx, ordersId, enum.REMOVE_MESSAGE, messageTo, "系统通知", content)
}

func (svc MessageService) RemoveAllMessage(ctx jet.Ctx) error {
	// 如果存在打手Id 说明是后面打手进行注册然后清理的
	if _, exists := ctx.Get(constant.LOGOUT_DASHER_ID); exists {
		return nil
	}
	return svc.messageRepo.RemoveAllMessage(ctx, middleware.MustGetUserId(ctx))
}

func (svc MessageService) ClearDispatchMessage(orderId uint64, userId uint) error {
	return svc.messageRepo.RemoveMessageByOrderIdAndUserId(orderId, userId)
}
