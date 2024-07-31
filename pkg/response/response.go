package response

import (
	"forum/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(ctx *gin.Context, data *AppData) {
	ctx.JSON(http.StatusOK, data)
}
func Failed(ctx *gin.Context, errors errors.AppError) {
	ctx.JSON(http.StatusBadRequest, errors)
}
