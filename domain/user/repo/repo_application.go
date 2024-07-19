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
	CreateAssistantApplication(ctx jet.Ctx, userID uint, phone string, memberNumber int64, name string) error
	// UpdateStatus 修改申请记录为通过
	UpdateStatus(ctx jet.Ctx, id uint, status string) error
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

func (repo AssistantApplicationRepo) CreateAssistantApplication(ctx jet.Ctx, userID uint, phone string, memberNumber int64, name string) error {
	application := &po.AssistantApplication{
		UserID:       userID,
		Phone:        phone,
		MemberNumber: memberNumber,
		Name:         name,
	}
	return repo.InsertOne(application)
}

func (repo AssistantApplicationRepo) UpdateStatus(ctx jet.Ctx, id uint, status string) error {
	update := xmysql.NewMysqlUpdate()
	update.SetFilter("id = ?", id)
	update.Set("status", status)
	return repo.UpdateByWrapper(update)
}
