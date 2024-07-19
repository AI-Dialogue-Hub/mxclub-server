package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/apps/mxclub-mini/entity/dto"
	"mxclub/apps/mxclub-mini/entity/req"
	"mxclub/apps/mxclub-mini/entity/vo"
	"mxclub/apps/mxclub-mini/middleware"
	"mxclub/domain/message/po"
	"mxclub/domain/message/repo"
	"mxclub/pkg/api"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NeMessageService)
}

type MessageService struct {
	messageRepo repo.IMessageRepo
}

func NeMessageService(repo repo.IMessageRepo) *MessageService {
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
	return svc.messageRepo.InsertOne(&po.Message{
		MessageType:   messageDTO.MessageType,
		Title:         messageDTO.Title,
		Content:       messageDTO.Content,
		MessageFrom:   messageDTO.MessageFrom,
		MessageTo:     messageDTO.MessageTo,
		MessageStatus: messageDTO.MessageStatus,
		Ext:           messageDTO.Ext,
	})
}

func (svc MessageService) PushSystemMessage(ctx jet.Ctx, messageTo uint, content string) error {
	return svc.messageRepo.PushNormalMessage(ctx, messageTo, "系统通知", content)
}
