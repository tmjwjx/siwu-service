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
		if tokenString == "" || len(tokenString) <= len("Bearer ") {
			response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("AuthMiddleware() : 缺少或无效的 Authorization 头"), nil))
			c.Abort()
			return
		}
		tokenString = tokenString[len("Bearer "):]

		//  检查 Token 是否在黑名单
		if IsTokenBlacklisted(globals.RDB, tokenString) {
			response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("AuthMiddleware() : token 已失效"), nil))
			c.Abort()
			return
		}

		// 验证 Token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("AuthMiddleware() : 无效的 token"), nil))
			c.Abort()
			return
		}

		// 保存用户 ID 到上下文
		c.Set("id", claims.ID)
		c.Next()
	}
}
