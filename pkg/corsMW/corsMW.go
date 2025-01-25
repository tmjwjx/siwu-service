package corsMW

import (
	"github.com/gin-gonic/gin"
)

// // CorsMiddleware 跨域中间件
// func CorsMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// 获取当前请求的方法
// 		method := c.Request.Method
// 		// 获取请求头中的 Origin 字段
// 		origin := c.Request.Header.Get("Origin")
// 		if origin != "" {
// 			c.Header("Access-Control-Allow-Origin", "*") // 可将 * 替换为指定的域名
// 			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
// 			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
// 			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
// 			c.Header("Access-Control-Allow-Credentials", "true")
// 		}
// 		// 在 预检请求（Preflight Request）中，浏览器会向服务器发送一个 OPTIONS 请求，询问服务器是否允许跨域请求。这个请求本身不会包含实际的数据，只是浏览器询问允许哪些 HTTP 方法、请求头等。
// 		// 如果请求方法是 OPTIONS，我们通过 AbortWithStatus(http.StatusNoContent) 直接响应并结束该请求的处理，返回 HTTP 状态码 204 No Content，表示没有内容，响应结束。
// 		if method == "OPTIONS" {
// 			c.AbortWithStatus(http.StatusNoContent)
// 		}
// 		c.Next()
// 	}
// }

// CorsMiddleware 跨域中间件
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") // 获取请求来源

		// 允许的前端域名列表
		allowedOrigins := map[string]bool{
			"http://192.168.10.7:9901": true,
			"http://192.168.10.7:9902": true,
		}

		// 判断 origin 是否在允许列表内
		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin) // 仅允许特定域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true") // 允许跨域请求携带 Cookie
		}

		// 处理预检请求（OPTIONS 请求）
		if method == "OPTIONS" {
			c.AbortWithStatus(204) // 直接返回 204 状态码，表示接受预检请求
			return
		}

		c.Next()
	}
}
