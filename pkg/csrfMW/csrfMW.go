package csrfMW

import (
	"errors"
	"fmt"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"net/http"
	"time"
)

/*
	当前端发送 Get 请求的时候会给前端发送一个 X-CSRF-Token 头部，当前端发送 Post、Put、Delete 等需要修改数据的请求的时候，后端会验证 X-CSRF-Token。
	CSRFTokenMW 负责生成 Token 并设置响应头，CSRFMW 负责验证。使用时需要确保 CSRFTokenMW 注册在 CSRFMW 之后，因为：
	csrf.Protect 生成的中间件会先执行，生成 Token 并存储到上下文。后续的 CSRFTokenMW 才能通过 csrf.Token(c.Request) 获取到这个 Token 并发送给前端。
*/

var (
	csrfMd       func(http.Handler) http.Handler
	tokenExpires = 5 * time.Minute // 设置 Token 过期时间为 5 分钟
)

// CSRFTokenMW 在响应中提供 CSRF 令牌给前端。
func CSRFTokenMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 仅在 GET 请求时生成 Token
		if c.Request.Method == http.MethodGet {
			// 从请求上下文中获取 CSRF Token
			token := csrf.Token(c.Request)

			// 将 Token 存储到 Redis 中，并设置过期时间
			err := SetCSRFTokenToRedis(globals.RDB, token, tokenExpires)
			if err != nil {
				globals.Log.Error(response.ErrStorageToRedis)
				response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrStorageToRedis), nil))
				c.Abort()
				return
			}

			// 将 Token 添加到响应头中
			c.Header("X-CSRF-Token", token)
			// 手动设置 CSRF-Cookie
			// c.SetCookie("CSRF-Cookie", token, int(tokenExpires.Seconds()), "/", "", false, false)
			// fmt.Println(token)
		}
	}
}

// CSRFMW 验证前端携带的 CSRF 令牌的有效性，确保应用免受跨站请求伪造攻击。
func CSRFMW() gin.HandlerFunc {
	// 初始化 CSRF 中间件
	// csrf.Protect 中间件的核心功能：生成 Token 并注入到上下文（供 csrf.Token(c.Request) 获取）。验证 Token。
	csrfMd = csrf.Protect(
		[]byte("8d7c2c6a1d4b7a3d9c8e4f6b5a2d1c3e4f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c"), // 32 字节的随机密钥
		csrf.Secure(false),   // 允许在 HTTP 中使用
		csrf.HttpOnly(false), // 设置 HttpOnly 属性，设置为false，让js可以获取到
		csrf.CookieName("CSRF-Cookie"),
		csrf.Path("/"),
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
	csrfHandler := adapter.Wrap(csrfMd)

	return func(c *gin.Context) {
		// 如果请求已被终止（如 Token 无效），直接返回
		if c.IsAborted() {
			return
		}

		// 对于非 GET 请求，执行自定义的 Redis 验证逻辑
		if c.Request.Method != http.MethodGet {
			token := c.GetHeader("X-CSRF-Token")

			if token == "" {
				globals.Log.Error("缺少 CSRF token")
				response.Failed(c, http.StatusForbidden, response.NewAppErr(globals.StatusForbidden, fmt.Errorf("缺少 CSRF token"), nil))
				c.Abort()
				return
			}

			// 检查 Token 是否存在且未使用
			status, err := checkCSRFTokenStatus(globals.RDB, token)
			if errors.Is(err, redis.Nil) {
				globals.Log.Error("CSRF token 已过期或无效")
				response.Failed(c, http.StatusForbidden, response.NewAppErr(globals.StatusForbidden, fmt.Errorf("CSRF token 已过期或无效"), nil))
				c.Abort()
				return

			} else if err != nil {
				globals.Log.Errorf("服务器内部错误,从Redis获取CSRF令牌失败:%v", err)
				response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("服务器内部错误"), nil))
				c.Abort()
				return
			}

			if status != "unused" {
				globals.Log.Error("CSRF token 已使用")
				response.Failed(c, http.StatusForbidden, response.NewAppErr(globals.StatusForbidden, fmt.Errorf("CSRF token 已过期或无效"), nil))
				c.Abort()
				return
			}

			// 标记 Token 为已使用
			err = markCSRFTokenAsUsed(globals.RDB, token, tokenExpires)
			if err != nil {
				globals.Log.Errorf("服务器内部错误,未能将CSRF令牌标记为在Redis中使用:%v", err)
				response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("服务器内部错误"), nil))
				c.Abort()
				return
			}
		}
		// 先执行原始的 CSRF 中间件（生成 Token 并验证）
		csrfHandler(c)
	}
}
