package requests

// MenuSearchReq 检索获取所有菜单列表请求
type MenuSearchReq struct {
	Icon          string `json:"icon"`           // 图标名
	Name          string `json:"name"`           // 菜单名称(权限点名称)
	Type          int    `json:"type"`           // 权限类型
	RouteName     string `json:"route_name"`     // 路由名称
	RoutePath     string `json:"route_path"`     // 路由路径
	Visible       int    `json:"isVisible"`      // 状态(显示隐藏)
	Code          string `json:"code"`           // 权限标识
	ComponentPath string `json:"component_path"` // 路由组件(文件路径)
	ParentId      uint   `json:"pid"`            // 父菜单(0表示没有父级)
	Page          int    `json:"page"`           // 分页查询的起始位置
	Limit         int    `json:"limit"`          // 分页查询返回数据的条数
}

// MenuSearchRes 检索获取所有菜单列表响应
type MenuSearchRes struct {
	Menus []*Menus `json:"menus"`
	Total int      `json:"total"`
}

type Menus struct {
	ID       uint   `json:"id"`   // 站点id
	ParentId uint   `json:"pid"`  //父菜单
	Name     string `json:"name"` // 菜单名称(权限点名称)
	Code     string `json:"code"` // 权限标识
	Icon     string `json:"icon"` // 图标名

	Type          int    `json:"type"`           // 权限类型
	RouteName     string `json:"route_name"`     // 路由名称
	RoutePath     string `json:"route_path"`     // 路由路径
	ComponentPath string `json:"component_path"` // 路由组件(文件路径)

	Visible int    `json:"isVisible"` // 状态(显示隐藏)
	Sort    int    `json:"sort"`      // 排序
	Desc    string `json:"desc"`      // 描述
}

// GetMenuIconRes 获取所有菜单图标
type GetMenuIconRes struct {
	IconList []*IconRes `json:"icon_list"`
}

type IconRes struct {
	//Icon string `json:"icon"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type CreateMenuReq struct {
	ParentId uint   `json:"pid"`  //父菜单
	Name     string `json:"name"` // 菜单名称(权限点名称)
	Code     string `json:"code"` // 权限标识
	Icon     string `json:"icon"` // 图标名

	Type          int    `json:"type"`           // 权限类型
	RouteName     string `json:"route_name"`     // 路由名称
	RoutePath     string `json:"route_path"`     // 路由路径
	ComponentPath string `json:"component_path"` // 路由组件(文件路径)

	Visible    int    `json:"isVisible"`   // 状态(显示隐藏)
	Sort       int    `json:"sort"`        // 排序
	Desc       string `json:"desc"`        // 描述
	RouteParam string `json:"route_param"` // 如果不为空以/:拼接在route_path后面

	ApiIds []uint `json:"api_id"` // 按钮拥有的api
}

// DeleteMenuReq 删除菜单
type DeleteMenuReq struct {
	IDs []uint
}

// UpdateMenuReq 修改菜单
type UpdateMenuReq struct {
	ParentId      uint   `json:"pid"`            //父菜单
	Type          int    `json:"type"`           // 权限类型
	Icon          string `json:"icon"`           // 图标名
	Name          string `json:"name"`           // 菜单名称(权限点名称)
	Sort          int    `json:"sort"`           // 排序
	Visible       int    `json:"isVisible"`      // 状态(显示隐藏)
	RouteName     string `json:"route_name"`     // 路由名称
	RoutePath     string `json:"route_path"`     // 路由路径
	ComponentPath string `json:"component_path"` // 路由组件(文件路径)
	Desc          string `json:"desc"`           // 描述
	Code          string `json:"code"`           // 权限标识
	RouteParam    string `json:"route_param"`    // 如果不为空以/:拼接在route_path后面
	ID            uint   `json:"id"`             // 被修改菜单id

	ApiIds []uint `json:"api_id"` // 按钮拥有的api
}

// GetMenuDetailRes 获取当前菜单详情
type GetMenuDetailRes struct {
	ID       uint   `json:"id"`   // 菜单id
	ParentId uint   `json:"pid"`  //父菜单
	Name     string `json:"name"` // 菜单名称(权限点名称)
	Code     string `json:"code"` // 权限标识
	Icon     string `json:"icon"` // 图标名

	Type          int    `json:"type"`           // 权限类型
	RouteName     string `json:"route_name"`     // 路由名称
	RoutePath     string `json:"route_path"`     // 路由路径
	ComponentPath string `json:"component_path"` // 路由组件(文件路径)

	Visible    int    `json:"isVisible"`   // 状态(显示隐藏)
	RouteParam string `json:"route_param"` // 路由尾部参数
	Sort       int    `json:"sort"`        // 排序
	Desc       string `json:"desc"`        // 描述

	ApiIds []uint `json:"api_id"` // 按钮拥有的api
}

// GetSpecificMenuRes
// @Description: 获取所有type为1和2的菜单
// @Author wangyulong 2024-10-15 11:06:19
type GetSpecificMenuRes struct {
	MenuTypeList *[]List `json:"menu_type_list"`
}

type List struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ParentID uint   `json:"pid"`
}
