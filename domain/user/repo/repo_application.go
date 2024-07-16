package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/user/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewAssistantApplicationRepo)
}

type IAssistantApplicationRepo interface {
	xmysql.IBaseRepo[po.AssistantApplication]
	CreateAssistantApplication(ctx jet.Ctx, userID uint, phone string, memberNumber int64) error
}

func NewAssistantApplicationRepo(db *gorm.DB) IAssistantApplicationRepo {
	repo := new(AssistantApplicationRepo)
	repo.Db = db
	repo.ModelPO = new(po.AssistantApplication)
	repo.Ctx = context.Background()
	return repo
}

type AssistantApplicationRepo struct {
	xmysql.BaseRepo[po.AssistantApplication]
}

func (repo *AssistantApplicationRepo) CreateAssistantApplication(ctx jet.Ctx, userID uint, phone string, memberNumber int64) error {
	application := &po.AssistantApplication{
		UserID:       userID,
		Phone:        phone,
		MemberNumber: memberNumber,
	}
	return repo.InsertOne(application)
}
