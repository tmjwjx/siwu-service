package token

import (
	"fmt"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware JWT 认证中间件。在 JWT 验证通过后，将 ID 存储在上下文中，供后续路由使用。
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("AuthMiddleware() : 缺少授权标头 Authorization"), nil))
			// 中止剩余的中间件和处理函数执行，直接返回响应
			c.Abort()
			return
		}

		// 提取 Token 部分，去掉 "Bearer " 前缀
		tokenString = tokenString[len("Bearer "):]

		// 验证并解析 Token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("AuthMiddleware() : 无效的 token"), nil))
			c.Abort()
			return
		}

		// 将用户 ID 保存到上下文中
		c.Set("id", claims.ID)
		c.Next()
	}
}
