package flowRestriction

import (
	"fmt"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoginRateLimitMiddleware 登录限流中间件
func LoginRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查限流
		allowed, remain, lockTime, message := CheckRateLimit(globals.RDB, c)

		// 设置响应头
		c.Header("X-Remaining-Attempts", fmt.Sprintf("%d", remain))
		// 次数不足，无法登录
		if !allowed {
			response.Failed(c, http.StatusTooManyRequests, response.NewAppErr(globals.StatusTooManyRequests, fmt.Errorf(message), gin.H{"remain": remain, "lock_time": lockTime}))
			c.Abort()
			return
		}

		c.Next()

		// 如果登录失败，记录尝试
		if c.Writer.Status() == http.StatusUnauthorized {
			RecordFailedAttempt(globals.RDB, c)
		}
	}
}
