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
	ReadByMessageId(ctx jet.Ctx, messageTo uint, messageId uint) error
	CountUnReadMessageById(ctx jet.Ctx, id uint) (int64, error)
	// PushNormalMessage messageTo 是
	PushNormalMessage(ctx jet.Ctx, messageType enum.MessageType, messageTo uint, title, content string) error
	PushOrderMessage(ctx jet.Ctx, ordersId uint, messageType enum.MessageType, messageTo uint, title, content string) error
	ClearCache(ctx jet.Ctx)
}

func NewMessageRepo(db *gorm.DB) IMessageRepo {
	repo := new(MessageRepo)
	repo.SetDB(db)
	repo.ModelPO = new(po.Message)
	repo.Ctx = context.Background()
	return repo
}

type MessageRepo struct {
	xmysql.BaseRepo[po.Message]
}

// ===========================================================

const cachePrefix = "message_product"
const listCachePrefix = cachePrefix + "_list"

func (repo MessageRepo) ListAroundCache(ctx jet.Ctx, params *api.PageParams, id uint) ([]*po.Message, int64, error) {
	parseIdString := utils.ParseString(id)
	// 根据页码参数生成唯一的缓存键
	cacheListKey := xredis.BuildListDataCacheKey(cachePrefix+parseIdString, params)
	cacheCountKey := xredis.BuildListCountCacheKey(listCachePrefix + parseIdString)

	list, count, err := xredis.GetListOrDefault[po.Message](ctx, cacheListKey, cacheCountKey, func() ([]*po.Message, int64, error) {
		// 如果缓存中未找到，则从数据库中获取
		list, count, err := repo.ListAndOrder(params.Page, params.PageSize, "created_at DESC", "message_to = ?", id)
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

func (repo MessageRepo) DeleteById(ctx jet.Ctx, id int64) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	err := repo.RemoveByID(id)
	if err != nil {
		ctx.Logger().Errorf("[productRepo]DeleteById ERROR:%v", err.Error())
		return errors.New("删除失败")
	}
	return nil
}

func (repo MessageRepo) Add(ctx jet.Ctx, po *po.Message) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	err := repo.InsertOne(po)
	if err != nil {
		ctx.Logger().Errorf("[productRepo]AddDeduction ERROR:%v", err.Error())
		return errors.New("添加失败")
	}
	return nil
}

func (repo MessageRepo) ReadAllMessage(ctx jet.Ctx, id uint) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	updateMap := map[string]any{"message_status": 1}
	err := repo.Update(updateMap, "message_to", id)
	if err != nil {
		ctx.Logger().Errorf("[ReadAllMessage]ERROR:%v", err.Error())
		return err
	}
	return nil
}

func (repo MessageRepo) CountUnReadMessageById(ctx jet.Ctx, id uint) (int64, error) {
	count, err := repo.Count("message_to = ? and message_status = ?", id, enum.UN_READ)
	if err != nil {
		ctx.Logger().Errorf("[CountUnReadMessageById]ERROR:%v", err.Error())
		return 0, err
	}
	return count, nil
}

func (repo MessageRepo) ReadByMessageId(ctx jet.Ctx, messageTo uint, messageId uint) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	updateMap := map[string]any{"message_status": 1}
	err := repo.Update(updateMap, "id = ? and message_to = ?", messageId, messageTo)
	if err != nil {
		ctx.Logger().Errorf("[ReadAllMessage]ERROR:%v", err.Error())
		return err
	}
	return nil
}

func (repo MessageRepo) PushNormalMessage(ctx jet.Ctx, messageType enum.MessageType, messageTo uint, title, content string) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	messagePO := &po.Message{
		MessageType:   messageType,
		Title:         title,
		Content:       content,
		MessageTo:     messageTo,
		MessageStatus: enum.UN_READ,
	}
	return repo.InsertOne(messagePO)
}

func (repo MessageRepo) PushOrderMessage(ctx jet.Ctx, ordersId uint, messageType enum.MessageType, messageTo uint, title, content string) error {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
	messagePO := &po.Message{
		MessageType:   messageType,
		Title:         title,
		OrderId:       ordersId,
		Content:       content,
		MessageTo:     messageTo,
		MessageStatus: enum.UN_READ,
	}
	return repo.InsertOne(messagePO)
}

func (repo MessageRepo) ClearCache(ctx jet.Ctx) {
	_ = xredis.DelMatchingKeys(ctx, cachePrefix)
}
