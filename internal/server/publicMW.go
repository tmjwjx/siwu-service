package server

import (
	"forum/pkg/corsMW"
	"github.com/gin-gonic/gin"
)

// HandlePublicMW 处理公共中间件
func HandlePublicMW(e *gin.Engine) {
	// 跨域
	e.Use(corsMW.CorsMiddleware())

}
