package logics

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/internalPkg/sqlUtils"
	"forum/internal/models"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"forum/pkg/casbin"
	"forum/pkg/globals"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"mime/multipart"
	"time"
)

// Reset 重置用户密码
func (u *UserReqContext) Reset(msg requests.ReseatReq) error {
	// 根据id查询到用户email
	user := repositories.QueryUserById(u.DB, msg.Id)
	// 默认密码为用户邮箱，将默认密码加密
	defaultPassword := user.Email
	encryptedPassword, err := internalUtils.HashPassword(defaultPassword)
	if err != nil {
		return fmt.Errorf("UserReqContext.Reset() : 密码%s加密失败", defaultPassword)
	}

	// 更新密码
	if err = sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: msg.Id}}, map[string]interface{}{"password": encryptedPassword}); err != nil {
		return fmt.Errorf("UserReqContext.Reset() err: %v", err)
	}

	return nil
}

// Add 添加用户
func (u *UserReqContext) Add(req requests.AddReq) (uint, error) {
	// 默认密码为用户邮箱，将默认密码加密
	defaultPassword := req.Email
	encryptedPassword, err := internalUtils.HashPassword(defaultPassword)
	if err != nil {
		return 0, fmt.Errorf("UserReqContext.Add() : 密码%s加密失败", defaultPassword)
	}

	// 判断该 email 是否使用过
	user := repositories.QueryUserByEmail(u.DB, req.Email)
	if user != nil {
		return 0, fmt.Errorf("UserReqContext.Add() err: email %v 已经被使用", req.Email)
	}

	for _, v := range req.RoleIds {
		// 判断 角色id 是否存在
		role := repositories.QueryRoleById(u.DB, v)
		if role == nil {
			return 0, fmt.Errorf("UserReqContext.Add() err: Role 表中不存在id为 %v 的角色", v)
		}
	}

	// 插入User表
	if err = sqlUtils.InsertObject(u.DB, &models.User{Nickname: req.NickName, Email: req.Email, Password: encryptedPassword, Status: req.UserStatus, LastLoginTime: time.Now()}); err != nil {
		return 0, fmt.Errorf("UserReqContext.Add() err: %v", err)
	}

	// 通过email查找id
	user = repositories.QueryUserByEmail(u.DB, req.Email)

	// 插入UserDetail表
	if err = sqlUtils.InsertObject(u.DB, &models.UserDetail{UserID: user.ID}); err != nil {
		return 0, fmt.Errorf("UserReqContext.Add() err: %v", err)
	}

	// 给用户分配角色id，调用 casbin 方法
	casbinService, err := casbin.NewCasbinService(globals.DB)
	if err != nil {
		return 0, fmt.Errorf("UserReqContext.Add() err: %v", err)
	}
	if err = casbinService.AssignRolesForAdminOrUser(user.Email, req.RoleIds); err != nil {
		return 0, fmt.Errorf("UserReqContext.Add() err: %v", err)
	}

	// 存储用户头像路径
	if err = internalUtils.StoreUrl(&internalUtils.UrlParam{
		UrlPath: []string{req.AvatarPath},
		Home:    globals.UserHome,
		HomeID:  user.ID,
		DB:      u.DB,
	}); err != nil {
		return 0, fmt.Errorf("UserReqContext.Add() err: %v", err)
	}

	return user.ID, nil
}

// Delete 删除用户
func (u *UserReqContext) Delete(req requests.DeleteReq) error {
	for _, v := range req.Ids {
		// 删除用户
		if _, err := sqlUtils.DeleteObjectsByModel(u.DB, &models.User{}, map[string]interface{}{"id": v}); err != nil {
			return fmt.Errorf("UserReqContext.Delete() err: %v", err)
		}

		// 删除该用户的详细信息
		if _, err := sqlUtils.DeleteObjectsByModel(u.DB, &models.UserDetail{}, map[string]interface{}{"user_id": v}); err != nil {
			return fmt.Errorf("UserReqContext.Delete() err: %v", err)
		}

		// 删除用户对应的角色id
		casbinService, err := casbin.NewCasbinService(globals.DB)
		if err != nil {
			return fmt.Errorf("UserReqContext.Delete() err: %v", err)
		}
		// 查询该id对应的邮箱
		email := repositories.QueryUserEmailById(u.DB, v)
		// 获取该用户对应的全部角色id
		roleIds, err := casbinService.GetRolesForAdminOrUser(email)
		if err != nil {
			return fmt.Errorf("UserReqContext.Delete() err: %v", err)
		}
		// 根据用户 id 删除 roleIds
		if err = casbinService.DeleteRoleForAdminOrUser("v", roleIds); err != nil {
			return fmt.Errorf("UserReqContext.Delete() err: %v", err)
		}
	}
	return nil
}

