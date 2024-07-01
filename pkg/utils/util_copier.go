// Copyright 2023 QINIU. All rights reserved
// @Description: 运行出错的报警
// @Version: 1.0.0
// @Date: 2023/08/17 10:27
// @Author: liangfengyuan@qiniu.com

package utils

import (
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"github.com/jinzhu/copier"
	"reflect"
)

// CopySlice 进行slice的拷贝
func CopySlice[In any, Out any](in []In) []Out {
	outs := make([]Out, len(in))
	outType := reflect.TypeOf((*Out)(nil)).Elem()
	for i := 0; i < len(in); i++ {
		var out Out
		if outType.Kind() == reflect.Ptr {
			out = reflect.New(outType.Elem()).Interface().(Out)
		} else {
			out = reflect.New(outType).Elem().Interface().(Out)
		}
		_ = copier.Copy(&out, in[i])
		outs[i] = out
	}
	return outs
}

func MustCopyByCtx[T any](ctx jet.Ctx, in any) (t *T) {
	var err error
	if t, err = Copy[T](in); err != nil {
		ctx.Logger().Errorf("MustCopyByCtx error: %v", err.Error())
	}
	return
}

func MustCopy[T any](in any) (t *T) {
	var err error
	if t, err = Copy[T](in); err != nil {
		xlog.Errorf("MustCopy error: %v", err.Error())
	}
	return
}

func Copy[T any](in any) (*T, error) {
	t := new(T)
	err := copier.Copy(&t, in)
	return t, err
}
