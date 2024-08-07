package controllers

import (
	"forum/internal/article/logics"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	logics.Search(c)
}
