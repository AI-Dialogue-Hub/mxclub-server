package utils

import "github.com/fengyuan-liang/GoKit/collection/maps"

// NewLinkedHashMapWithPairs 封装了创建和初始化 LinkedHashMap 的逻辑
func NewLinkedHashMapWithPairs[K comparable, V any](pairs []*maps.Pair[K, V]) maps.IMap[K, V] {
	linkedHashMap := maps.NewLinkedHashMap[K, V]()
	linkedHashMap.PutAll(pairs)
	return linkedHashMap
}

func ReverseLinkedHashMap[K comparable, V comparable](linkedHashMap maps.IMap[K, V]) maps.IMap[V, K] {
	reversedLinkedHashMap := maps.NewLinkedHashMap[V, K]()
	linkedHashMap.ForEach(func(k K, v V) {
		reversedLinkedHashMap.Put(v, k)
	})
	return reversedLinkedHashMap
}

func SliceToMap[T any, KEY comparable](arr []T, f func(ele T) KEY) maps.IMap[KEY, []T] {
	linkedHashMap := maps.NewLinkedHashMap[KEY, []T]()
	for _, ele := range arr {
		parseKey := f(ele)
		if linkedHashMap.ContainsKey(parseKey) {
			linkedHashMap.Put(parseKey, append(linkedHashMap.MustGet(parseKey), ele))
		} else {
			elements := []T{ele}
			linkedHashMap.Put(parseKey, elements)
		}
	}
	return linkedHashMap
}

func SliceToRawMap[T any, KEY comparable](arr []T, f func(ele T) KEY) map[KEY]T {
	hashMap := make(map[KEY]T)
	for _, ele := range arr {
		hashMap[f(ele)] = ele
	}
	return hashMap
}

func SliceToSingleMap[T any, KEY comparable](arr []T, f func(ele T) KEY) maps.IMap[KEY, T] {
	linkedHashMap := maps.NewLinkedHashMap[KEY, T]()
	for _, ele := range arr {
		parseKey := f(ele)
		linkedHashMap.Put(parseKey, ele)
	}
	return linkedHashMap
}
