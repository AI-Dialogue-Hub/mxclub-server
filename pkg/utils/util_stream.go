// Copyright 2023 QINIU. All rights reserved
// @Description: 常用的集合操作
// @Version: 1.0.0
// @Date: 2023/10/09 15:08
// @Author: liangfengyuan@qiniu.com

package utils

func Filter[T any](in []T, f func(in T) bool) []T {
	out := make([]T, 0)
	for _, entity := range in {
		if f(entity) {
			out = append(out, entity)
		}
	}
	return out
}

func Map[In any, Out any](in []In, f func(in In) Out) []Out {
	out := make([]Out, len(in))
	for index, entity := range in {
		out[index] = f(entity)
	}
	return out
}

// MapFromSliceToSlice slice -> func(in) slice -> slice  接收一个数组并拼接到最终的结果中
func MapFromSliceToSlice[In any, Out any](in []In, f func(in In) []Out) []Out {
	out := make([]Out, len(in))
	for _, entity := range in {
		out = append(out, f(entity)...)
	}
	return out
}

func ForEach[T any](entities []*T, f func(t *T)) {
	for _, entity := range entities {
		f(entity)
	}
}
