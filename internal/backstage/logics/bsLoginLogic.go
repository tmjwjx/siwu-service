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
	"forum/pkg/token"
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
func (b *BsManageContext) BsLogin(msg requests.BackstageLoginReq) (*requests.FinalBackstageLoginRes, error) {
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

	// 查询超级管理员的ID
	superAdminId, err := repositories.QuerySuperAdminId(b.DB)
	if err != nil {
		return nil, fmt.Errorf("UserReqContext.BsLogin() %v", err)
	}

	var menuPerm *[]requests.MenuPerm // 用于存储当前用户角色的所有菜单权限并集
	var permCode *[]string            // 用于存储当前用户角色的权限标识并集

	var flag int // 用于判断该用户是否是超级管理员

	// 判断该用户是否是超级管理员
	for _, roleId := range roleIds {
		if roleId == superAdminId {
			flag = 1
			break
		}
	}

	if flag == 1 {
		menuPerm, err = b.GetMenuPermRep(b.DB, nil, 1)
		if err != nil {
			return nil, fmt.Errorf("UserReqContext.BsLogin() -> menuPerm -> %v", err)
		}
		permCode, err = b.GetPermCodeRep(casbinService, b.DB, nil, 1)
		if err != nil {
			return nil, fmt.Errorf("UserReqContext.BsLogin() -> permCode -> %v", err)
		}

	} else {
		menuPerm, err = b.GetMenuPermRep(b.DB, roleIds, 0)
		if err != nil {
			return nil, fmt.Errorf("UserReqContext.BsLogin()2 %v", err)
		}
		permCode, err = b.GetPermCodeRep(casbinService, b.DB, roleIds, 0)
		if err != nil {
			return nil, fmt.Errorf("UserReqContext.BsLogin() -> permCode -> %v", err)
		}
	}

	res := &requests.FinalBackstageLoginRes{
		UserInfo: backstageLoginRes,
		Perm:     menuPerm,
		CodeList: permCode,
	}

	// return backstageLoginRes, nil

	// 补丁：只允许管理者进入后台（以后要额外添加一张表用来存储管理者）
	for _, v := range backstageLoginRes.RoleNames {
		if v == "管理员" || v == "超级管理员" || v == "测试" {
			return res, nil
		}
	}
	return nil, fmt.Errorf("用户id为%d的用户没有权限进入后台", user.ID)
}

// BsLogout 后台登出
func (b *BsManageContext) BsLogout(tokenString string) error {

	// 设置过期时间为 Token 剩余时间
	claims, err := token.ValidateToken(tokenString)
	if err != nil {
		return fmt.Errorf("BsManageContext.BsLogout() : 无效的 token")
		// response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("Logout() : 无效的 token"), nil))
		// return
	}

	expiration := time.Until(claims.ExpiresAt.Time)
	// 如果该token还没有失效，就把它添加到黑名单中，让它失效
	if expiration > 0 {
		if err = token.AddTokenToBlacklist(globals.RDB, tokenString, expiration); err != nil {
			return fmt.Errorf("BsManageContext.BsLogout() : err -> %v", err)
			// response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("Logout() : err -> %v", err), nil))
			// return
		}
	}
	return nil
}

// GetMenuPermRep
// @Description: 获取当前角色的菜单权限(用于渲染侧边栏，只要type1和2)
// @Author wangyulong 2025-02-07 20:29:59
// @receiver     b
// @param        db *gorm.DB
// @param        roleIds []uint
// @param        flag int 用于判断该用户是否是超级管理员
// @return       *requests.GetMenuPermRes
// @return       error
func (b *BsManageContext) GetMenuPermRep(db *gorm.DB, roleIds []uint, flag int) (*[]requests.MenuPerm, error) {

	var menuIds []uint

	// flag 为 1 表示该用户是超级管理员
	if flag == 1 {
		// 查询所有菜单ID
		err := db.Model(&models.Menu{}).Pluck("id", &menuIds).Error
		if err != nil {
			return nil, fmt.Errorf("GetMenuPermRep -> 1查询菜单id异常 -> %s", err)
		}
	} else {
		// flag 为 0 表示该用户不是超级管理员

		// seen 用于高效地去除切片中的重复数据
		seen := make(map[uint]struct{}) // 使用空 struct{} 节省内存

		for _, roleId := range roleIds {
			menuId := make([]uint, 0)
			// 查询角色拥有的菜单的ID
			err := b.DB.Model(&models.RoleMenu{}).Select("MenuId").Where("role_id = ?", roleId).Scan(&menuId).Error
			if err != nil {
				return nil, fmt.Errorf("GetMenuPermRep -> 2查询菜单id异常 -> %s", err)
			}
			// 将角色拥有的菜单ID去重
			internalUtils.RemoveDuplicates(&seen, menuId)
		}
		for menuId, _ := range seen {
			menuIds = append(menuIds, menuId)
		}
	}

	// 查询权限
	var menuPerm []requests.MenuPerm
	err := db.Model(models.Menu{}).Where("id IN ? and type IN ?", menuIds, []uint{1, 2}).Scan(&menuPerm).Error
	if err != nil {
		return nil, fmt.Errorf("GetMenuPermRep -> 查询权限失败 -> %s", err)
	}

	var menuPermIds []uint
	var menuParent []uint

	for _, menu := range menuPerm {
		// 将查询出来的符合条件的菜单的ID存入menuPermIds中
		menuPermIds = append(menuPermIds, menu.ID)

		if menu.ParentId != 0 {
			// 将查询出来的符合条件的菜单的父ID不为0的父菜单的ID存入menuParent中
			menuParent = append(menuParent, menu.ParentId)
		}

	}

	// seen 用于高效地去除切片中的重复数据
	seen := make(map[uint]struct{}) // 使用空 struct{} 节省内存
	// 将menuIds中的元素现存入seen中
	internalUtils.RemoveDuplicates(&seen, menuPermIds)

	var menuParent2 []uint
	// 从menuParent中提取出menuIds中没有的元素
	menuParent2 = internalUtils.RemoveDuplicates2(&seen, menuParent)

	// 查询父ID不为0的菜单的信息
	err = db.Model(models.Menu{}).Where("id IN ? and type IN ?", menuParent2, []uint{1, 2}).Scan(&menuPerm).Error
	if err != nil {
		return nil, fmt.Errorf("GetMenuPermRep -> 查询权限失败 -> %s", err)
	}

	return &menuPerm, nil
}

