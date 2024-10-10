package casbin

import (
	"fmt"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

func NewCasbinAuth(cbs *CasbinService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 加载策略文件
		err := cbs.Enforcer.LoadPolicy()
		if err != nil {
			e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("NewCasbinAuth -> 加载策略文件失败"), nil)
			response.Failed(c, 500, e)
			c.Abort()
			return
		}

		// 获取用户id
		id, exists := c.Get("id")
		if !exists {
			e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("NewCasbinAuth -> 获取用户id失败"), nil)
			response.Failed(c, 500, e)
			c.Abort()
			return
		}

		//权限验证

		// 获取apiId
		apiId, err := SelApiId(c.Request.URL.Path)
		if err != nil || apiId == 0 {
			e := response.NewAppErr(globals.StatusBadRequest, err, nil)
			response.Failed(c, 400, e)
			c.Abort()
			return
		}
		ok, err := cbs.Enforcer.Enforce(id, apiId)
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
