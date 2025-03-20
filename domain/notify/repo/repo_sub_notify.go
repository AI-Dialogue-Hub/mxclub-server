package repo

import (
	"context"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/notify/po"
	"mxclub/pkg/common/xmysql"
)

func init() {
	jet.Provide(NewSubNotifyRepo)
}

type ISubNotifyRepo interface {
	xmysql.IBaseRepo[po.SubNotifyRecord]
	AddSubNotifyRecord(ctx jet.Ctx, userId uint, templateId string) error
	ExistsSubNotifyRecord(ctx jet.Ctx, userId uint, templateId string) bool
}

func NewSubNotifyRepo(db *gorm.DB) ISubNotifyRepo {
	repo := new(SubNotifyRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.SubNotifyRecord)
	repo.Ctx = context.Background()
	return repo
}

type SubNotifyRepo struct {
	xmysql.BaseRepo[po.SubNotifyRecord]
}

func (repo SubNotifyRepo) AddSubNotifyRecord(ctx jet.Ctx, userId uint, templateId string) error {
	return repo.InsertOne(&po.SubNotifyRecord{
		UserID:     userId,
		TemplateID: templateId,
	})
}

func (repo SubNotifyRepo) ExistsSubNotifyRecord(ctx jet.Ctx, userId uint, templateId string) bool {
	count, err := repo.Count("user_id = ? and template_id = ?", userId, templateId)
	if err != nil {
		ctx.Logger().Errorf("ExistsSubNotifyRecord ERROR, %v", err)
		return false
	}
	return count == 1
}
