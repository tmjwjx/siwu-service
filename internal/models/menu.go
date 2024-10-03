package models

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	ParentId *uint  `json:"pid"`  //父菜单
	Name     string `json:"name"` // 菜单名称(权限点名称)
	Code     string `json:"code"` // 权限标识
	Icon     string `json:"icon"` // 图标名

	Type          int    `json:"type"`           // 权限类型
	RouteName     string `json:"route_name"`     // 路由名称
	RoutePath     string `json:"route_path"`     // 路由路径
	RouteParam    string `json:"route_param"`    // 路由尾部参数
	ComponentPath string `json:"component_path"` // 路由组件(文件路径)

	Visible int    `json:"isVisible"` // 状态(显示隐藏)
	Sort    int    `json:"sort"`      // 排序
	Desc    string `json:"desc"`      // 描述

	//Redirect      string `json:"redirect"`      // 跳转路径
	//Roles string `json:"roles"` // 角色
	//RequestUrl    string `json:"requestUrl "`   // 请求接口
	//RequestMethod string `json:"requestMethod"` // 请求方法
}
