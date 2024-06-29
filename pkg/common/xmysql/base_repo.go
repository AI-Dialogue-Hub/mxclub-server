package xmysql

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

type IBaseRepo[T any] interface {
	DB() *gorm.DB
	InsertOne(t *T) error
	InsertBatch(entities []interface{}) (int, error)
	InsertMany(entities []*T) (int, error)
	RemoveByID(id interface{}) error
	RemoveOne(filter interface{}) error
	Update(filter interface{}, update interface{}) error
	FindByID(id interface{}) (*T, error)
	Find(filter interface{}) ([]*T, error)
	FindOne(filter interface{}) (*T, error)
	FindAll() ([]*T, error)
	FindOrCreate(findFunc func() bool, t *T) (*T, error)
	List(pageNo int64, pageSize int64, filter interface{}) ([]*T, int64, error)
	Count(filter interface{}) (int64, error)
}

type BaseRepo[T any] struct {
	Db  *gorm.DB
	Ctx context.Context
}

func (r *BaseRepo[T]) DB() *gorm.DB {
	return r.Db
}

func (r *BaseRepo[T]) InsertOne(t *T) error {
	return r.Db.WithContext(r.Ctx).Create(t).Error
}

func (r *BaseRepo[T]) InsertBatch(entities []interface{}) (int, error) {
	if len(entities) == 0 {
		return 0, errors.New("cannot insert empty array")
	}
	result := r.Db.WithContext(r.Ctx).Create(entities)
	return int(result.RowsAffected), result.Error
}

func (r *BaseRepo[T]) InsertMany(entities []*T) (int, error) {
	if len(entities) == 0 {
		return 0, errors.New("cannot insert empty array")
	}
	result := r.Db.WithContext(r.Ctx).Create(entities)
	return int(result.RowsAffected), result.Error
}

func (r *BaseRepo[T]) RemoveByID(id interface{}) error {
	return r.Db.WithContext(r.Ctx).Delete(new(T), id).Error
}

func (r *BaseRepo[T]) RemoveOne(filter interface{}) error {
	return r.Db.WithContext(r.Ctx).Where(filter).Delete(new(T)).Error
}

func (r *BaseRepo[T]) Update(filter interface{}, update interface{}) error {
	return r.Db.WithContext(r.Ctx).Model(new(T)).Where(filter).Updates(update).Error
}

func (r *BaseRepo[T]) FindByID(id interface{}) (*T, error) {
	var t T
	err := r.Db.WithContext(r.Ctx).First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *BaseRepo[T]) Find(filter interface{}) ([]*T, error) {
	var entities []*T
	err := r.Db.WithContext(r.Ctx).Where(filter).Find(&entities).Error
	return entities, err
}

func (r *BaseRepo[T]) FindOne(filter interface{}) (*T, error) {
	var entity T
	err := r.Db.WithContext(r.Ctx).Where(filter).First(&entity).Error
	return &entity, err
}

func (r *BaseRepo[T]) FindAll() ([]*T, error) {
	var entities []*T
	err := r.Db.WithContext(r.Ctx).Find(&entities).Error
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

func (r *BaseRepo[T]) List(pageNo int64, pageSize int64, filter any) ([]*T, int64, error) {
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
	err := r.Db.WithContext(r.Ctx).Where(filter).Offset(int((pageNo - 1) * pageSize)).Limit(int(pageSize)).Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}
	var count int64
	err = r.Db.WithContext(r.Ctx).Where(filter).Count(&count).Error
	return entities, count, err
}

func (r *BaseRepo[T]) Count(filter any) (int64, error) {
	var count int64
	err := r.Db.WithContext(r.Ctx).Where(filter).Count(&count).Error
	return count, err
}
