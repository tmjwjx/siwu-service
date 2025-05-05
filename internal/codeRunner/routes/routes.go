package routes

import (
	"forum/internal/codeRunner/controllers"
	"github.com/gin-gonic/gin"
)

func CodeRunner(c *gin.Engine) {
	c.POST("/codeRunner/execute", controllers.GetCodeInfoCtr)
	c.POST("/codeRunner/getResult", controllers.GetCodeRunResultCtrl)
}
