package response

import (
	"github.com/gin-gonic/gin"
)

// status 指的是 http.StatusOK等之类的状态码

func Success(ctx *gin.Context, status int, data *AppData) {
	ctx.JSON(status, data)
}

func Failed(ctx *gin.Context, status int, err *AppErr) {
	ctx.JSON(status, err)
}