// GetPermCodeRep
// @Description: 获取当前角色的所有权限标识
// @Author wangyulong 2025-02-07 22:57:51
// @receiver     b
// @param        casbinService *casbin.CasbinService
// @param        db *gorm.DB
// @param        roleIds []uint
// @param        flag int 用于判断该用户是否是超级管理员
// @return       *[]string
// @return       error
func (b *BsManageContext) GetPermCodeRep(casbinService *casbin.CasbinService, db *gorm.DB, roleIds []uint, flag int) (*[]string, error) {

	var apiIds []uint

	// flag 为 1 表示该用户是超级管理员
	if flag == 1 {
		// 获取所有api权限
		err := db.Model(&models.Api{}).Pluck("id", &apiIds).Error
		if err != nil {
			return nil, fmt.Errorf("GetPermCodeRep -> 获取当前角色的api权限失败 -> %s", err)
		}
	} else {
		// flag 为 0 表示该用户不是超级管理员

		// seen 用于高效地去除切片中的重复数据
		seen := make(map[uint]struct{}) // 使用空 struct{} 节省内存

		for _, roleId := range roleIds {
			// 获取当前角色的api权限
			apiId, err := casbinService.GetApiPerm(fmt.Sprintf("%v", roleId))
			if err != nil {
				return nil, fmt.Errorf("GetPermCodeRep -> 获取当前角色的api权限失败 -> %s", err)
			}
			// 将角色拥有的菜单ID去重
			internalUtils.RemoveDuplicates(&seen, apiId)
		}
		for apiId, _ := range seen {
			apiIds = append(apiIds, apiId)
		}
	}

	// 获取type为3的菜单(即按钮)
	var menuIds []uint
	err := db.Model(&models.Menu{}).Where("type = ?", 3).Pluck("id", &menuIds).Error
	if err != nil {
		return nil, fmt.Errorf("GetPermCodeRep -> 获取type为3的菜单(即按钮)失败 -> %s", err)
	}

	// 通过菜单id查询其拥有的api的id
	//menuApi := make(map[uint]uint)
	var endMenuId []uint
	var apiID []uint
	for _, menuId := range menuIds {
		err = db.Model(&models.MenuApi{}).Where("menu_id = ?", menuId).Pluck("api_id", &apiID).Error
		if err != nil {
			return nil, fmt.Errorf("GetPermCodeRep -> 通过菜单id查询其拥有的api的id失败 -> %s", err)
		}

		// 通过查询到的菜单拥有的api,去角色拥有的api权限中寻找相等的api，如果相等，这该菜单(即按钮)属于该角色
		for _, mApi := range apiID {
			flag2 := 0
			for _, rApi := range apiIds {
				if rApi == mApi {
					endMenuId = append(endMenuId, menuId)
					flag2 = 1
					break
				}
			}
			if flag2 == 1 {
				break
			}
		}
	}

	//通过查到的角色拥有的菜单(即按钮),去查询菜单的code字段
	var codes []string
	err = db.Model(&models.Menu{}).Where("id IN ?", endMenuId).Pluck("code", &codes).Error
	if err != nil {
		return nil, fmt.Errorf("GetPermCodeRep -> 通过菜单id查询其拥有的api的id失败 -> %s", err)
	}

	return &codes, nil
}
