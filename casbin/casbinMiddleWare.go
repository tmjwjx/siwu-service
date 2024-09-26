package casbin

//import (
//	"fmt"
//	"forum/pkg/globals"
//	"forum/pkg/response"
//	"github.com/gin-gonic/gin"
//)
//
//func NewCasbinAuth(cbs *CasbinService) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 加载策略文件
//		err := cbs.Enforcer.LoadPolicy()
//		if err != nil {
//			e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("NewCasbinAuth -> 加载策略文件失败"), nil)
//			response.Failed(c, 500, e)
//			c.Abort()
//			return
//		}
//		// 获取请求头token解析出username
//		token := c.GetHeader("token")
//		// 查询账号对应的角色
//		role, err := token2.GetRole(token)
//		if err != nil || role == "" {
//			e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("NewCasbinAuth -> token中没有携带角色"), nil)
//			response.Failed(c, 400, e)
//			c.Abort()
//			return
//		}
//		// 该部分写菜单
//		menuId, err := SelMenId(c.Request.URL.Path, c.Request.Method)
//		if err != nil || menuId == "" {
//			e := response.NewAppErr(globals.StatusBadRequest, err, nil)
//			response.Failed(c, 400, e)
//			c.Abort()
//			return
//		}
//		// 权限验证
//		ok, err := cbs.Enforcer.Enforce(role, menuId)
//		if err != nil {
//			e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("NewCasbinAuth -> 权限验证失败"), nil)
//			response.Failed(c, 500, e)
//			c.Abort()
//			return
//		} else if !ok {
//			e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("NewCasbinAuth -> 该用户没有该权限"), nil)
//			response.Failed(c, 400, e)
//			c.Abort()
//			return
//		}
//		c.Next()
//	}
//}
