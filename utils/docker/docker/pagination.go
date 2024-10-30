package docker

import "math"

// SlicePagination  切片分页
func SlicePagination[T interface{}](page int, limit int, data []T) []T {

	var total, start, end int
	total = len(data)
	// limit 大于总记录数,返回所有记录
	if limit > total {
		return data
	}

	// 总页数计算
	count := int(math.Ceil(float64(total) / float64(limit)))
	if page > count {
		return nil
	}

	start = (page - 1) * limit
	end = start + limit

	// 总记录数小于end,total 赋值给end
	if end > total {
		end = total
	}

	return data[start:end]
}
