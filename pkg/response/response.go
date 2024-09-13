package response

import (
	"github.com/gin-gonic/gin"
)

// status 指的是 http.StatusOK等之类的状态码

func Success(c *gin.Context, status int, data *AppData) {
	c.JSON(status, data)
}

func Failed(c *gin.Context, status int, err *AppErr) {
	c.JSON(status, err)
}