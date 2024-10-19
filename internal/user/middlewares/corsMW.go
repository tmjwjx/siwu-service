package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
