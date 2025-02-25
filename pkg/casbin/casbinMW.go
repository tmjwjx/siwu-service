package casbin

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/pkg/globals"
	"forum/pkg/response"
	"forum/pkg/utils"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func CasbinAuth(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 加载策略文件
		err := e.LoadPolicy()
		if err != nil {
			e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("NewCasbinAuth -> 加载策略文件失败 -> %s", err), nil)
			response.Failed(c, 500, e)
			c.Abort()
			return
		}

		// 获取管理员id
		adminID, exists := c.Get("id")
		if !exists {
			e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("NewCasbinAuth -> 获取用户id失败"), nil)
			response.Failed(c, 500, e)
			c.Abort()
			return
		}

		// 查询管理员的email
		uintValue, err := internalUtils.ChangeAnyToUint(adminID)
		if err != nil {
			e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
			response.Failed(c, 500, e)
			c.Abort()
			return
		}

		email, err := utils.SelEmailForAdmin(uintValue)
		if err != nil {
			e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
			response.Failed(c, 500, e)
			c.Abort()
			return
		}
		//权限验证

		/*// 查询超级管理员对应的id
		superAdminID, err := utils.SelIdForSuperAdmin()
		if superAdminID == "" && err != nil {
			e := response.NewAppErr(globals.StatusBadRequest, err, nil)
			response.Failed(c, 400, e)
			c.Abort()
			return
		}*/

		// 验证用户是否是超级管理员
		casbinServer := &CasbinService{
			Enforcer: e,
		}

		ok, err := casbinServer.VerifySuperAdministrator(email, "超级管理员")
		if err != nil {
			e := response.NewAppErr(globals.StatusBadRequest, err, nil)
			response.Failed(c, 400, e)
			c.Abort()
			return
		}
		if ok {
			// 如果用户是超级管理员，就不用进行权限验证了，直接结束即可
			c.Next()
		} else {
			// 获取apiId
			apiId, err := utils.SelApiId(c.Request.URL.Path)
			if err != nil {
				e := response.NewAppErr(globals.StatusBadRequest, err, nil)
				response.Failed(c, 400, e)
				c.Abort()
				return
			}

			ok, err = e.Enforce(email, fmt.Sprintf("%v", apiId))
			if err != nil {
				e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("NewCasbinAuth -> 权限验证失败"), nil)
				response.Failed(c, 500, e)
				c.Abort()
				return
			} else if !ok {
				e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("NewCasbinAuth -> 该用户没有该权限"), nil)
				response.Failed(c, 400, e)
				c.Abort()
				return
			}
			c.Next()
		}
	}
}
