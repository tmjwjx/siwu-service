package requests

// AssignMenuPermReq 为角色分配菜单权限
type AssignMenuPermReq struct {
	ID      uint   `json:"id"`
	PermIds []uint `json:"permIds"`
}

// GetMenuPermRes 获取当前角色的菜单权限(用于渲染侧边栏，只要type1和2)
type GetMenuPermRes struct {
	Perm []MenuPerm `json:"perm"`
}

type MenuPerm struct {
	ID        uint   `json:"id"`
	ParentId  *uint  `json:"pid"`        //父菜单
	Name      string `json:"name"`       // 菜单名称(权限点名称)
	Icon      string `json:"icon"`       // 图标名
	RouteName string `json:"route_name"` // 路由名称
}

// GetApiPermRes 获取当前角色的api权限
type GetApiPermRes struct {
	Group *[]uint `json:"group"`
}

//// GetApiPermRes 获取当前角色的api权限
//type GetApiPermRes struct {
//	Group []*Group `json:"group"`
//}
//
//type Group struct {
//	GroupId  uint     `json:"group_id"`
//	Children []*Child `json:"children"`
//}
//
//type Child struct {
//	ID uint `json:"id"`
//}

//// AssignApiPermReq 为角色分配api权限
//type AssignApiPermReq struct {
//	ID   uint   `json:"id"`
//	Apis []uint `json:"apis"`
//}

// AssignApiPermReq 为角色分配api权限
type AssignApiPermReq struct {
	ID   string   `json:"id"`
	Apis []string `json:"apis"`
}

// GetPermCodeRes 获取当前角色的所有权限标识
type GetPermCodeRes struct {
	CodeList []string `json:"code_list"`
}
