package casbin

import (
	"fmt"
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
			e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("NewCasbinAuth -> 加载策略文件失败"), nil)
			response.Failed(c, 500, e)
			c.Abort()
			return
		}

		// 获取用户id
		id, exists := c.Get("id")
		fmt.Println("---------------->", id)
		if !exists {
			e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("NewCasbinAuth -> 获取用户id失败"), nil)
			response.Failed(c, 500, e)
			c.Abort()
			return
		}

		//权限验证

		// 获取apiId
		apiId, err := utils.SelApiId(c.Request.URL.Path)
		fmt.Println("=================>", apiId)
		if err != nil {
			e := response.NewAppErr(globals.StatusBadRequest, err, nil)
			response.Failed(c, 400, e)
			c.Abort()
			return
		}
		fmt.Println("+++++++++---->", fmt.Sprintf("%v", id), fmt.Sprintf("%v", apiId))
		//ok, err := cbs.Enforcer.Enforce(id, apiId)
		//ok, err := cbs.Enforcer.Enforce(fmt.Sprintf("%v", id), fmt.Sprintf("%v", apiId))
		ok, err2 := e.Enforce("1", "16")
		fmt.Println(err2)
		if err2 != nil {
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
