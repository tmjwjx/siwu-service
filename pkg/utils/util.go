package utils

// UintToMap uint 切片类型转换为map
func UintToMap(slice []uint) map[uint]bool {
	result := make(map[uint]bool)
	for _, v := range slice {
		result[v] = true
	}
	return result
}
