// package csrfMW
//
// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/csrf"
// 	adapter "github.com/gwatts/gin-adapter"
// 	"net/http"
// )
//
// // 当前端发送 Get 请求的时候会给前端发送一个 X-CSRF-Token 头部，当前端发送 Post、Put、Delete 等需要修改数据的请求的时候，后端会验证 X-CSRF-Token
//
// var csrfMd func(http.Handler) http.Handler
//
// // CSRFTokenMW 在响应中提供 CSRF 令牌给前端。
// // 在每个 HTTP 响应的头部添加 X-CSRF-Token 字段，该字段的值是根据当前请求生成的 CSRF（Cross-Site Request Forgery，跨站请求伪造）令牌。前端可以从响应头部获取这个令牌，并在后续的请求中携带该令牌，以通过后端的 CSRF 验证。
// func CSRFTokenMW() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Header("X-CSRF-Token", csrf.Token(c.Request))
// 	}
// }
//
// // CSRFMW 验证前端携带的 CSRF 令牌的有效性，确保应用免受跨站请求伪造攻击。
// func CSRFMW() gin.HandlerFunc {
// 	// 初始化 CSRF 中间件
// 	csrfMd = csrf.Protect(
// 		[]byte("8d7c2c6a1d4b7a3d9c8e4f6b5a2d1c3e4f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c"),
// 		csrf.Secure(false),  // 是否仅在 HTTPS 连接中使用 CSRF 保护。false 表示在 HTTP 连接中也可以使用。
// 		csrf.HttpOnly(true), // 设置 CSRF 令牌的 HttpOnly 属性为 true，这样 JavaScript 无法访问该令牌，提高安全性。
// 		// 指定当 CSRF 验证失败时的错误处理函数。当 CSRF 令牌无效时，会调用这个函数。
// 		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			w.WriteHeader(http.StatusForbidden)
// 			w.Write([]byte(`{
// 				"code": 4003,
// 				"data": {},
// 				"err": "无效的 CSRF token"
// 			}`))
// 		})),
// 	)
//
// 	// 转换为 Gin 中间件
// 	return adapter.Wrap(csrfMd)
// }

// 可以通过下面的方法生成一个 32 字节的随机十六进制字符串，如 8d7c2c6a1d4b7a3d9c8e4f6b5a2d1c3e4f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c
// keyLength := 32
// key, err := generateRandomKey(keyLength)
//
// func generateRandomKey(length int) (string, error) {
// 	key := make([]byte, length)
// 	_, err := rand.Read(key)
// 	if err != nil {
// 		return "", err
// 	}
// 	return hex.EncodeToString(key), nil
// }

package csrfMW

import (
	"errors"
	"fmt"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"net/http"
	"time"
)

var (
	csrfMd       func(http.Handler) http.Handler
	tokenExpires = 5 * time.Minute // 设置 token 过期时间为 5 分钟
)

// CSRFTokenMW 在响应中提供 CSRF 令牌给前端。
func CSRFTokenMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := csrf.Token(c.Request)
		fmt.Println(token)
		// 将 token 存储到 Redis 中，并设置过期时间
		err := setCSRFTokenToRedis(globals.RDB, token, tokenExpires)
		if err != nil {
			fmt.Printf("Failed to store CSRF token in Redis: %v\n", err)
		}
		c.Header("X-CSRF-Token", token)
	}
}

// CSRFMW 验证前端携带的 CSRF 令牌的有效性，确保应用免受跨站请求伪造攻击。
func CSRFMW() gin.HandlerFunc {
	// 初始化 CSRF 中间件
	csrfMd = csrf.Protect(
		[]byte("8d7c2c6a1d4b7a3d9c8e4f6b5a2d1c3e4f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c"),
		csrf.Secure(false),
		csrf.HttpOnly(true),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{
                "code": 4003,
                "data": {},
                "err": "无效的 CSRF token"
            }`))
		})),
	)

	// 自定义验证逻辑
	return func(c *gin.Context) {
		// 如果是 GET 请求，跳过 CSRF token 检查
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 4003,
				"data": nil,
				"err":  "缺少 CSRF token",
			})
			return
		}

		// 检查 token 是否存在且未使用
		status, err := checkCSRFTokenStatus(globals.RDB, token)
		if errors.Is(err, redis.Nil) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 4003,
				"data": nil,
				"err":  "CSRF token 已过期或无效",
			})
			return
		} else if err != nil {
			fmt.Printf("Failed to get CSRF token from Redis: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"data": nil,
				"err":  "服务器内部错误",
			})
			return
		}

		if status != "unused" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 4003,
				"data": nil,
				"err":  "CSRF token 已使用",
			})
			return
		}

		// 标记 token 为已使用
		err = markCSRFTokenAsUsed(globals.RDB, token)
		if err != nil {
			fmt.Printf("Failed to mark CSRF token as used in Redis: %v\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"data": nil,
				"err":  "服务器内部错误",
			})
			return
		}

		// 继续进行原始的 CSRF 验证
		adapter.Wrap(csrfMd)(c)
	}
}
