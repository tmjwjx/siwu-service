package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(ctx *gin.Context, data *AppData) {
	ctx.JSON(http.StatusOK, data)
}

func Failed(ctx *gin.Context, err *AppErr) {
	ctx.JSON(http.StatusBadRequest, err)
}