// Edit 编辑用户
func (u *UserReqContext) Edit(req requests.EditReq) error {
	var user *models.User
	// 判断是否存在该 userId
	if user = repositories.QueryUserById(u.DB, req.UserId); user == nil {
		return fmt.Errorf("UserReqContext.Edit() err: 不存在id为 %v 的用户", req.UserId)
	}
	// 判断邮箱是否已经使用
	if user = repositories.QueryUserByEmail(u.DB, req.Email); user != nil && user.ID != req.UserId {
		return fmt.Errorf("UserReqContext.Edit() err: 邮箱 %v 已经被使用", req.Email)
	}
	// 判断是否存在要添加的 RoleId
	for _, v := range req.RoleIds {
		if role := repositories.QueryRoleById(u.DB, v); role == nil {
			return fmt.Errorf("UserReqContext.Edit() err: 不存在id为 %v 的角色", v)
		}
	}

	// 更改 nickname、email、status
	if err := sqlUtils.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: req.UserId}}, map[string]interface{}{"nickname": req.NickName, "email": req.Email, "status": req.UserStatus}); err != nil {
		return fmt.Errorf("UserReqContext.Edit() -> %v", err)
	}
	// 更改用户对应的 role_id
	casbinService, err := casbin.NewCasbinService(globals.DB)
	if err != nil {
		return fmt.Errorf("UserReqContext.Delete() err: %v", err)
	}
	if err = casbinService.UpdateRoleForAdminOrUser(req.Email, req.RoleIds); err != nil {
		return fmt.Errorf("UserReqContext.Delete() err: %v", err)
	}

	// 更改用户对应的 role_id 后还要在管理者表中添加这个用户
	for _, v := range req.RoleIds {
		if v != 1 && user != nil { // 如果角色id != 1，那么说明这个角色不是用户，而是管理员
			if err = sqlUtils.InsertObject(u.DB, &models.Administrator{Email: user.Email, Password: user.Password}); err != nil {
				return err
			}
		}
	}

	// 存储用户头像路径
	if err = internalUtils.StoreUrl(&internalUtils.UrlParam{
		UrlPath: []string{req.AvatarPath},
		Home:    globals.UserHome,
		HomeID:  req.UserId,
		DB:      u.DB,
	}); err != nil {
		return fmt.Errorf("UserReqContext.Edit() -> %v", err)
	}

	return nil
}

