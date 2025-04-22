package utils

import (
	"fmt"
	"forum/internal/article/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"math"
	"sort"
	"strconv"
	"time"
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

// SelIdForSuperAdmin
// @Description: 查询超级管理员对应的ID
// @Author wangyulong 2025-01-24 18:44:49
// @return       string
// @return       error
func SelIdForSuperAdmin() (string, error) {
	var superAdmin models.Role

	// 查询超级管理员对应的ID
	result := globals.DB.Model(&models.Role{}).Select("id").Where("name = ?", "超级管理员").First(&superAdmin)
	if result.Error != nil {
		return "", fmt.Errorf("SelApiId -> 查询超级管理员对应的ID异常 -> %s", result.Error)
	} else if result.RowsAffected == 0 {
		return "", fmt.Errorf("SelApiId -> 没有查询到超级管理员对应的ID")
	}
	return fmt.Sprintf("%v", superAdmin.ID), nil
}

// SelEmailForAdmin
// @Description: 根据管理员的ID查询其Email
// @Author wangyulong 2025-02-25 17:39:39
// @param        ID string
// @return       string
// @return       error
func SelEmailForAdmin(ID uint) (string, error) {

	var email string

	// 根据管理员的ID查询其Email
	res := globals.DB.Model(&models.Administrator{}).Select("email").Where("id = ?", ID).First(&email)
	if res.Error != nil {
		return "", fmt.Errorf("SelEmailForAdmin -> 根据管理员的ID查询其Email 异常 -> %v", res.Error)
	}
	return email, nil
}

// RedditHot 计算 Reddit 热度排序得分
func RedditHot(articles []requests.SearchArticleListRes) []requests.SearchArticleListRes {
	// 最终返回结果
	resArticlers := make([]requests.SearchArticleListRes, len(articles))
	// 创建映射
	type mp struct {
		id  int
		val float64
	}
	tmpArticles := make([]mp, len(articles))
	for i, v := range articles {
		tmpArticles[i].id = i
		tmpArticles[i].val = heat(v.Heat, v.PublishedAt)
	}
	sort.Slice(tmpArticles, func(i, j int) bool {
		return tmpArticles[i].val > tmpArticles[j].val
	})

	for i, v := range tmpArticles {
		//resArticlers = append(resArticlers, articles[v.id])
		resArticlers[i] = articles[v.id]
	}

	return resArticlers
}

// 返回加热度
func heat(score int, createdTime *time.Time) float64 {
	// 计算 log(得分)
	order := math.Log10(math.Max(math.Abs(float64(score)), 1))

	// 确定符号
	var sign float64
	if score > 0 {
		sign = 1
	} else if score < 0 {
		sign = -1
	} else {
		sign = 0
	}

	// Reddit 的基准时间戳：2005-12-08T07:46:43Z
	epoch := int64(1134028003)

	// 距离基准时间的秒数
	createdUTC := createdTime.Unix()
	seconds := createdUTC - epoch

	// 返回热度得分
	return math.Round((sign*order+float64(seconds)/10000)*1e7) / 1e7
}
