package inits

import (
	"forum/internal/server"
	"github.com/gin-gonic/gin"
)

func routerInit() {
	server.Router = gin.Default()
}
