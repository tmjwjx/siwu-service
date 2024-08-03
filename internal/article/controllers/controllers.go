package controllers

import (
	"forum/internal/article/logic"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	logic.Search(c)
}
func Test1(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "ok"})
}
