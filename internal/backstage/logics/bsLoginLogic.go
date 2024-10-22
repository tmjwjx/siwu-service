package logics

import (
	"fmt"
	"forum/internal/backstage/repositories"
	"forum/internal/backstage/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/internalPkg/sqlUtils"
	"forum/internal/models"
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

	// 获取用户头像
	avatarPath := ""
	//userImgs, err := imageCtrl.GetImagesControllers("user", user.ID)
	//if err != nil { // 数据库中没有该用户的头像，使用默认的头像
	//	avatarPath = internalUtils.UserDefaultImage
	//} else {
	//	avatarPath = (*userImgs)[0].Path
	//}

	var roleNames = make([]string, 0)
	sort := math.MaxInt
	code := ""
	// 获取用户拥有的角色id
	roleIds := repositories.QueryAdminRoleByUserId(b.DB, user.ID)
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
		RoleNames:  roleNames,
		AvatarPath: avatarPath,
		Code:       code,
	}
	return backstageLoginRes, nil
}