// List 获取所有用户列表
func (u *UserReqContext) List(req requests.ListReq) ([]*requests.ListRes, int, error) {
	// 查询符合条件的角色（除了符合 req.RoleNames）（nickname使用模糊查询）
	users, total, err := repositories.QueryUserListByPage(u.DB, req)
	if err != nil {
		return nil, 0, fmt.Errorf("UserReqContext.List() %v", err)
	}

	// 整理返回切片数据
	var listRes = make([]*requests.ListRes, 0)
	for _, v := range users {
		// 查询用户的头像路径
		userImages, err := internalUtils.GetImages(u.DB, globals.UserHome, v.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("UserReqContext.List() %v", err)
		}
		// 没有图片
		if userImages == nil {
			return nil, 0, fmt.Errorf("UserReqContext.List() err = 无法找到id为%d的用户头像图片", v.ID)
		}
		avatarPath := (*userImages)[0]

		// 获取该用户对应的全部角色id
		casbinService, err := casbin.NewCasbinService(globals.DB)
		if err != nil {
			return nil, 0, fmt.Errorf("UserReqContext.List() %v", err)
		}
		roleIds, err := casbinService.GetRolesForAdminOrUser(v.Email)
		if err != nil {
			return nil, 0, fmt.Errorf("UserReqContext.List() %v", err)
		}
		// 通过角色id查找角色名称
		roleNames := make([]string, 0)
		for _, v := range roleIds {
			role := repositories.QueryRoleById(u.DB, v)
			if role == nil {
				return nil, 0, fmt.Errorf("UserReqContext.List() err -> 没有查询到角色 id:%d", roleIds)
			}
			roleNames = append(roleNames, role.Name)
		}

		listRes = append(listRes, &requests.ListRes{
			Id:         v.ID,
			AvatarPath: avatarPath,
			NickName:   v.Nickname,
			Email:      v.Email,
			Heat:       v.Heat,
			FansCount:  v.FansCount,
			RoleNames:  roleNames,
			UserStatus: v.Status,
			// 将 time.Time 格式化为字符串
			LastLoginTime: v.LastLoginTime.Format("2006-01-02 15:04:05"),
			CreateTime:    v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 根据前端的要求，翻转 listRes
	// slices.Reverse(listRes)

	return listRes, total, nil
}

// Import 导入用户表
func (u *UserReqContext) Import(file *multipart.FileHeader) error {
	// 打开上传的文件
	openFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("UserReqContext.Import() err: Failed to open file")
	}
	defer openFile.Close()
	// 解析 Excel 文件
	f, err := excelize.OpenReader(openFile)
	if err != nil {
		return fmt.Errorf("UserReqContext.Import() err: Failed to read Excel file")
	}
	// 获取第一个工作表的所有行
	rows, err := f.GetRows("User Template")
	if err != nil {
		return fmt.Errorf("UserReqContext.Import() err: Failed to get rows from Excel")
	}

	var l = 0 // 表头长度
	// 从第二行开始读取数据
	for i, row := range rows {
		// 跳过表头
		if i == 0 {
			l = len(row)
			continue
		}
		// 表格第i行后面的一些数据是nil，数据不够
		if len(row) < l {
			return fmt.Errorf("UserReqContext.Import() err: 第 %v 行数据，数据量不足", i)
		}

		// 根据列顺序解析每一行数据
		nickname := row[0]
		email := row[1]
		password := row[2]
		status := 0
		if n, err := fmt.Sscanf(row[3], "%d", &status); n == 0 || err != nil {
			return fmt.Errorf("UserReqContext.Import() err: Status为 %v 不合法", row[3])
		}

		// 判断邮箱是否合法，是否已经使用
		if !internalUtils.IsValidEmail(email) {
			return fmt.Errorf("UserReqContext.Import() err: email为 %v 不合法", email)
		}
		if user := repositories.QueryUserByEmail(u.DB, email); user != nil {
			return fmt.Errorf("UserReqContext.Import() err: email为 %v 已经被使用", email)
		}
		encryptedPassword := ""
		// 判断密码是否合法，并加密
		if password == "" { // 如果密码是空，就选择默认密码，默认密码为用户邮箱
			password = email
		} else if !internalUtils.IsValidPassword(password) { // 如果不为空，就判断是否合法
			return fmt.Errorf("UserReqContext.Import() err: 密码必须要同时包含字母、数字、特殊字符，长度在8到20位之间")
		}
		// 加密
		encryptedPassword, err = internalUtils.HashPassword(password)
		if err != nil {
			return fmt.Errorf("UserReqContext.Import() err: 密码%s加密失败", password)
		}
		// 检查Status，Status只能是1或者2
		if status != 1 && status != 2 {
			return fmt.Errorf("UserReqContext.Import() err: status为 %v 不合法", status)
		}

		// 将数据存储到数据库中
		user := &models.User{
			Nickname:      nickname,
			Email:         email,
			Password:      encryptedPassword,
			Status:        status,
			LastLoginTime: time.Now(),
		}

		// 保存到 User 表中
		if err = sqlUtils.InsertObject(u.DB, user); err != nil {
			return fmt.Errorf("UserReqContext.Import() err: 没有将邮箱为%v的用户存储在User表中", email)
		}
		// 查询该 user 的id，并 添加到 UserDetail 表中
		user = repositories.QueryUserByEmail(u.DB, email)
		if user == nil {
			return fmt.Errorf("UserReqContext.Import() err: 没有将邮箱为%v的用户存储在User表中", email)
		}
		// 保存到 UserDetail 表中
		if err = sqlUtils.InsertObject(u.DB, &models.UserDetail{UserID: user.ID}); err != nil {
			return fmt.Errorf("UserReqContext.Import() err: 没有将邮箱为%v的用户存储在Userdetail表中", email)
		}
	}

	return nil
}

// Export 导出用户表
func (u *UserReqContext) Export() error {
	// 从数据库获取所有用户
	users, err := repositories.QueryAllUser(u.DB)
	if err != nil {
		return fmt.Errorf("UserReqContext.Export() -> %v", err)
	}

	// 创建一个新的 Excel 文件
	f := excelize.NewFile()
	// 创建一个新的 Sheet，工作表名称为 "Users"
	index, _ := f.NewSheet("Users")
	// 添加列头到第一行
	headers := []string{"Id", "创造时间", "更新时间", "昵称", "邮箱", "密码", "热度", "关注数量", "粉丝数量", "状态", "最后登录时间"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue("Users", cell, header)
	}

	// 写入用户数据
	for i, user := range users {
		row := i + 2 // 从第二行开始写入数据
		f.SetCellValue("Users", fmt.Sprintf("A%d", row), user.ID)
		f.SetCellValue("Users", fmt.Sprintf("B%d", row), user.CreatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue("Users", fmt.Sprintf("C%d", row), user.UpdatedAt.Format("2006-01-02 15:04:05"))
		f.SetCellValue("Users", fmt.Sprintf("D%d", row), user.Nickname)
		f.SetCellValue("Users", fmt.Sprintf("E%d", row), user.Email)
		f.SetCellValue("Users", fmt.Sprintf("F%d", row), user.Password)
		f.SetCellValue("Users", fmt.Sprintf("G%d", row), user.Heat)
		f.SetCellValue("Users", fmt.Sprintf("H%d", row), user.AttentionCount)
		f.SetCellValue("Users", fmt.Sprintf("I%d", row), user.FansCount)
		f.SetCellValue("Users", fmt.Sprintf("J%d", row), user.Status)
		f.SetCellValue("Users", fmt.Sprintf("K%d", row), user.LastLoginTime.Format("2006-01-02 15:04:05"))
	}

	// 设置活动工作表
	f.SetActiveSheet(index)

	// 写入响应
	if err = f.Write(u.Ctx.Writer); err != nil {
		return fmt.Errorf("UserReqContext.Export() -> %v", err)
	}
	return nil
}

// DownloadTemplate 下载导入用户模版excel
func (u *UserReqContext) DownloadTemplate() error {
	// 创建一个新的 Excel 文件
	f := excelize.NewFile()
	// 创建一个新的 Sheet，工作表名称为 "User Template"
	index, _ := f.NewSheet("User Template")

	// 添加列头到第一行
	headers := []string{"昵称", "邮箱", "密码", "状态"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i))) // 将索引转换为对应的列号 (A, B, C...)
		// 设置单元格值
		if err := f.SetCellValue("User Template", cell, header); err != nil {
			return fmt.Errorf("UserReqContext.DownloadTemplate() err: %v", err)
		}
	}
	// 将 "User Template" 设为活动工作表
	f.SetActiveSheet(index)

	// 将文件内容写入响应中
	if err := f.Write(u.Ctx.Writer); err != nil {
		return fmt.Errorf("UserReqContext.DownloadTemplate() err: Failed to create Excel file")
	}
	return nil
}

