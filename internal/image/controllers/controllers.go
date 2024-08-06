package controllers

import (
	"forum/internal/image/logics"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UploadHandlerControllers 将前端传过来的图片文件存到文件系统中
// kind 用来判断图片是用户的还是文章的, ID 是用户或文章的ID
func UploadHandlerControllers(c *gin.Context, kind int, ID uint) {
	// 将文件系统中的目录映射到 URL 路径
	globals.Router.Static("/images", "./static/images")

	// 解析multipart/form-data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 设置最大大小为 32MB
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse form"})
		return
	}
	// 具体逻辑实现
	logics.UploadHandlerLogic(c, kind, ID)
}
