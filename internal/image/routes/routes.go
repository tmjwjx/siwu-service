package routes

import (
	"forum/internal/image/controllers"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

// ProduceImageUrl 生成图片的url
func ProduceImageUrl(e *gin.Engine) {

	// 根据上传的图片，生成url路径，并将图片存到静态文件中，url路径存到数据库中
	e.POST("/produce_image_url", controllers.ProduceUrlCtrl)

	// 压缩图片
	e.GET(globals.SConfig.Prefix+"/:path", controllers.CompressImageCtrl)

}
