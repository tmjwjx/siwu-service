package csrfMW

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"net/http"
)

var csrfMd func(http.Handler) http.Handler

// CSRFTokenMW 在响应中提供 CSRF 令牌给前端。
// 在每个 HTTP 响应的头部添加 X-CSRF-Token 字段，该字段的值是根据当前请求生成的 CSRF（Cross-Site Request Forgery，跨站请求伪造）令牌。前端可以从响应头部获取这个令牌，并在后续的请求中携带该令牌，以通过后端的 CSRF 验证。
func CSRFTokenMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-CSRF-Token", csrf.Token(c.Request))
	}
}

// CSRFMW 验证前端携带的 CSRF 令牌的有效性，确保应用免受跨站请求伪造攻击。
func CSRFMW() gin.HandlerFunc {
	// 初始化 CSRF 中间件
	csrfMd = csrf.Protect(
		[]byte("8d7c2c6a1d4b7a3d9c8e4f6b5a2d1c3e4f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c"),
		csrf.Secure(false),  // 是否仅在 HTTPS 连接中使用 CSRF 保护。false 表示在 HTTP 连接中也可以使用。
		csrf.HttpOnly(true), // 设置 CSRF 令牌的 HttpOnly 属性为 true，这样 JavaScript 无法访问该令牌，提高安全性。
		// 指定当 CSRF 验证失败时的错误处理函数。当 CSRF 令牌无效时，会调用这个函数。
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{
				"code": 4003,
				"data": {},
				"err": "无效的 CSRF token"
			}`))
		})),
	)

	// 转换为 Gin 中间件
	return adapter.Wrap(csrfMd)
}

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
