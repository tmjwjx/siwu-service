package middlewares

import (
	"fmt"
	"forum/internal/internal_pkg/internal_utils"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// JWTAuth 中间件，验证请求
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("令牌为空"), nil))
			c.Abort()
			return
		}

		user, err := internal_utils.ParseToken(token)
		if err != nil {
			response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("无效的令牌"), nil))
			c.Abort()
			return
		}

		// 将一个键值对存储在 gin.Context 中，供同一请求的后续处理中使用。
		c.Set("user", user)
		c.Next()
	}
}
