package logics

import (
	"fmt"
	"forum/internal/internal_pkg/internal_utils"
	"forum/internal/models"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"gorm.io/gorm"
)

// Reset 重置用户密码
func (u *UserReqContext) Reset(msg requests.ReseatReq) error {
	// 将默认密码加密
	defaultPassword := "abc123"
	encryptedPassword, err := internal_utils.HashPassword(defaultPassword)
	if err != nil {
		return fmt.Errorf("UserReqContext.Reset() : 密码%s加密失败", defaultPassword)
	}

	// 更新密码
	if err := repositories.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: msg.Id}}, map[string]interface{}{"password": encryptedPassword}); err != nil {
		return fmt.Errorf("UserReqContext.Reset() err: %v", err)
	}

	return nil
}

// Add 添加用户
func (u *UserReqContext) Add(req requests.AddAndEditReq) error {
	// 将默认密码加密
	defaultPassword := "abc123"
	encryptedPassword, err := internal_utils.HashPassword(defaultPassword)
	if err != nil {
		return fmt.Errorf("UserReqContext.Add() : 密码%s加密失败", defaultPassword)
	}

	// 判断该 email 是否已经存在
	user := repositories.QueryUserByEmail(u.DB, req.Email)
	if user != nil { // 已经存在该 email
		return fmt.Errorf("UserReqContext.Add() err: email %v 已经被使用", req.Email)
	}

	for _, v := range req.RoleIds {
		// 判断 角色id 是否存在
		role := repositories.QueryRoleById(u.DB, v)
		if role == nil {
			return fmt.Errorf("UserReqContext.Add() err: Role 表中不存在id为 %v 的角色", v)
		}
	}

	// 插入User表
	err = repositories.InsertObject(u.DB, &models.User{Nickname: req.NickName, Email: req.Email, Password: encryptedPassword, Status: req.UserStatus})
	if err != nil {
		return fmt.Errorf("UserReqContext.Add() err: %v", err)
	}

	// 通过email查找id
	user = repositories.QueryUserByEmail(u.DB, req.Email)

	// 插入 AdminRole 表（插入用户对应的角色）
	for _, v := range req.RoleIds {
		if err = repositories.InsertObject(u.DB, &models.AdminRole{AdminId: user.ID, RoleId: v}); err != nil {
			return fmt.Errorf("UserReqContext.Add() err: %v", err)
		}
	}

	return nil
}

// Delete 删除用户
func (u *UserReqContext) Delete(req requests.DeleteReq) error {
	num := 0
	for _, v := range req.Ids {
		n, err := repositories.DeleteObjectsByModel(u.DB, &models.User{}, map[string]interface{}{"id": v})
		if err != nil {
			return fmt.Errorf("UserReqContext.Delete() err: %v", err)
		}
		num += int(n)
	}
	fmt.Printf("应该删除 %d 条数据，实际删除 %v 条数据\n", len(req.Ids), num)

	return nil
}

// Edit 编辑用户
func (u *UserReqContext) Edit(req requests.AddAndEditReq) error {

	// // 更改 nickname、email、status
	// if err := repositories.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: req.Id}}, map[string]interface{}{"nickname": req.NickName, "email": req.Email, "status": req.UserStatus}); err != nil {
	// 	return fmt.Errorf("UserReqContext.Edit() err: %v", err)
	// }

	// 更改 role_ids
	if err := repositories.UpdateObjects(u.DB, &models.AdminRole{}, map[string]interface{}{"nickname": req.NickName, "email": req.Email, "status": req.UserStatus}); err != nil {
		return fmt.Errorf("UserReqContext.Edit() err: %v", err)
	}

	return nil
}

// List 获取所有用户列表
func (u *UserReqContext) List(req requests.ListReq) ([]*requests.ListRes, error) {
	conditions := make(map[string]interface{})
	// 判断是否该添加某些查询条件（如果某些条件为空，那么就不查询这个条件）
	if req.NickName != "" {
		conditions["nick_name"] = req.NickName
	}
	if req.Email != "" {
		conditions["email"] = req.Email
	}

	// len(req.RoleIds) != 0

	// req.UserStatus == 0 代表着全部
	conditions["status"] = req.UserStatus

	// limit, page, m

	return nil, nil
}

// // Import 导入用户表
// func (u *UserReqContext) Import(req requests.AddAndEditReq) ([]*requests.UserRankRes, error) {
//
// }

// // Export 导出用户表
// func (u *UserReqContext) Export(req requests.AddAndEditReq) ([]*requests.UserRankRes, error) {
//
// }
//
// // ImportTemplate 下载导入用户模版excel
// func (u *UserReqContext) ImportTemplate(req requests.AddAndEditReq) ([]*requests.UserRankRes, error) {
//
// }
//
// // GetInfo 获取当前用户基本信息
// func (u *UserReqContext) GetInfo(req requests.AddAndEditReq) ([]*requests.UserRankRes, error) {
//
// }
