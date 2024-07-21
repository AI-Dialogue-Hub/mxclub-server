package service

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/domain/message/entity/enum"
	"mxclub/domain/message/repo"
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

func (svc MessageService) PushSystemMessage(ctx jet.Ctx, messageTo uint, content string) error {
	return svc.messageRepo.PushNormalMessage(ctx, enum.SYSTEM_NOTIFICATION, messageTo, "系统通知", content)
}
