package utils

import "github.com/fengyuan-liang/GoKit/collection/maps"

// NewLinkedHashMapWithPairs 封装了创建和初始化 LinkedHashMap 的逻辑
func NewLinkedHashMapWithPairs[K comparable, V any](pairs []*maps.Pair[K, V]) maps.IMap[K, V] {
	linkedHashMap := maps.NewLinkedHashMap[K, V]()
	linkedHashMap.PutAll(pairs)
	return linkedHashMap
}
