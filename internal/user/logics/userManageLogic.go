package logics

import (
	"fmt"
	"forum/internal/image/controllers"
	"forum/internal/internal_pkg/internal_utils"
	"forum/internal/models"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"mime/multipart"
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
func (u *UserReqContext) Add(req requests.AddReq) (uint, error) {
	// 将默认密码加密
	defaultPassword := "abc123"
	encryptedPassword, err := internal_utils.HashPassword(defaultPassword)
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
	if err = repositories.InsertObject(u.DB, &models.User{Nickname: req.NickName, Email: req.Email, Password: encryptedPassword, Status: req.UserStatus}); err != nil {
		return 0, fmt.Errorf("UserReqContext.Add() err: %v", err)
	}

	// 通过email查找id
	user = repositories.QueryUserByEmail(u.DB, req.Email)

	// 插入UserDetail表
	if err = repositories.InsertObject(u.DB, &models.UserDetail{UserID: user.ID}); err != nil {
		return 0, fmt.Errorf("UserReqContext.Add() err: %v", err)
	}

	// 插入 AdminRole 表（插入用户对应的角色）
	for _, v := range req.RoleIds {
		if err = repositories.InsertObject(u.DB, &models.AdminRole{AdminId: user.ID, RoleId: v}); err != nil {
			return 0, fmt.Errorf("UserReqContext.Add() err: %v", err)
		}
	}

	return user.ID, nil
}

// Delete 删除用户
func (u *UserReqContext) Delete(req requests.DeleteReq) error {
	num := 0
	for _, v := range req.Ids {
		// 删除用户
		n, err := repositories.DeleteObjectsByModel(u.DB, &models.User{}, map[string]interface{}{"id": v})
		if err != nil {
			return fmt.Errorf("UserReqContext.Delete() err: %v", err)
		}
		num += int(n)

		// 删除用户对应的角色id
		if _, err = repositories.DeleteObjectsByModel(u.DB, &models.AdminRole{}, map[string]interface{}{"admin_id": v}); err != nil {
			return fmt.Errorf("UserReqContext.Delete() err: %v", err)
		}
	}
	fmt.Printf("应该删除 %d 条数据，实际删除 %v 条数据\n", len(req.Ids), num)

	return nil
}

// Edit 编辑用户
func (u *UserReqContext) Edit(req requests.EditReq) error {
	// 判断是否存在该 userId
	if user := repositories.QueryUserById(u.DB, req.UserId); user == nil {
		return fmt.Errorf("UserReqContext.Edit() err: 不存在id为 %v 的用户", req.UserId)
	}
	// 判断邮箱是否已经使用
	if user := repositories.QueryUserByEmail(u.DB, req.Email); user != nil && user.ID != req.UserId {
		return fmt.Errorf("UserReqContext.Edit() err: 邮箱 %v 已经被使用", req.Email)
	}
	// 判断是否存在要添加的 RoleId
	for _, v := range req.RoleIds {
		if role := repositories.QueryRoleById(u.DB, v); role == nil {
			return fmt.Errorf("UserReqContext.Edit() err: 不存在id为 %v 的角色", v)
		}
	}

	// 更改 nickname、email、status
	if err := repositories.UpdateObjects(u.DB, &models.User{Model: gorm.Model{ID: req.UserId}}, map[string]interface{}{"nickname": req.NickName, "email": req.Email, "status": req.UserStatus}); err != nil {
		return fmt.Errorf("UserReqContext.Edit() -> %v", err)
	}
	// 更改用户对应的 role_id
	if err := repositories.UpdateAdminRoles(u.DB, req.UserId, req.RoleIds); err != nil {
		return fmt.Errorf("UserReqContext.Edit() -> %v", err)
	}

	return nil
}

// List 获取所有用户列表
func (u *UserReqContext) List(req requests.ListReq) ([]*requests.ListRes, int, error) {
	conditions := make(map[string]interface{})
	// 判断是否该添加某些查询条件（如果某些条件为空或为0，代表着不查询这个条件）
	if req.NickName != "" {
		conditions["nickname"] = req.NickName
	}
	if req.Email != "" {
		conditions["email"] = req.Email
	}
	// status = 0 表示全部
	if req.UserStatus != 0 {
		conditions["status"] = req.UserStatus
	}

	// 查询符合条件的角色（除了符合 req.RoleIds）
	users, total, err := repositories.QueryUserListByPage(u.DB, conditions, req.Page, req.Limit, req.RoleIds, req.Heat, req.FansCount, req.CreateTime, req.LastLoginTime)
	if err != nil {
		return nil, 0, fmt.Errorf("UserReqContext.List() %v", err)
	}

	// 整理返回切片数据
	var listRes = make([]*requests.ListRes, 0)
	for _, v := range users {
		// 查询用户的头像路径
		avatarPath := ""
		userImgs, err := controllers.GetImagesControllers("用户", v.ID)
		if err != nil { // 数据库中没有该用户的头像，使用默认的头像
			avatarPath = internal_utils.UserDefaultImage
		} else {
			avatarPath = (*userImgs)[0].Path
		}

		// 查询这个用户的全部角色id
		nowRoleIds := repositories.QueryAdminRoleByUserId(u.DB, v.ID)
		// 查询用户拥有的 role name
		roles := make([]requests.Role, 0)
		for _, roleId := range nowRoleIds {
			role := repositories.QueryRoleById(u.DB, roleId)
			if role == nil {
				return nil, 0, fmt.Errorf("UserReqContext.List() err: 不存在id为 %v 的角色", roleId)
			}
			roles = append(roles, requests.Role{Id: roleId, Name: role.Name})
		}

		listRes = append(listRes, &requests.ListRes{
			Id:         v.ID,
			AvatarPath: avatarPath,
			NickName:   v.Nickname,
			Email:      v.Email,
			Heat:       v.Heat,
			FansCount:  v.FansCount,
			Roles:      roles,
			UserStatus: v.Status,
			// 将 time.Time 格式化为字符串
			LastLoginTime: v.LastLoginTime.Format("2006-01-02 15:04:05"),
			CreateTime:    v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return listRes, total, nil
}

//
// // Import 导入用户表
// func (u *UserReqContext) Import(file *multipart.FileHeader) error {
// 	// 打开上传的文件
// 	openFile, err := file.Open()
// 	if err != nil {
// 		return fmt.Errorf("UserReqContext.Import() err: Failed to open file")
// 	}
// 	defer openFile.Close()
// 	// 解析 Excel 文件
// 	f, err := excelize.OpenReader(openFile)
// 	if err != nil {
// 		return fmt.Errorf("UserReqContext.Import() err: Failed to read Excel file")
// 	}
// 	// 获取第一个工作表的所有行
// 	rows, err := f.GetRows("User Template")
// 	if err != nil {
// 		return fmt.Errorf("UserReqContext.Import() err: Failed to get rows from Excel")
// 	}
//
// 	var l = 0 // 表头长度
// 	layout := "1/2/06 15:04"
// 	loc, _ := time.LoadLocation("Asia/Shanghai") // 使用上海时区 (UTC+8)，否则存入到数据库中的数据是世界时区
// 	// 从第二行开始读取数据
// 	for i, row := range rows {
// 		// 跳过表头
// 		if i == 0 {
// 			l = len(row)
// 			continue
// 		}
// 		// 表格第i行后面的一些数据是nil，数据不够
// 		if len(row) < l {
// 			return fmt.Errorf("UserReqContext.Import() err: 第 %v 行数据，数据量不足", i)
// 		}
// 		// 根据列顺序解析每一行数据
// 		id := 0
// 		if n, err := fmt.Sscanf(row[0], "%d", &id); n == 0 || err != nil {
// 			return fmt.Errorf("UserReqContext.Import() err: Id为 %v 不合法", row[0])
// 		}
// 		// createAt 可以为空
// 		var createAt time.Time
// 		if row[1] != "" {
// 			createAt, err = time.ParseInLocation(layout, row[1], loc)
// 			if err != nil {
// 				return fmt.Errorf("UserReqContext.Import() err: CreateAt为 %v 不合法", row[1])
// 			}
// 		} else {
// 			createAt = time.Now()
// 		}
// 		var updateAt time.Time
// 		if row[2] != "" {
// 			// updateAt, err = time.Parse("2006-01-02", row[2])
// 			// updateAt, err = time.Parse("1/2/06 15:04", row[2])
// 			updateAt, err = time.ParseInLocation(layout, row[2], loc)
// 			if err != nil {
// 				return fmt.Errorf("UserReqContext.Import() err: UpdateAt为 %v 不合法", row[2])
// 			}
// 		} else {
// 			updateAt = time.Now()
// 		}
// 		nickname := row[3]
// 		email := row[4]
// 		password := row[5]
// 		heat := 0
// 		if n, err := fmt.Sscanf(row[6], "%d", &heat); n == 0 || err != nil {
// 			return fmt.Errorf("UserReqContext.Import() err: heat为 %v 不合法", row[6])
// 		}
// 		attentionCount := 0
// 		if n, err := fmt.Sscanf(row[7], "%d", &attentionCount); n == 0 || err != nil {
// 			return fmt.Errorf("UserReqContext.Import() err: AttentionCount为 %v 不合法", row[7])
// 		}
// 		fansCount := 0
// 		if n, err := fmt.Sscanf(row[8], "%d", &fansCount); n == 0 || err != nil {
// 			return fmt.Errorf("UserReqContext.Import() err: FansCount为 %v 不合法", row[8])
// 		}
// 		status := 0
// 		if n, err := fmt.Sscanf(row[9], "%d", &status); n == 0 || err != nil {
// 			return fmt.Errorf("UserReqContext.Import() err: Status为 %v 不合法", row[9])
// 		}
// 		// lastLoginTime 可以为空
// 		var lastLoginTime time.Time
// 		if row[10] != "" {
// 			lastLoginTime, err = time.ParseInLocation(layout, row[10], loc)
// 			if err != nil {
// 				return fmt.Errorf("UserReqContext.Import() err: LastLoginTime为 %v 不合法", row[10])
// 			}
// 		}
//
// 		// 判断数据
// 		// 判断id是否已经存在
// 		if user := repositories.QueryUserByIdIncludeSoftDelete(u.DB, uint(id)); user != nil {
// 			return fmt.Errorf("UserReqContext.Import() err: id为 %v 已经被使用", id)
// 		}
// 		// 判断邮箱是否合法，是否已经使用
// 		if !internal_utils.IsValidEmail(email) {
// 			return fmt.Errorf("UserReqContext.Import() err: email为 %v 不合法", email)
// 		}
// 		if user := repositories.QueryUserByEmail(u.DB, email); user != nil {
// 			return fmt.Errorf("UserReqContext.Import() err: email为 %v 已经被使用", email)
// 		}
// 		// 密码是否合法，加密
// 		if !internal_utils.IsValidPassword(password) {
// 			return fmt.Errorf("UserReqContext.Import() err: password为 %v 不合法", password)
// 		}
// 		password, err = internal_utils.HashPassword(password)
// 		if err != nil {
// 			return fmt.Errorf("UserReqContext.Import() err: password为 %v 加密失败", password)
// 		}
// 		// 检查Status，Status只能是1或者2
// 		if status != 1 && status != 2 {
// 			return fmt.Errorf("UserReqContext.Import() err: status为 %v 不合法", status)
// 		}
// 		// 判断 CreateAt <= UpdateAt，CreateAt <= LastLoginTime
// 		if !createAt.Before(updateAt) && !createAt.Equal(updateAt) {
// 			return fmt.Errorf("UserReqContext.Import() err: CreateAt %v 应该早于 UpdateAt %v", createAt, updateAt)
// 		}
// 		if !createAt.Before(lastLoginTime) && !createAt.Equal(lastLoginTime) {
// 			return fmt.Errorf("UserReqContext.Import() err: CreateAt %v 应该早于 LastLoginTime %v", createAt, updateAt)
// 		}
//
// 		// 将数据存储到数据库中
// 		user := &models.User{
// 			Model: gorm.Model{
// 				ID:        uint(id),
// 				CreatedAt: createAt,
// 				UpdatedAt: updateAt,
// 			},
// 			Nickname:       nickname,
// 			Email:          email,
// 			Password:       password,
// 			Heat:           heat,
// 			AttentionCount: uint(attentionCount),
// 			FansCount:      fansCount,
// 			Status:         status,
// 			LastLoginTime:  lastLoginTime,
// 		}
// 		// 保存到数据库
// 		if err = repositories.InsertObject(u.DB, user); err != nil {
// 			return fmt.Errorf("UserReqContext.Import() err: Failed to save user to database")
// 		}
// 	}
//
// 	return nil
// }
//
// // Export 导出用户表
// func (u *UserReqContext) Export() error {
// 	// 从数据库获取所有用户
// 	users, err := repositories.QueryAllUser(u.DB)
// 	if err != nil {
// 		return fmt.Errorf("UserReqContext.Export() -> %v", err)
// 	}
//
// 	// 创建一个新的 Excel 文件
// 	f := excelize.NewFile()
// 	// 创建一个新的 Sheet，工作表名称为 "Users"
// 	index, _ := f.NewSheet("Users")
// 	// 添加列头到第一行
// 	headers := []string{"Id", "CreateAt", "UpdateAt", "Nickname", "Email", "Password", "Heat", "AttentionCount", "FansCount", "Status", "LastLoginTime"}
// 	for i, header := range headers {
// 		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
// 		f.SetCellValue("Users", cell, header)
// 	}
//
// 	// 写入用户数据
// 	for i, user := range users {
// 		row := i + 2 // 从第二行开始写入数据
// 		f.SetCellValue("Users", fmt.Sprintf("A%d", row), user.ID)
// 		f.SetCellValue("Users", fmt.Sprintf("B%d", row), user.CreatedAt.Format("2006-01-02 15:04:05"))
// 		f.SetCellValue("Users", fmt.Sprintf("C%d", row), user.UpdatedAt.Format("2006-01-02 15:04:05"))
// 		f.SetCellValue("Users", fmt.Sprintf("D%d", row), user.Nickname)
// 		f.SetCellValue("Users", fmt.Sprintf("E%d", row), user.Email)
// 		f.SetCellValue("Users", fmt.Sprintf("F%d", row), user.Password)
// 		f.SetCellValue("Users", fmt.Sprintf("G%d", row), user.Heat)
// 		f.SetCellValue("Users", fmt.Sprintf("H%d", row), user.AttentionCount)
// 		f.SetCellValue("Users", fmt.Sprintf("I%d", row), user.FansCount)
// 		f.SetCellValue("Users", fmt.Sprintf("J%d", row), user.Status)
// 		f.SetCellValue("Users", fmt.Sprintf("K%d", row), user.LastLoginTime.Format("2006-01-02 15:04:05"))
// 	}
//
// 	// 设置活动工作表
// 	f.SetActiveSheet(index)
//
// 	// // 设置响应头，返回 Excel 文件
// 	// c.Header("Content-Disposition", "attachment; filename=users.xlsx")
// 	// c.Header("Content-Type", "application/octet-stream")
// 	// c.Header("Content-Transfer-Encoding", "binary")
//
// 	// 将文件内容写入响应中
// 	if err := f.Write(u.Ctx.Writer); err != nil {
// 		return fmt.Errorf("UserReqContext.Export() err: Failed to create Excel file")
// 	}
//
// 	return nil
// }
//
// // DownloadTemplate 下载导入用户模版excel
// func (u *UserReqContext) DownloadTemplate() error {
// 	// 创建一个新的 Excel 文件
// 	f := excelize.NewFile()
// 	// 创建一个新的 Sheet，工作表名称为 "User Template"
// 	index, _ := f.NewSheet("User Template")
//
// 	// 添加列头到第一行
// 	// headers := []string{"Id", "Nickname", "Email", "Password", "Heat", "AttentionCount", "FansCount", "PrivateSettings", "Status", "LastLoginTime"}
// 	headers := []string{"Id", "CreateAt", "UpdateAt", "Nickname", "Email", "Password", "Heat", "AttentionCount", "FansCount", "Status", "LastLoginTime"}
//
// 	for i, header := range headers {
// 		cell := fmt.Sprintf("%s1", string(rune('A'+i))) // 将索引转换为对应的列号 (A, B, C...)
// 		// 设置单元格值
// 		if err := f.SetCellValue("User Template", cell, header); err != nil {
// 			return fmt.Errorf("UserReqContext.DownloadTemplate() err: %v", err)
// 		}
// 	}
//
// 	// 将 "User Template" 设为活动工作表
// 	f.SetActiveSheet(index)
//
// 	// // 设置响应头，返回 Excel 文件
// 	// c.Header("Content-Disposition", "attachment; filename=user_template.xlsx")
// 	// c.Header("Content-Type", "application/octet-stream")
// 	// c.Header("Content-Transfer-Encoding", "binary")
//
// 	// 将文件内容写入响应中
// 	if err := f.Write(u.Ctx.Writer); err != nil {
// 		return fmt.Errorf("UserReqContext.DownloadTemplate() err: Failed to create Excel file")
// 	}
// 	return nil
// }

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
		if !internal_utils.IsValidEmail(email) {
			return fmt.Errorf("UserReqContext.Import() err: email为 %v 不合法", email)
		}
		if user := repositories.QueryUserByEmail(u.DB, email); user != nil {
			return fmt.Errorf("UserReqContext.Import() err: email为 %v 已经被使用", email)
		}
		// 密码是否合法，加密
		if !internal_utils.IsValidPassword(password) {
			return fmt.Errorf("UserReqContext.Import() err: password为 %v 不合法", password)
		}
		password, err = internal_utils.HashPassword(password)
		if err != nil {
			return fmt.Errorf("UserReqContext.Import() err: password为 %v 加密失败", password)
		}
		// 检查Status，Status只能是1或者2
		if status != 1 && status != 2 {
			return fmt.Errorf("UserReqContext.Import() err: status为 %v 不合法", status)
		}

		// 将数据存储到数据库中
		user := &models.User{
			Nickname: nickname,
			Email:    email,
			Password: password,
			Status:   status,
		}
		// 保存到数据库
		if err = repositories.InsertObject(u.DB, user); err != nil {
			return fmt.Errorf("UserReqContext.Import() err: Failed to save user to database")
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

	// 将文件内容写入响应中
	if err := f.Write(u.Ctx.Writer); err != nil {
		return fmt.Errorf("UserReqContext.Export() err: Failed to create Excel file")
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
	avatarPath := ""
	userImgs, err := controllers.GetImagesControllers("用户", user.ID)
	if err != nil { // 数据库中没有该用户的头像，使用默认的头像
		avatarPath = internal_utils.UserDefaultImage
	} else {
		avatarPath = (*userImgs)[0].Path
	}

	// 查询这个用户的全部角色
	roleIds := repositories.QueryAdminRoleByUserId(u.DB, user.ID)
	roles := make([]requests.Role, 0)
	for _, roleId := range roleIds {
		role := repositories.QueryRoleById(u.DB, roleId)
		if role == nil {
			return nil, fmt.Errorf("UserReqContext.GetInfo() err: 不存在id为 %v 的角色", roleId)
		}
		roles = append(roles, requests.Role{Id: roleId, Name: role.Name})
	}

	var getInfoRes = &requests.GetInfoRes{
		Id:         user.ID,
		AvatarPath: avatarPath,
		NickName:   user.Nickname,
		Email:      user.Email,
		UserStatus: user.Status,
		Roles:      roles,
	}

	return getInfoRes, nil
}

// UploadHeadshot 上传用户头像
func (u *UserReqContext) UploadHeadshot(id uint) error {
	// 判断id是否存在
	user := repositories.QueryUserById(u.DB, id)
	if user == nil {
		return fmt.Errorf("UserReqContext.UploadHeadshot() err: 不存在id为 %v 的用户", id)
	}

	// 存头像
	err, _ := controllers.UploadImagesControllers(u.Ctx, "用户", id)
	if err != nil {
		return fmt.Errorf("UserReqContext.UploadHeadshot() 存头像错误err: %v", err)
	}
	return nil
}