package inits

import (
	"forum/internal/server"
	"forum/pkg/utils"
	"github.com/gin-gonic/gin"
)

func Init() {
	// Initialize Router
	server.Router = gin.Default()

	// Initialize MYSQL
	mysqlInit()

	utils.InitFile("logs", "forum")
}
