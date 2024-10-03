package casbin

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// casbin 结构体
type CasbinService struct {
	Enforcer *casbin.Enforcer
	Adapter  *gormadapter.Adapter
}

//// RolePolicy 对应于 'CasbinRule' 表中的(v0, v1)
//type RolePolicy struct {
//	RoleId   uint   // 角色id
//	SingleId uint   // 菜单id或者api的id
//	Kind     string // 用于判断是菜单还是api
//}
