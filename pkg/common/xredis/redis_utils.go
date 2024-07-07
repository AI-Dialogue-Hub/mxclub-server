package xredis

import (
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"mxclub/pkg/api"
	"mxclub/pkg/constant"
)

func GetOrDefault[T any](ctx jet.Ctx, key string, defaultFunc func() (*T, error)) (*T, error) {
	got, err := GetByString[T](ctx, key)
	if err == nil {
		return got, nil
	}
	got, err = defaultFunc()
	if err != nil {
		return nil, err
	}
	err = SetJSONStr(key, got, constant.Duration_7_Day)
	if err != nil {
		ctx.Logger().Errorf("SetJSONStr error:%v", err.Error())
	}
	return got, nil
}

func GetListOrDefault[T any](ctx jet.Ctx, listKey string, countKey string, defaultFunc func() ([]*T, int64, error)) ([]*T, int64, error) {
	list, listErr := GetByString[[]*T](ctx, listKey)
	count, countErr := GetByString[int64](ctx, countKey)
	if listErr == nil && countErr == nil {
		return *list, *count, nil
	}
	gotList, gotCount, err := defaultFunc()
	if err != nil {
		return nil, 0, err
	}
	if gotList != nil {
		if err = SetJSONStr(listKey, gotList, constant.Duration_7_Day); err != nil {
			_ = Del(listKey)
			ctx.Logger().Errorf("SetJSONStr error:%v", err.Error())
		}
	}
	if err = SetJSONStr(countKey, gotCount, constant.Duration_7_Day); err != nil {
		_ = Del(countKey)
		ctx.Logger().Errorf("SetJSONStr error:%v", err.Error())
	}
	return gotList, gotCount, nil
}

func BuildListDataCacheKey(prefix string, params *api.PageParams) string {
	return fmt.Sprintf("%s_ListData:Page%d:Size%d", prefix, params.Page, params.PageSize)
}

func BuildListCountCacheKey(prefix string) string {
	return prefix + "_ListCount"
}
