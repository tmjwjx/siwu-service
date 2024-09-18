package utils

// UintToMap uint 切片类型转换为map
func UintToMap(slice []uint) map[uint]bool {
	result := make(map[uint]bool)
	for _, v := range slice {
		result[v] = true
	}
	return result
}

// IsUintSliContainUint 检查一个 uint slice 是否包含一个 uint
func IsUintSliContainUint(slice []uint, item uint) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
