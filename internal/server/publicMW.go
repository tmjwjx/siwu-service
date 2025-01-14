package server

import (
	"forum/pkg/crossMW"
	"github.com/gin-gonic/gin"
)

// HandlePublicMW 处理公共中间件
func HandlePublicMW(e *gin.Engine) {
	// 跨域
	e.Use(crossMW.CorsMiddleware())

}
