package routes

import (
	"forum/internal/api/controlles"
	"github.com/gin-gonic/gin"
)

func Api(e *gin.Engine) {

	r := e.Group("/api")
	{
		r.GET("/init", controlles.ApiInitCtrl)
	}
}
