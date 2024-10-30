package paginate

import (
	"fmt"
	"reflect"
)

// Paginate 切片分页
func Paginate(items interface{}, page, pageSize int) (interface{}, error) {
	slice := reflect.ValueOf(items)

	if slice.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Paginate: items is not a slice")
	}

	startIndex := (page - 1) * pageSize
	endIndex := page * pageSize

	// 确保索引范围不会超出切片范围
	if startIndex > slice.Len() {
		startIndex = slice.Len()
	}
	if endIndex > slice.Len() {
		endIndex = slice.Len()
	}

	paginatedSlice := slice.Slice(startIndex, endIndex).Interface()

	// 创建一个新的切片指针
	slicePtr := reflect.New(reflect.TypeOf(items))
	slicePtr.Elem().Set(reflect.ValueOf(paginatedSlice))

	return slicePtr.Interface(), nil
}
