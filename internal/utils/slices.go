package utils

import "slices"

// DelElements 删除切片中的元素
func DelElements[T comparable](slice []T, elems ...T) []T {
	if len(slice) == 0 || len(elems) == 0 {
		return slice
	}
	return slices.DeleteFunc(slice, func(e T) bool {
		return slices.Contains(elems, e)
	})
}

// OnDemandAppend 按需扩容切片,
// 为长期驻留内存的数据使用
func OnDemandAppend[T any](slice []T, elems ...T) []T {
	if len(elems) == 0 {
		return slice
	}
	newSlice := make([]T, len(slice), len(slice)+len(elems))
	copy(newSlice, slice)
	return append(newSlice, elems...)
}

// NoReduplicateAppend 去重并按需扩容切片,
// 为长期驻留内存的数据使用
func NoReduplicateAppend[T comparable](slice []T, elems ...T) []T {
	DelElements(elems, slice...)
	return OnDemandAppend(slice, elems...)
}
