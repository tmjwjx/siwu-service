package requests

// AddAdministratorReq
// @Description: 添加管理员
// @Author tianjiajie 2025-02-21 15:25:26
type AddAdministratorReq struct {
	Email string `json:"email" binding:"email"`
}

// DeleteAdministratorReq
// @Description: 删除管理员
// @Author tianjiajie 2025-02-21 15:30:38
type DeleteAdministratorReq struct {
	ID int `json:"id" binding:"required"`
}
