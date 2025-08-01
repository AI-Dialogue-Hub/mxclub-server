package xmysql

import (
	"context"
	"errors"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"gorm.io/gorm"
)

type IBaseRepo[T any] interface {
	DB() *gorm.DB
	InsertOne(t *T) error
	InsertBatch(entities []any) (int, error)
	InsertMany(entities []*T) (int, error)
	RemoveByID(id interface{}) error
	Remove(filter any, data ...any) error
	Update(update any, filter any, data ...any) error
	UpdateById(update any, id any) error
	UpdateByWrapper(updateWrap *MysqlUpdate) error
	FindByID(id interface{}) (*T, error)
	Find(filter any, data ...any) ([]*T, error)
	FindOne(filter any, data ...any) (*T, error)
	FindOneByWrapper(query *MysqlQuery) (*T, error)
	FindByWrapper(query *MysqlQuery) ([]*T, error)
	FindAll() ([]*T, error)
	FindOrCreate(findFunc func() bool, t *T) (*T, error)
	List(pageNo int64, pageSize int64, filter any, data ...any) ([]*T, int64, error)
	ListByWrapper(ctx jet.Ctx, query *MysqlQuery) ([]*T, int64, error)
	ListAndOrder(pageNo int64, pageSize int64, order string, filter any, data ...any) ([]*T, int64, error)
	ListNoCount(pageNo int64, pageSize int64, order string, filter any, data ...any) ([]*T, error)
	ListNoCountByQuery(query *MysqlQuery) ([]*T, error)
	Count(filter any, data ...any) (int64, error)
}

type BaseRepo[T any] struct {
	db      *gorm.DB
	Ctx     context.Context
	ModelPO *T
}

func (r *BaseRepo[T]) DB() *gorm.DB {
	return r.db.Model(r.ModelPO)
}

func (r *BaseRepo[T]) SetDB(db *gorm.DB) {
	r.db = db
}

func (r *BaseRepo[T]) InsertOne(t *T) error {
	return r.db.Model(r.ModelPO).WithContext(r.Ctx).Create(t).Error
}

func (r *BaseRepo[T]) InsertBatch(entities []interface{}) (int, error) {
	if len(entities) == 0 {
		return 0, errors.New("cannot insert empty array")
	}
	result := r.db.Model(r.ModelPO).WithContext(r.Ctx).Create(entities)
	return int(result.RowsAffected), result.Error
}

func (r *BaseRepo[T]) InsertMany(entities []*T) (int, error) {
	if len(entities) == 0 {
		return 0, errors.New("cannot insert empty array")
	}
	result := r.db.Model(r.ModelPO).WithContext(r.Ctx).Create(entities)
	return int(result.RowsAffected), result.Error
}

func (r *BaseRepo[T]) RemoveByID(id interface{}) error {
	return r.db.Model(r.ModelPO).WithContext(r.Ctx).Delete(new(T), id).Error
}

func (r *BaseRepo[T]) Remove(filter any, data ...any) error {
	return r.db.Model(r.ModelPO).WithContext(r.Ctx).Where(filter, data...).Delete(new(T)).Error
}

func (r *BaseRepo[T]) Update(update any, filter any, data ...any) error {
	return r.db.Model(r.ModelPO).WithContext(r.Ctx).Where(filter, data...).Updates(update).Error
}

func (r *BaseRepo[T]) UpdateById(update any, id any) error {
	return r.Update(update, "id = ?", id)
}

func (r *BaseRepo[T]) UpdateByWrapper(updateWrap *MysqlUpdate) error {
	return r.db.Model(r.ModelPO).Where(updateWrap.Query, updateWrap.Args...).Updates(updateWrap.Values).Error
}

