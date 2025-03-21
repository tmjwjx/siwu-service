package routes

import (
	"forum/internal/AI/controllers"
	"github.com/gin-gonic/gin"
)

func AIRouters(c *gin.Engine) {
	c.POST("/AI/codeExplain", controllers.GetCodeExplain)
}
