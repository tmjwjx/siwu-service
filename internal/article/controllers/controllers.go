package controllers

import (
	"forum/internal/article/logics"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	logics.Search(c)
}
func Test1(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "ok"})
}
