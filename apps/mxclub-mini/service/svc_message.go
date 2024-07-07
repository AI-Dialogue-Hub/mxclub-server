package service

import (
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
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
	userId := middleware.MustGetUserInfo(ctx)
	if err := svc.messageRepo.ReadAllMessage(ctx, userId); err != nil {
		return errors.New("标记已读失败")
	}
	return nil
}

func (svc MessageService) CountUnReadMessage(ctx jet.Ctx) (int64, error) {
	count, err := svc.messageRepo.CountUnReadMessageById(ctx, middleware.MustGetUserInfo(ctx))
	if err != nil {
		return 0, errors.New("获取失败")
	}
	return count, nil
}
