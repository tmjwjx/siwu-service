package response

import (
	"github.com/gin-gonic/gin"
)

// status 指的是 http.StatusOK等之类的状态码

func Success(ctx *gin.Context, data *AppData, status int) {
	ctx.JSON(status, data)
}

func Failed(ctx *gin.Context, err *AppErr, status int) {
	ctx.JSON(status, err)
}
