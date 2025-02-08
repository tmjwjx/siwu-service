package requests

// BackstageLoginReq 后台登录请求
type BackstageLoginReq struct {
	Email    string `json:"email"`    // 邮箱
	Password string `json:"password"` // 密码
}

// BackstageLoginRes 后台登录响应
type BackstageLoginRes struct {
	Id         uint     `json:"id"`          // 用户id
	Nickname   string   `json:"nickname"`    // 用户名
	Email      string   `json:"email"`       // 邮箱
	UserStatus int      `json:"user_status"` // 用户状态
	RoleIds    []uint   `json:"role_ids"`    // 角色的id集合
	RoleNames  []string `json:"role_names"`  // 角色的名称集合
	AvatarPath string   `json:"avatar_path"` // 用户头像
	Code       string   `json:"code"`        // 当前用户所拥有的排序最高的角色编码
}

// FinalBackstageLoginRes 最终版的后台登录响应
type FinalBackstageLoginRes struct {
	Token    string             `json:"token"`
	UserInfo *BackstageLoginRes `json:"userInfo"`
	Perm     *[]MenuPerm        `json:"perm"`
	CodeList *[]string          `json:"code_list"`
}

type MenuPerm struct {
	ID            uint   `json:"id"`
	ParentId      uint   `json:"pid"`            //父菜单
	RouteName     string `json:"route_name"`     // 路由名称
	Name          string `json:"name"`           // 菜单名称(权限点名称)
	Icon          string `json:"icon"`           // 图标名
	Sort          int    `json:"sort"`           // 排序
	RoutePath     string `json:"route_path"`     // 路由路径
	ComponentPath string `json:"component_path"` // 路由组件(文件路径)(前端要求存储的东西，后端不用管，只管存就好)
}
