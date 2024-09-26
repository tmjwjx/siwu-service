package casbin

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/globals"
	"regexp"
	"strings"
)

// SelMenId 根据请求路径，获取菜单ID
func SelMenId(requestUrl, requestMethod string) (string, error) {
	var menuId string
	// 获取最后一个/后面的信息
	re := regexp.MustCompile("/([^/]+)$")
	matches := re.FindStringSubmatch(requestUrl)
	// 判断是不是角色
	err := SelMRole(matches[1])
	if err != nil {
		requestUrl = strings.Replace(requestUrl, matches[0], "", 1)
	} else {
		return "", err
	}
	// 根据请求中的路由接口,查询menuId
	err = globals.DB.Model(models.Menu{}).Select("id").Where("request_url = ? and request_method", requestUrl, requestMethod).Scan(&menuId).Error
	return menuId, fmt.Errorf("SelMenId -> 根据请求路径，获取菜单ID -> %s", err)
}

// SelMRole 模糊查询角色字段
func SelMRole(role string) error {

	var count int64

	err := globals.DB.Table("casbin_rule").Where("v1 LIKE ?", role+"%").Count(&count).Error
	if err != nil {
		return fmt.Errorf("SelMRole -> 模糊查询角色字段失败 -> %s", err)
	}

	// 如果count>0,表示找到了记录
	if count > 0 {
		return nil
	}

	// 如果没有找到记录，返回false
	return fmt.Errorf("SelMRole -> 该角色不存在")

}
