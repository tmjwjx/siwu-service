package controllers

import (
	"forum/internal/image/logics"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadHandlerControllers(c *gin.Context) {
	// 解析multipart/form-data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 设置最大大小为 32MB
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse form"})
		return
	}
	// 具体逻辑实现
	logics.UploadHandlerLogic(c)
}