// GetInfo 获取当前用户基本信息
func (u *UserReqContext) GetInfo(id uint) (*requests.GetInfoRes, error) {
	// 判断id是否存在
	user := repositories.QueryUserById(u.DB, id)
	if user == nil {
		return nil, fmt.Errorf("UserReqContext.GetInfo() err: 不存在id为 %v 的用户", id)
	}

	// 查询用户的头像路径
	userImages, err := internalUtils.GetImages(u.DB, globals.UserHome, id)
	if err != nil {
		return nil, fmt.Errorf("UserReqContext.GetInfo() %v", err)
	}
	// 没有图片
	if userImages == nil {
		return nil, fmt.Errorf("UserReqContext.GetInfo() err = 无法找到id为%d的用户头像图片", id)
	}
	avatarPath := (*userImages)[0]

	// 获取该用户对应的全部角色id
	casbinService, err := casbin.NewCasbinService(globals.DB)
	if err != nil {
		return nil, fmt.Errorf("UserReqContext.GetInfo() %v", err)
	}
	roleIds, err := casbinService.GetRolesForAdminOrUser(user.Email)
	if err != nil {
		return nil, fmt.Errorf("UserReqContext.GetInfo() %v", err)
	}

	var getInfoRes = &requests.GetInfoRes{
		Id:         user.ID,
		AvatarPath: avatarPath,
		NickName:   user.Nickname,
		Email:      user.Email,
		UserStatus: user.Status,
		RoleIds:    roleIds,
	}

	return getInfoRes, nil
}
