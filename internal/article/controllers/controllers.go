package controllers

import (
	"forum/internal/article/logic"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	logic.Search(c)
}
