package models

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	ParentId uint   `json:"parent_id"` //父菜单(0表示没有父级)
	Name     string `json:"name"`      // 菜单名称(权限点名称)
	Code     string `json:"code"`      // 权限标识
	Icon     string `json:"icon"`      // 图标名(只是一个图片名，没有其他任何别的东西)

	Type          int    `json:"type"`           // 权限类型(用于判断是目录还是菜单或者按钮)
	RouteName     string `json:"route_name"`     // 路由名称
	RoutePath     string `json:"route_path"`     // 路由路径
	RouteParam    string `json:"route_param"`    // 路由尾部参数
	ComponentPath string `json:"component_path"` // 路由组件(文件路径)(前端要求存储的东西，后端不用管，只管存就好)

	Visible int    `json:"isVisible"` // 状态(显示隐藏)
	Sort    int    `json:"sort"`      // 排序
	Desc    string `json:"desc"`      // 描述

}
