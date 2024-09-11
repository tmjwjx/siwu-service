package middlewares

import (
	"fmt"
	"forum/pkg/globals"
	"forum/pkg/response"
	"forum/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

// // JWTAuth 中间件，验证请求
// func JWTAuth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		token := c.GetHeader("Authorization")
// 		if token == "" {
// 			response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("令牌为空"), nil))
// 			c.Abort()
// 			return
// 		}
//
// 		user, err := token.ParseToken(token)
// 		if err != nil {
// 			response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("无效的令牌"), nil))
// 			c.Abort()
// 			return
// 		}
//
// 		// 将一个键值对存储在 gin.Context 中，供同一请求的后续处理中使用。
// 		c.Set("user", user)
// 		c.Next()
// 	}
// }

// AuthMiddleware JWT 认证中间件。在 JWT 验证通过后，将 email 存储在上下文中，供后续路由使用。
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
		claims, err := token.ValidateToken(tokenString)
		fmt.Println("AuthMiddleware  ", claims)

		if err != nil {
			response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("AuthMiddleware() : 无效的 token"), nil))
			c.Abort()
			return
		}

		// 将用户 email 和 password 信息保存到上下文中
		c.Set("email", claims.Email)

		c.Next()
	}
}
