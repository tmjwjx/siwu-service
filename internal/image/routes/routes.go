package routes

import (
	"forum/internal/image/controllers"
	"github.com/gin-gonic/gin"
)

// ProduceImageUrl 生成图片的url
func ProduceImageUrl(e *gin.Engine) {

	e.POST("/produce_image_url", controllers.ProduceUrlCtrl)

}