func (r *BaseRepo[T]) FindByID(id interface{}) (*T, error) {
	var t T
	err := r.db.Model(r.ModelPO).WithContext(r.Ctx).Take(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *BaseRepo[T]) Find(filter any, data ...any) ([]*T, error) {
	var entities []*T
	err := r.db.Model(r.ModelPO).WithContext(r.Ctx).Where(filter, data...).Find(&entities).Error
	return entities, err
}

func (r *BaseRepo[T]) FindOneByWrapper(query *MysqlQuery) (*T, error) {
	var entity *T
	err := r.db.Model(r.ModelPO).WithContext(r.Ctx).Where(query.Query, query.Args...).Find(&entity).Error
	return entity, err
}

func (r *BaseRepo[T]) FindByWrapper(query *MysqlQuery) ([]*T, error) {
	var entityList = make([]*T, 0)
	err := r.db.Model(r.ModelPO).WithContext(r.Ctx).Where(query.Query, query.Args...).Find(&entityList).Error
	return entityList, err
}

func (r *BaseRepo[T]) FindOne(filter any, data ...any) (*T, error) {
	var entity T
	err := r.db.Model(r.ModelPO).WithContext(r.Ctx).Where(filter, data...).First(&entity).Error
	return &entity, err
}

func (r *BaseRepo[T]) FindAll() ([]*T, error) {
	var entities []*T
	err := r.db.Model(r.ModelPO).WithContext(r.Ctx).Find(&entities).Error
	return entities, err
}

func (r *BaseRepo[T]) FindOrCreate(findFunc func() bool, t *T) (*T, error) {
	if !findFunc() {
		err := r.InsertOne(t)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func (r *BaseRepo[T]) List(pageNo int64, pageSize int64, filter any, data ...any) ([]*T, int64, error) {
	return r.ListAndOrder(pageNo, pageSize, "", filter, data...)
}

func (r *BaseRepo[T]) ListByWrapper(ctx jet.Ctx, query *MysqlQuery) ([]*T, int64, error) {
	entities := make([]*T, 0)
	if query.Limit <= 0 {
		query.Limit = 10
	}
	err := r.db.Model(r.ModelPO).WithContext(r.Ctx).
		Where(query.Query, query.Args...).
		Offset(query.Offset).
		Order(query.Sort).
		Limit(query.Limit).
		Find(&entities).Error
	if err != nil {
		ctx.Logger().Errorf("ERROR:%v", err.Error())
		return nil, 0, err
	}
	count, _ := r.Count(query.Query, query.Args...)
	return entities, count, nil
}

func (r *BaseRepo[T]) ListAndOrder(pageNo int64, pageSize int64, order string, filter any, data ...any) ([]*T, int64, error) {
	entities, err := r.ListNoCount(pageNo, pageSize, order, filter, data...)
	if err != nil {
		return nil, 0, err
	}
	count, err := r.Count(filter, data...)
	return entities, count, err
}

func (r *BaseRepo[T]) ListNoCount(pageNo int64, pageSize int64, order string, filter any, data ...any) ([]*T, error) {
	if filter == nil {
		filter = map[string]interface{}{}
	}
	entities := make([]*T, 0)
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	var err error
	tx := r.db.Model(r.ModelPO).WithContext(r.Ctx).
		Where(filter, data...).
		Offset(int((pageNo - 1) * pageSize)).
		Limit(int(pageSize))
	if order == "" {
		err = tx.Find(&entities).Error
	} else {
		err = tx.Order(order).Find(&entities).Error
	}
	return entities, err
}

func (r *BaseRepo[T]) ListNoCountByQuery(query *MysqlQuery) ([]*T, error) {
	entities := make([]*T, 0)
	if query.Limit <= 0 {
		query.Limit = 1000
	}
	var err error
	err = r.db.Model(r.ModelPO).WithContext(r.Ctx).
		Where(query.Query, query.Args...).
		Where("deleted_at IS NULL").
		Offset(query.Offset).
		Order(query.Sort).
		Limit(query.Limit).
		Find(&entities).Error
	return entities, err
}

func (r *BaseRepo[T]) Count(filter any, data ...any) (int64, error) {
	var count int64
	err := r.db.Model(r.ModelPO).WithContext(r.Ctx).Where(filter, data...).Count(&count).Error
	return count, err
}
