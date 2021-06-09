package util

func GetPage(pageNum int, pageSize int) int {
	result := 0
	if pageNum > 0 {
		result = (pageNum - 1) * pageSize
	}
	return result
}
