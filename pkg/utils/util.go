package utils

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/globals"
	"strconv"
)

// ChangeStringToUint
// @Description: 将string类型的值转换成uint类型
// @Author wangyulong 2024-10-09 15:26:15
// @param        str string
// @return       uint
// @return       error
func ChangeStringToUint(str string) (uint, error) {

	// 将字符串转换成uint64, 基数为 10, 位大小为 64 位
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("ChangeType -> 将string类型的值转换成uint类型失败 -> %s", err)
	}

	// 将 uint64 转换成 uint 类型
	uintNum := uint(num)

	return uintNum, nil
}

// SelApiId 根据请求路径查找api的id
func SelApiId(requestUrl string) (uint, error) {
	var apiId uint
	// 根据请求中的路由接口,查询apiId
	err := globals.DB.Model(models.Api{}).Select("ID").Where("path = ?", requestUrl).Scan(&apiId).Error
	if err != nil {
		return 0, fmt.Errorf("SelApiId -> 根据请求路径，获取apiID失败，该用户没有该权限 -> %s", err)
	}
	if err == nil && apiId == 0 {
		return 0, fmt.Errorf("SelApiId -> 该api不存在，该用户没有该权限")
	}
	return apiId, nil
}
