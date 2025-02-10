package recovery

import (
	"fmt"
	"forum/pkg/globals"
	"forum/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CustomRecovery recover panic
// @Description: 在发生 panic 时，将错误打印到日志中，并返回 500。
// @Author lizhuang 2025-01-23 09:47:40
// @return       gin.HandlerFunc
func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录到日志文件
				globals.Log.Panicf("Panic recovered: %v\n", err)

				// 返回 500 错误
				response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("服务器内部错误"), nil))

				c.Abort()
			}
		}()
		c.Next()
	}
}
