package repo

import (
	"context"
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
	"mxclub/domain/message/entity/enum"
	"mxclub/domain/message/po"
	"mxclub/pkg/api"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/common/xredis"
	"mxclub/pkg/utils"
)

func init() {
	jet.Provide(NewMessageRepo)
}

type IMessageRepo interface {
	xmysql.IBaseRepo[po.Message]
	ListAroundCache(ctx jet.Ctx, params *api.PageParams, id uint) ([]*po.Message, int64, error)
	DeleteById(ctx jet.Ctx, id int64) error
	Add(ctx jet.Ctx, po *po.Message) error
	ReadAllMessage(ctx jet.Ctx, id uint) error
	CountUnReadMessageById(ctx jet.Ctx, id uint) (int64, error)
}

func NewMessageRepo(db *gorm.DB) IMessageRepo {
	repo := new(MessageRepo)
	repo.Db = db.Model(new(po.Message))
	repo.Ctx = context.Background()
	return repo
}

type MessageRepo struct {
	xmysql.BaseRepo[po.Message]
}

// ===========================================================

const cachePrefix = "message_product"
const listCachePrefix = cachePrefix + "_list"

func (repo *MessageRepo) ListAroundCache(ctx jet.Ctx, params *api.PageParams, id uint) ([]*po.Message, int64, error) {
	parseIdString := utils.ParseString(id)
	// 根据页码参数生成唯一的缓存键
	cacheListKey := xredis.BuildListDataCacheKey(cachePrefix+parseIdString, params)
	cacheCountKey := xredis.BuildListCountCacheKey(listCachePrefix + parseIdString)

	list, count, err := xredis.GetListOrDefault[po.Message](ctx, cacheListKey, cacheCountKey, func() ([]*po.Message, int64, error) {
		// 如果缓存中未找到，则从数据库中获取
		list, count, err := repo.ListAndOrder(params.Page, params.PageSize, "created_at DESC", nil)
		if err != nil {
			return nil, 0, err
		}
		return list, count, nil
	})
	if err != nil {
		ctx.Logger().Errorf("ListAroundCache 错误: %v", err)
		return nil, 0, err
	}

	return list, count, nil
}

func (repo *MessageRepo) DeleteById(ctx jet.Ctx, id int64) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	err := repo.RemoveByID(id)
	if err != nil {
		ctx.Logger().Errorf("[ProductRepo]DeleteById ERROR:%v", err.Error())
		return errors.New("删除失败")
	}
	return nil
}

func (repo *MessageRepo) Add(ctx jet.Ctx, po *po.Message) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	err := repo.InsertOne(po)
	if err != nil {
		ctx.Logger().Errorf("[ProductRepo]Add ERROR:%v", err.Error())
		return errors.New("添加失败")
	}
	return nil
}

func (repo *MessageRepo) ReadAllMessage(ctx jet.Ctx, id uint) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	updateMap := map[string]any{"message_status": 1}
	err := repo.Update(updateMap, "message_to", id)
	if err != nil {
		ctx.Logger().Errorf("[ReadAllMessage]ERROR:%v", err.Error())
		return err
	}
	return nil
}

func (repo *MessageRepo) CountUnReadMessageById(ctx jet.Ctx, id uint) (int64, error) {
	count, err := repo.Count("message_to = ? and message_status = ?", id, enum.UN_READ)
	if err != nil {
		ctx.Logger().Errorf("[CountUnReadMessageById]ERROR:%v", err.Error())
		return 0, err
	}
	return count, nil
}