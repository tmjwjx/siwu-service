package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware 设置 CORS 中间件。跨域问题。
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},                                // AllowOrigins: 允许的来源，可以设置为特定的域名，也可以使用 * 表示允许所有来源（不建议在生产环境中使用）。（http://192.168.10.7:8081）
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},     // AllowMethods: 允许的 HTTP 方法。
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept"}, // AllowHeaders: 允许的请求头部。
			ExposeHeaders:    []string{"Content-Length"},                   // ExposeHeaders: 允许客户端访问的响应头部。
			AllowCredentials: true,                                         // AllowCredentials: 是否允许发送 Cookies。
			MaxAge:           12 * 3600,                                    // MaxAge: 预检请求的有效期，单位为秒。
		})
	}
}
