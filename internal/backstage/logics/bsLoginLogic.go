package logics

import (
	"fmt"
	"forum/internal/backstage/repositories"
	"forum/internal/backstage/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/internalPkg/sqlUtils"
	"forum/internal/models"
	"forum/pkg/casbin"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"time"
)

// BsManageContext
// @Description: 用于在处理请求时传递数据库连接和请求上下文信息
// @Author lizhuang 2024-10-21 20:26:10
type BsManageContext struct {
	DB  *gorm.DB
	Ctx *gin.Context
}

// NewBsManageContext
// @Description: 新建BsManageContext对象
// @Author lizhuang 2024-10-21 20:26:15
// @param        db *gorm.DB
// @param        c *gin.Context
// @return       *BsManageContext
func NewBsManageContext(db *gorm.DB, c *gin.Context) *BsManageContext {
	return &BsManageContext{
		DB:  db,
		Ctx: c,
	}
}

// BsLogin
// @Description: 后台登陆业务
// @Author lizhuang 2024-10-21 20:36:41
// @receiver     b
// @param        msg requests.BackstageLoginReq
// @return       error
func (b *BsManageContext) BsLogin(msg requests.BackstageLoginReq) (*requests.BackstageLoginRes, error) {
	// 判断邮箱和密码是否匹配
	email := msg.Email
	password := msg.Password

	// 根据邮箱查用户
	user := repositories.QueryUserByEmail(b.DB, email)
	if user == nil {
		return nil, fmt.Errorf("BsManageContext.BsLogin() err: 不存在该邮箱用户")
	}

	// 比较加密密码
	encryptedPassword := user.Password
	if !internalUtils.CheckPasswordHash(password, encryptedPassword) {
		return nil, fmt.Errorf("BsManageContext.BsLogin() err: 密码错误")
	}

	// 改变 LastLoginTime
	now := time.Now() // 获取当前时间
	if err := sqlUtils.UpdateObjects(b.DB, &models.User{Model: gorm.Model{ID: user.ID}}, map[string]interface{}{"last_login_time": now}); err != nil {
		return nil, fmt.Errorf("BsManageContext.BsLogin() -> %v", err)
	}

	// 查询用户的头像路径
	userImages, err := internalUtils.GetImages(b.DB, globals.UserHome, user.ID)
	if err != nil {
		return nil, fmt.Errorf("BsManageContext.BsLogin() -> %v", err)
	}
	// 没有图片
	if userImages == nil {
		return nil, fmt.Errorf("BsManageContext.BsLogin() err = 无法找到id为%d的用户头像图片", user.ID)
	}
	avatarPath := (*userImages)[0]

	var roleNames = make([]string, 0)
	sort := math.MaxInt
	code := ""

	// 获取该用户对应的全部角色id
	casbinService, err := casbin.NewCasbinService(globals.DB)
	if err != nil {
		return nil, fmt.Errorf("UserReqContext.BsLogin() %v", err)
	}
	roleIds, err := casbinService.GetRolesForUser(user.ID)
	if err != nil {
		return nil, fmt.Errorf("UserReqContext.BsLogin() %v", err)
	}

	for _, v := range roleIds {
		// 查询角色id对应的角色信息
		role := repositories.QueryRoleById(b.DB, v)
		roleNames = append(roleNames, role.Name)

		// 获取当前用户所拥有的排序最高的角色编码
		if sort > role.Sort {
			sort = role.Sort
			code = role.Code
		}
	}

	var backstageLoginRes = &requests.BackstageLoginRes{
		Id:         user.ID,
		Nickname:   user.Nickname,
		Email:      user.Email,
		UserStatus: user.Status,
		RoleIds:    roleIds,
		RoleNames:  roleNames,
		AvatarPath: avatarPath,
		Code:       code,
	}

	// return backstageLoginRes, nil

	// 补丁：只允许管理者进入后台
	for _, v := range backstageLoginRes.RoleNames {
		if v == "管理员" || v == "超级管理员" {
			return backstageLoginRes, nil
		}
	}
	return nil, fmt.Errorf("用户id为%d的用户没有权限进入后台", user.ID)
}

// // BsLogout 后台登出
// func (b *BsManageContext) BsLogout(tokenString string) error {
// 	// 使token无效
// 	if err := token.InvalidateToken(tokenString); err != nil {
// 		return fmt.Errorf("BsManageContext.BsLogout() err = %v", err)
// 	}
// 	return nil
// }
