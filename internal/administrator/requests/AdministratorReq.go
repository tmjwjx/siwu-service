package requests

// AddAdministratorReq
// @Description: 添加管理员
// @Author tianjiajie 2025-02-21 15:25:26
type AddAdministratorReq struct {
	Email string `json:"email" binding:"email"`
}

// AdministratorReq
// @Description: 删除管理员
// @Author tianjiajie 2025-02-21 15:30:38
type AdministratorReq struct {
	ID int `json:"id" binding:"required"`
}

// UpdateAdministratorReq
// @Description: 编辑管理员信息
// @Author tianjiajie 2025-02-21 17:43:05
type UpdateAdministratorReq struct {
	ID       int    `json:"id" binding:"required"`
	Avatar   string `json:"avatar"`   // 头像
	Name     string `json:"name"`     // 名称
	Password string `json:"password"` // 密码
}

// GetAdministratorListReq
// @Description: 查询管理员列表
// @Author tianjiajie 2025-02-21 20:31:44
type GetAdministratorListReq struct {
	Page  int `json:"page" form:"page"`   // 页码
	Limit int `json:"limit" form:"limit"` // 每页数量
}

// GetAdministratorListRes
// @Description: 查询管理员列表响应
// @Author tianjiajie 2025-02-21 20:39:34
type GetAdministratorListRes struct {
	ID     uint   `json:"id"`     // ID
	Email  string `json:"email"`  // 邮箱，唯一
	Avatar string `json:"avatar"` // 头像
	Name   string `json:"name"`   // 名称
}

// GetAdministratorInfoRes
// @Description: 查询管理员信息响应
// @Author tianjiajie 2025-02-21 20:39:34
type GetAdministratorInfoRes struct {
	ID     uint   `json:"id"`     // ID
	Email  string `json:"email"`  // 邮箱，唯一
	Avatar string `json:"avatar"` // 头像
	Name   string `json:"name"`   // 名称
}
